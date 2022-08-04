package loader

import (
	jschema "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/errors"
	"github.com/jsightapi/jsight-schema-go-library/internal/lexeme"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema/internal/schema"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema/internal/schema/constraint"
)

// orRuleSetLoader loads data from rule-set into the type. Specifies the name of
// this user type in the parent node (in types list constraint).
type orRuleSetLoader struct {
	// The node.
	node schema.Node

	// A rootSchema to which the type from the "or" rule will be added.
	rootSchema *schema.Schema

	// embeddedValueLoader a loader for "enum" value.
	embeddedValueLoader embeddedLoader

	// rules all available rules.
	rules map[string]jschema.Rule

	// stateFunc a function for running a state machine (the current state of the
	// state machine).
	stateFunc func(lexeme.LexEvent)

	// typeRoot a node (type mixed) to which constraints from rule-set are
	// added. This node will become the root node for the type created from
	// the rule-set.
	typeRoot *schema.MixedNode

	// ruleNameLex the last found key in rule-set.
	ruleNameLex lexeme.LexEvent

	// inProgress indicates are we already done or not.
	inProgress bool
}

var _ embeddedLoader = (*orRuleSetLoader)(nil)

// Loader for rule-set value. Ex: {type: "integer", min: 0}
func newOrRuleSetLoader(
	node schema.Node,
	rootSchema *schema.Schema,
	rules map[string]jschema.Rule,
) *orRuleSetLoader {
	if _, ok := node.(*schema.MixedValueNode); ok {
		panic(errors.ErrCannotSpecifyOtherRulesWithTypeReference)
	}

	s := &orRuleSetLoader{
		node:       node,
		rootSchema: rootSchema,
		rules:      rules,
		typeRoot:   schema.NewMixedNode(node.BasisLexEventOfSchemaForNode()),
		inProgress: true,
	}
	s.stateFunc = s.objectBegin
	return s
}

func (s *orRuleSetLoader) Load(lex lexeme.LexEvent) bool {
	defer lexeme.CatchLexEventError(lex)
	s.stateFunc(lex)
	return s.inProgress
}

func (s *orRuleSetLoader) embeddedLoad(lex lexeme.LexEvent) {
	if !s.embeddedValueLoader.Load(lex) {
		s.embeddedValueLoader = nil
		s.stateFunc = s.valueEnd
		if lex.Type() == lexeme.TypesShortcutEnd {
			s.stateFunc = s.afterShortcutEnd
		}
	}
}

func (s *orRuleSetLoader) afterShortcutEnd(lex lexeme.LexEvent) {
	if lex.Type() != lexeme.MixedValueEnd {
		panic(errors.ErrLoader)
	}
	s.stateFunc = s.valueEnd
}

// objectBegin begin of object "{"
func (s *orRuleSetLoader) objectBegin(lex lexeme.LexEvent) {
	if lex.Type() != lexeme.ObjectBegin {
		panic(errors.ErrLoader)
	}
	s.stateFunc = s.keyOrObjectEnd
}

// keyOrObjectEnd object key or object end
// ex: {"key" <--
// ex: {...} <--
func (s *orRuleSetLoader) keyOrObjectEnd(lex lexeme.LexEvent) {
	switch lex.Type() {
	case lexeme.ObjectKeyBegin:
		return
	case lexeme.ObjectKeyEnd:
		s.ruleNameLex = lex
		s.stateFunc = s.valueBegin
		if s.ruleNameLex.Value().String() == "enum" {
			s.stateFunc = s.enumValueBegin
		}
	case lexeme.ObjectEnd:
		s.stateFunc = s.endOfLoading
		s.inProgress = false
		s.makeTypeFromRuleSet()
	default:
		panic(errors.ErrLoader)
	}
}

// valueBegin object value begin
// ex: {"key": <--
func (s *orRuleSetLoader) valueBegin(lex lexeme.LexEvent) {
	if lex.Type() != lexeme.ObjectValueBegin {
		panic(errors.ErrLoader)
	}
	s.stateFunc = s.valueLiteral
}

func (s *orRuleSetLoader) enumValueBegin(lex lexeme.LexEvent) {
	if lex.Type() != lexeme.ObjectValueBegin {
		panic(errors.ErrLoader)
	}
	enumConstraint := constraint.NewEnum()
	s.typeRoot.AddConstraint(enumConstraint)
	s.embeddedValueLoader = newEnumValueLoader(enumConstraint, s.rules)
	s.stateFunc = s.embeddedLoad
}

// valueLiteral literal value
// ex: {"key": ... <--
func (s *orRuleSetLoader) valueLiteral(lex lexeme.LexEvent) {
	switch lex.Type() {
	case lexeme.LiteralBegin:
		return
	case lexeme.LiteralEnd:
		c := constraint.NewConstraintFromRule(s.ruleNameLex, lex.Value(), s.node.Value()) // can panic
		s.typeRoot.AddConstraint(c)
		s.stateFunc = s.valueEnd
	default:
		panic(errors.ErrLiteralValueExpected)
	}
}

// valueEnd object value end
// ex: {"key": "value" <--
func (s *orRuleSetLoader) valueEnd(lex lexeme.LexEvent) {
	if lex.Type() != lexeme.ObjectValueEnd {
		panic(errors.ErrLoader)
	}
	s.stateFunc = s.keyOrObjectEnd
}

// endOfLoading the method should not be called during normal operation. Ensures
// that the loader will not continue to work after the load is complete.
func (*orRuleSetLoader) endOfLoading(lexeme.LexEvent) {
	panic(errors.ErrLoader)
}

// return TypesList constraint for node
func (s *orRuleSetLoader) nodeTypesListConstraint() *constraint.TypesList {
	c := s.node.Constraint(constraint.TypesListConstraintType)
	if c == nil {
		panic(errors.ErrLoader) // constraint not found
	}
	return c.(*constraint.TypesList)
}

// makeTypeFromRuleSet appends new type based on rule-set.
func (s *orRuleSetLoader) makeTypeFromRuleSet() {
	if s.typeRoot.NumberOfConstraints() == 0 {
		panic(errors.ErrEmptyRuleSet)
	}

	c := s.nodeTypesListConstraint()
	an := s.makeTypeASTNode(c.Source())

	typeConstraint := s.typeRoot.Constraint(constraint.TypeConstraintType)
	if typeConstraint != nil && s.typeRoot.NumberOfConstraints() == 1 {
		typeValue := typeConstraint.(constraint.BytesKeeper).Bytes().Unquote()
		if typeValue.IsUserTypeName() {
			c.AddNameWithASTNode(typeValue.String(), typeValue.String(), an)
			return
		}
	}

	typ := schema.New()
	typ.SetRootNode(s.typeRoot)

	CompileBasic(&typ, false)

	lex := s.node.BasisLexEventOfSchemaForNode()
	name := s.rootSchema.AddUnnamedType(&typ, lex.File(), lex.Begin())

	c.AddNameWithASTNode(name, s.typeRoot.Type().String(), an)
}

func (s *orRuleSetLoader) makeTypeASTNode(
	source jschema.RuleASTNodeSource,
) jschema.RuleASTNode {
	cc := s.typeRoot.ConstraintMap()

	an := jschema.RuleASTNode{
		JSONType:   jschema.JSONTypeObject,
		Properties: jschema.MakeRuleASTNodes(cc.Len()),
		Source:     source,
	}

	cc.EachSafe(func(k constraint.Type, v constraint.Constraint) {
		an.Properties.Set(k.String(), v.ASTNode())
	})
	return an
}

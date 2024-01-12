package loader

import (
	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/errs"

	"github.com/jsightapi/jsight-schema-core/json"
	"github.com/jsightapi/jsight-schema-core/lexeme"
	"github.com/jsightapi/jsight-schema-core/notations/jschema/ischema"
	"github.com/jsightapi/jsight-schema-core/notations/jschema/ischema/constraint"
)

type orValueLoader struct {
	// A node to add type constraint.
	node ischema.Node

	// A rootSchema into which types can be added from the "or" rule.
	rootSchema *ischema.ISchema

	// rules all available rules.
	rules map[string]schema.Rule

	// stateFunc a function for running a state machine (the current state of the
	// state machine).
	stateFunc func(lexeme.LexEvent)

	// ruleSetLoader a loader for rule-set value. Ex: {type: "integer", min: 0}.
	ruleSetLoader *orRuleSetLoader

	// inProgress indicates are we already done or not.
	inProgress bool
}

var _ embeddedLoader = (*orValueLoader)(nil)

// newOrValueLoader creates loader for "or" rule value (array of rule-sets).
// Ex: [{type: "@typeName-1"}, "@typeName-2", {type: "integer", min: 0}]
func newOrValueLoader(
	node ischema.Node,
	rootSchema *ischema.ISchema,
	rules map[string]schema.Rule,
) *orValueLoader {
	a := &orValueLoader{
		node:       node,
		rootSchema: rootSchema,
		rules:      rules,
		inProgress: true,
	}
	a.stateFunc = a.begin
	return a
}

func (a *orValueLoader) Load(lex lexeme.LexEvent) bool {
	defer lexeme.CatchLexEventError(lex)
	if a.ruleSetLoader != nil {
		if !a.ruleSetLoader.Load(lex) {
			a.ruleSetLoader = nil
		}
	} else {
		a.stateFunc(lex)
	}
	return a.inProgress
}

// nodeTypesListConstraint returns TypesList constraint for node.
func (a *orValueLoader) nodeTypesListConstraint() *constraint.TypesList {
	c := a.node.Constraint(constraint.TypesListConstraintType)
	if c == nil {
		panic(errs.ErrLoader.F()) // constraint not found
	}
	return c.(*constraint.TypesList)
}

// begin of array "["
func (a *orValueLoader) begin(lex lexeme.LexEvent) {
	if lex.Type() != lexeme.ArrayBegin {
		panic(errs.ErrArrayWasExpectedInOrRule.F())
	}
	a.stateFunc = a.itemBeginOrArrayEnd
}

// itemBeginOrArrayEnd begin of array item or array end
// ex: [{ <--
// ex: [" <--
// ex: ] <--
func (a *orValueLoader) itemBeginOrArrayEnd(lex lexeme.LexEvent) {
	switch lex.Type() {
	case lexeme.ArrayItemBegin:
		a.stateFunc = a.itemInner
	case lexeme.ArrayEnd:
		switch a.nodeTypesListConstraint().Len() {
		case 0:
			panic(errs.ErrEmptyArrayInOrRule.F())
		case 1:
			panic(errs.ErrOneElementInArrayInOrRule.F())
		}
		a.stateFunc = a.endOfLoading
		a.inProgress = false
	default:
		panic(errs.ErrLoader.F())
	}
}

// itemInner array item value (literal begin or object begin)
// ex: [{ <--
// ex: [" <--
func (a *orValueLoader) itemInner(lex lexeme.LexEvent) {
	switch lex.Type() {
	case lexeme.LiteralBegin:
		a.stateFunc = a.literal
	case lexeme.ObjectBegin:
		a.ruleSetLoader = newOrRuleSetLoader(a.node, a.rootSchema, a.rules)
		a.ruleSetLoader.Load(lex)
		a.stateFunc = a.itemEnd
	default:
		panic(errs.ErrIncorrectArrayItemTypeInOrRule.F()) // ex: array
	}
}

// literal parse name of user or JSON type.
// ex: ["@type" <--
// ex: ["string" <--
func (a *orValueLoader) literal(lex lexeme.LexEvent) {
	if lex.Type() != lexeme.LiteralEnd {
		panic(errs.ErrLoader.F())
	}

	if json.Guess(lex.Value()).LiteralJsonType() != json.TypeString {
		panic(errs.ErrIncorrectArrayItemTypeInOrRule.F())
	}

	val := lex.Value().Unquote()
	valStr := val.String()

	c := a.nodeTypesListConstraint()
	if val.IsUserTypeName() {
		c.AddNameWithASTNode(valStr, valStr, schema.RuleASTNode{
			TokenType:  schema.TokenTypeShortcut,
			Value:      valStr,
			Properties: &schema.RuleASTNodes{},
			Source:     schema.RuleASTNodeSourceManual,
		})
	} else {
		root := ischema.NewMixedNode(a.node.BasisLexEventOfSchemaForNode())
		root.AddConstraint(constraint.NewType(val, schema.RuleASTNodeSourceManual))

		typ := ischema.New()
		typ.SetRootNode(root)

		CompileBasic(&typ, false)

		lex := a.node.BasisLexEventOfSchemaForNode()
		name := a.rootSchema.AddUnnamedType(&typ, lex.File(), lex.Begin())

		a.
			nodeTypesListConstraint().
			AddName(name, string(root.SchemaType()), schema.RuleASTNodeSourceManual)
	}

	a.stateFunc = a.itemEnd
}

// itemEnd array item end
// ex: ["@type" <--
// ex: [{...} <--
func (a *orValueLoader) itemEnd(lex lexeme.LexEvent) {
	if lex.Type() != lexeme.ArrayItemEnd {
		panic(errs.ErrLoader.F())
	}
	a.stateFunc = a.itemBeginOrArrayEnd
}

// endOfLoading the method should not be called during normal operation. Ensures
// that the loader will not continue to work after the load is complete.
func (*orValueLoader) endOfLoading(lexeme.LexEvent) {
	panic(errs.ErrLoader.F())
}

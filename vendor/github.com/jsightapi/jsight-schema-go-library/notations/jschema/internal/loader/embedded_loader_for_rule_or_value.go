package loader

import (
	jschema "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/errors"
	"github.com/jsightapi/jsight-schema-go-library/internal/json"
	"github.com/jsightapi/jsight-schema-go-library/internal/lexeme"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema/internal/schema"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema/internal/schema/constraint"
)

type orValueLoader struct {
	// A node to add type constraint.
	node schema.Node

	// A rootSchema into which types can be added from the "or" rule.
	rootSchema *schema.Schema

	// stateFunc a function for running a state machine (the current state of the
	// state machine).
	stateFunc func(lexeme.LexEvent)

	// ruleSetLoader a loader for rule-set value. Ex: {type: "integer", min: 0}.
	ruleSetLoader *orRuleSetLoader

	// inProgress indicates are we already done or not.
	inProgress bool
}

// newOrValueLoader creates loader for "or" rule value (array of rule-sets).
// Ex: [{type: "@typeName-1"}, "@typeName-2", {type: "integer", min: 0}]
func newOrValueLoader(node schema.Node, rootSchema *schema.Schema) embeddedLoader {
	a := &orValueLoader{
		node:       node,
		rootSchema: rootSchema,
		inProgress: true,
	}
	a.stateFunc = a.begin
	return a
}

// Returns false when the load is complete.
func (a *orValueLoader) load(lex lexeme.LexEvent) bool {
	defer lexeme.CatchLexEventError(lex)
	if a.ruleSetLoader != nil {
		if !a.ruleSetLoader.load(lex) {
			a.ruleSetLoader = nil
		}
	} else {
		a.stateFunc(lex)
	}
	return a.inProgress
}

// return TypesList constraint for node
func (a *orValueLoader) nodeTypesListConstraint() *constraint.TypesList {
	c := a.node.Constraint(constraint.TypesListConstraintType)
	if c == nil {
		panic(errors.ErrLoader) // constraint not found
	}
	return c.(*constraint.TypesList)
}

// begin of array "["
func (a *orValueLoader) begin(lex lexeme.LexEvent) {
	if lex.Type() != lexeme.ArrayBegin {
		panic(errors.ErrArrayWasExpectedInOrRule)
	}
	a.stateFunc = a.itemBeginOrArrayEnd
}

// begin of array item or array end
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
			panic(errors.ErrEmptyArrayInOrRule)
		case 1:
			panic(errors.ErrOneElementInArrayInOrRule)
		}
		a.stateFunc = a.endOfLoading
		a.inProgress = false
	default:
		panic(errors.ErrLoader)
	}
}

// array item value (literal begin or object begin)
// ex: [{ <--
// ex: [" <--
func (a *orValueLoader) itemInner(lex lexeme.LexEvent) {
	switch lex.Type() {
	case lexeme.LiteralBegin:
		a.stateFunc = a.literal
	case lexeme.ObjectBegin:
		a.ruleSetLoader = newOrRuleSetLoader(a.node, a.rootSchema)
		a.ruleSetLoader.load(lex)
		a.stateFunc = a.itemEnd
	default:
		panic(errors.ErrIncorrectArrayItemTypeInOrRule) // ex: array
	}
}

// literal parse name of user or JSON type.
// ex: ["@type" <--
// ex: ["string" <--
func (a *orValueLoader) literal(lex lexeme.LexEvent) {
	if lex.Type() != lexeme.LiteralEnd {
		panic(errors.ErrLoader)
	}

	if json.Guess(lex.Value()).LiteralJsonType() != json.TypeString {
		panic(errors.ErrIncorrectArrayItemTypeInOrRule)
	}

	val := lex.Value().Unquote()

	c := a.nodeTypesListConstraint()
	if val.IsUserTypeName() {
		c.AddName(val.String(), val.String(), jschema.RuleASTNodeSourceManual)
	} else {
		root := schema.NewMixedNode(a.node.BasisLexEventOfSchemaForNode())
		root.AddConstraint(constraint.NewType(val, jschema.RuleASTNodeSourceManual))

		typ := schema.New()
		typ.SetRootNode(root)

		CompileBasic(&typ, false)

		lex := a.node.BasisLexEventOfSchemaForNode()
		name := a.rootSchema.AddUnnamedType(&typ, lex.File(), lex.Begin())

		a.
			nodeTypesListConstraint().
			AddName(name, root.Type().String(), jschema.RuleASTNodeSourceManual)
	}

	a.stateFunc = a.itemEnd
}

// array item end
// ex: ["@type" <--
// ex: [{...} <--
func (a *orValueLoader) itemEnd(lex lexeme.LexEvent) {
	if lex.Type() != lexeme.ArrayItemEnd {
		panic(errors.ErrLoader)
	}
	a.stateFunc = a.itemBeginOrArrayEnd
}

// endOfLoading the method should not be called during normal operation. Ensures
// that the loader will not continue to work after the load is complete.
func (*orValueLoader) endOfLoading(lexeme.LexEvent) {
	panic(errors.ErrLoader)
}

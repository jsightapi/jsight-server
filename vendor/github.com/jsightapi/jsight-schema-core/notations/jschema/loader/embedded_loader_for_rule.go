package loader

import (
	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/errs"
	"github.com/jsightapi/jsight-schema-core/lexeme"
	"github.com/jsightapi/jsight-schema-core/notations/jschema/ischema"
	"github.com/jsightapi/jsight-schema-core/notations/jschema/ischema/constraint"
)

// ruleLoader responsible for creating constraints for SCHEMA internal representation
// nodes from the RULES described in the SCHEMA file.
type ruleLoader struct {
	// A node to add constraint.
	node ischema.Node

	// rootSchema a scheme into which types can be added from the "or" rule.
	rootSchema *ischema.ISchema

	// rules all available rules.
	rules map[string]schema.Rule

	// stateFunc a function for running a state machine (the current state of the
	// state machine) to parse RULE that occur in the schema.
	stateFunc func(lexeme.LexEvent)

	// embeddedValueLoader a loader for "or" and "enum" value.
	embeddedValueLoader embeddedLoader

	// ruleNameLex the last found object key.
	ruleNameLex lexeme.LexEvent

	// nodesPerCurrentLineCount the number of nodes in a line. To check because
	// the rule cannot be added if there is more than one node suitable for this
	// in the row.
	nodesPerCurrentLineCount uint
}

func newRuleLoader(
	node ischema.Node,
	nodesPerCurrentLineCount uint,
	rootSchema *ischema.ISchema,
	rules map[string]schema.Rule,
) *ruleLoader {
	rl := &ruleLoader{
		node:                     node,
		rootSchema:               rootSchema,
		nodesPerCurrentLineCount: nodesPerCurrentLineCount,
		rules:                    rules,
	}
	rl.stateFunc = rl.begin
	return rl
}

func (rl *ruleLoader) load(lex lexeme.LexEvent) {
	defer lexeme.CatchLexEventError(lex)
	rl.stateFunc(lex)
}

func (rl *ruleLoader) begin(lex lexeme.LexEvent) {
	switch lex.Type() {
	case lexeme.NewLine:
		// Do nothing

	case lexeme.InlineAnnotationTextBegin, lexeme.MultiLineAnnotationTextBegin:
		rl.stateFunc = rl.commentTextEnd

	case lexeme.ObjectBegin:
		rl.stateFunc = rl.ruleKeyOrObjectEnd

	default:
		panic(errs.ErrLoader.F())
	}
}

func (rl *ruleLoader) commentTextBegin(lex lexeme.LexEvent) {
	switch lex.Type() {
	case lexeme.NewLine:
		// Do nothing keep waiting for annotation start.

	case lexeme.InlineAnnotationTextBegin, lexeme.MultiLineAnnotationTextBegin:
		rl.stateFunc = rl.commentTextEnd
	default:
		panic(errs.ErrLoader.F())
	}
}

func (rl *ruleLoader) commentTextEnd(lex lexeme.LexEvent) {
	switch lex.Type() {
	case lexeme.InlineAnnotationTextEnd, lexeme.MultiLineAnnotationTextEnd:
		if rl.node != nil {
			rl.node.SetComment(lex.Value().TrimSpaces().String())
		}
		rl.stateFunc = rl.endOfLoading
	default:
		panic(errs.ErrLoader.F())
	}
}

func (rl *ruleLoader) ruleKeyOrObjectEnd(lex lexeme.LexEvent) {
	switch lex.Type() {
	case lexeme.ObjectKeyBegin, lexeme.NewLine:
	case lexeme.ObjectKeyEnd:
		rl.ruleNameLex = lex
		rl.stateFunc = rl.ruleValueBegin
	case lexeme.ObjectEnd:
		rl.stateFunc = rl.commentTextBegin
	default:
		panic(errs.ErrLoader.F())
	}
}

func (rl *ruleLoader) objectEndAfterRuleName(lex lexeme.LexEvent) {
	switch lex.Type() {
	case lexeme.ObjectKeyBegin, lexeme.ObjectValueEnd, lexeme.NewLine:
	case lexeme.ObjectKeyEnd:
		rl.ruleNameLex = lex
		rl.stateFunc = rl.ruleValueBegin
	case lexeme.ObjectEnd:
		rl.stateFunc = rl.commentTextBegin
	default:
		panic(errs.ErrLoader.F())
	}
}

func (rl *ruleLoader) ruleValueBegin(lex lexeme.LexEvent) {
	if lex.Type() != lexeme.ObjectValueBegin {
		panic(errs.ErrLoader.F())
	}
	rl.stateFunc = rl.ruleValue
}

func (rl *ruleLoader) ruleValue(lex lexeme.LexEvent) {
	if rl.nodesPerCurrentLineCount == 0 {
		panic(errs.ErrIncorrectRuleWithoutExample.F())
	} else if rl.nodesPerCurrentLineCount != 1 {
		panic(errs.ErrIncorrectRuleForSeveralNode.F())
	}

	ruleName := rl.ruleNameLex.Value().TrimSpaces().Unquote().String()

	switch ruleName {
	case "or":
		rl.node.AddConstraint(constraint.NewTypesList(schema.RuleASTNodeSourceManual))
		rl.node.AddConstraint(constraint.NewOr(schema.RuleASTNodeSourceManual)) // Used for compile-time checking.
		rl.embeddedValueLoader = newOrValueLoader(rl.node, rl.rootSchema, rl.rules)
		rl.stateFunc = rl.loadEmbeddedValue
		rl.stateFunc(lex)

	case "enum":
		enumConstraint := constraint.NewEnum()
		rl.node.AddConstraint(enumConstraint)
		rl.embeddedValueLoader = newEnumValueLoader(enumConstraint, rl.rules)
		rl.stateFunc = rl.loadEmbeddedValue
		rl.stateFunc(lex)

	case "allOf":
		allOfConstraint := constraint.NewAllOf()
		rl.node.AddConstraint(allOfConstraint)
		rl.embeddedValueLoader = newAllOfValueLoader(allOfConstraint)
		rl.stateFunc = rl.loadEmbeddedValue
		rl.stateFunc(lex)

	default:
		if lex.Type() != lexeme.LiteralBegin {
			panic(errs.ErrIncorrectRuleValueType.F())
		}

		rl.stateFunc = rl.ruleValueLiteral
	}
}

func (rl *ruleLoader) ruleValueLiteral(ruleValue lexeme.LexEvent) {
	if ruleValue.Type() != lexeme.LiteralEnd {
		panic(errs.ErrLoader.F())
	}
	c := constraint.NewConstraintFromRule(rl.ruleNameLex, ruleValue.Value(), rl.node.Value()) // can panic
	rl.node.AddConstraint(c)

	rl.stateFunc = rl.ruleValueEnd
}

func (rl *ruleLoader) loadEmbeddedValue(lex lexeme.LexEvent) {
	if lex.Type() == lexeme.NewLine {
		return
	}
	if !rl.embeddedValueLoader.Load(lex) {
		rl.embeddedValueLoader = nil
		rl.stateFunc = rl.ruleValueEnd
	}
}

func (rl *ruleLoader) ruleValueEnd(lex lexeme.LexEvent) {
	switch lex.Type() {
	case lexeme.ObjectValueEnd:
		rl.stateFunc = rl.ruleKeyOrObjectEnd
	case lexeme.MixedValueEnd:
		rl.stateFunc = rl.objectEndAfterRuleName
	default:
		panic(errs.ErrLoader.F())
	}
}

// The method should not be called during normal operation. Ensures that the loader will not continue to work after
// the load is complete.
func (*ruleLoader) endOfLoading(lexeme.LexEvent) {
	panic(errs.ErrLoader.F())
}

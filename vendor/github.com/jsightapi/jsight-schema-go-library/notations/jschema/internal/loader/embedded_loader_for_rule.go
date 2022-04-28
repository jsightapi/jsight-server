package loader

import (
	jschema "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/errors"
	"github.com/jsightapi/jsight-schema-go-library/internal/lexeme"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema/internal/schema"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema/internal/schema/constraint"
)

// This class is responsible for creating constraints for SCHEMA internal representation nodes from the RULES described
// in the SCHEMA file.

type ruleLoader struct {
	// A node to add constraint.
	node schema.Node

	// rootSchema a scheme into which types can be added from the "or" rule.
	rootSchema *schema.Schema

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

func newRuleLoader(node schema.Node, nodesPerCurrentLineCount uint, rootSchema *schema.Schema) *ruleLoader {
	rl := new(ruleLoader)
	rl.node = node
	rl.rootSchema = rootSchema
	rl.nodesPerCurrentLineCount = nodesPerCurrentLineCount
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
		panic(errors.ErrLoader)
	}
}

func (rl *ruleLoader) commentTextBegin(lex lexeme.LexEvent) {
	switch lex.Type() {
	case lexeme.NewLine:
		// Do nothing keep waiting for annotation start.

	case lexeme.InlineAnnotationTextBegin, lexeme.MultiLineAnnotationTextBegin:
		rl.stateFunc = rl.commentTextEnd
	default:
		panic(errors.ErrLoader)
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
		panic(errors.ErrLoader)
	}
}

func (rl *ruleLoader) ruleKeyOrObjectEnd(lex lexeme.LexEvent) {
	switch lex.Type() {
	case lexeme.ObjectKeyBegin, lexeme.NewLine:
		return
	case lexeme.ObjectKeyEnd:
		rl.ruleNameLex = lex
		rl.stateFunc = rl.ruleValueBegin
	case lexeme.ObjectEnd:
		rl.stateFunc = rl.commentTextBegin
	default:
		panic(errors.ErrLoader)
	}
}

func (rl *ruleLoader) ruleValueBegin(lex lexeme.LexEvent) {
	switch lex.Type() {
	case lexeme.ObjectValueBegin:
		rl.stateFunc = rl.ruleValue
	default:
		panic(errors.ErrLoader)
	}
}

func (rl *ruleLoader) ruleValue(lex lexeme.LexEvent) {
	if rl.nodesPerCurrentLineCount == 0 {
		panic(errors.ErrIncorrectRuleWithoutExample)
	} else if rl.nodesPerCurrentLineCount != 1 {
		panic(errors.ErrIncorrectRuleForSeveralNode)
	}

	ruleName := rl.ruleNameLex.Value().TrimSpaces().Unquote().String()

	if ruleName == "or" { //nolint:gocritic // todo rewrite this logic to switch
		rl.node.AddConstraint(constraint.NewTypesList(jschema.RuleASTNodeSourceManual))
		rl.node.AddConstraint(constraint.NewOr(jschema.RuleASTNodeSourceManual)) // Used for compile-time checking.
		rl.embeddedValueLoader = newOrValueLoader(rl.node, rl.rootSchema)
		rl.stateFunc = rl.loadEmbeddedValue
		rl.stateFunc(lex)
		return
	} else if ruleName == "enum" {
		enumConstraint := constraint.NewEnum()
		rl.node.AddConstraint(enumConstraint)
		rl.embeddedValueLoader = newEnumValueLoader(enumConstraint)
		rl.stateFunc = rl.loadEmbeddedValue
		rl.stateFunc(lex)
		return
	} else if ruleName == "allOf" {
		allOfConstraint := constraint.NewAllOf()
		rl.node.AddConstraint(allOfConstraint)
		rl.embeddedValueLoader = newAllOfValueLoader(allOfConstraint)
		rl.stateFunc = rl.loadEmbeddedValue
		rl.stateFunc(lex)
		return
	} else if lex.Type() == lexeme.LiteralBegin {
		rl.stateFunc = rl.ruleValueLiteral
		return
	}

	panic(errors.ErrIncorrectRuleValueType)
}

func (rl *ruleLoader) ruleValueLiteral(ruleValue lexeme.LexEvent) {
	switch ruleValue.Type() {
	case lexeme.LiteralEnd:
		c := constraint.NewConstraintFromRule(rl.ruleNameLex, ruleValue.Value(), rl.node.Value()) // can panic
		rl.node.AddConstraint(c)

		rl.stateFunc = rl.ruleValueEnd
	default:
		panic(errors.ErrLoader)
	}
}

func (rl *ruleLoader) loadEmbeddedValue(lex lexeme.LexEvent) {
	if lex.Type() == lexeme.NewLine {
		return
	}
	if !rl.embeddedValueLoader.load(lex) {
		rl.embeddedValueLoader = nil
		rl.stateFunc = rl.ruleValueEnd
	}
}

func (rl *ruleLoader) ruleValueEnd(lex lexeme.LexEvent) {
	switch lex.Type() {
	case lexeme.ObjectValueEnd:
		rl.stateFunc = rl.ruleKeyOrObjectEnd
	default:
		panic(errors.ErrLoader)
	}
}

// The method should not be called during normal operation. Ensures that the loader will not continue to work after the load is complete.
func (*ruleLoader) endOfLoading(lexeme.LexEvent) {
	panic(errors.ErrLoader)
}

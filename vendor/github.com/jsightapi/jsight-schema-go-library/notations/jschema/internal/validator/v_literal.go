package validator

import (
	"github.com/jsightapi/jsight-schema-go-library/errors"
	"github.com/jsightapi/jsight-schema-go-library/internal/lexeme"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema/internal/schema"
)

// Validates json according to jSchema's LiteralNode.

type literalValidator struct {
	node_   schema.Node
	parent_ validator
}

func newLiteralValidator(node schema.Node, parent validator) *literalValidator {
	switch node.(type) {
	case *schema.LiteralNode, *schema.MixedNode, *schema.MixedValueNode:
		v := literalValidator{
			node_:   node,
			parent_: parent,
		}
		return &v
	default:
		panic(errors.ErrValidator)
	}
}

func (v literalValidator) node() schema.Node {
	return v.node_
}

func (v literalValidator) parent() validator {
	return v.parent_
}

func (v *literalValidator) setParent(parent validator) {
	v.parent_ = parent
}

// return array (pointers to validators, or nil if not found) and bool (true if validator is done)
func (v *literalValidator) feed(jsonLexeme lexeme.LexEvent) ([]validator, bool) {
	defer lexeme.CatchLexEventError(jsonLexeme)

	switch jsonLexeme.Type() { //nolint:exhaustive // We will throw a panic in over cases.
	case lexeme.LiteralBegin:
		return nil, false
	case lexeme.LiteralEnd:
		ValidateLiteralValue(v.node_, jsonLexeme.Value()) // can panic
		return nil, true
	}

	panic(errors.ErrUnexpectedLexInLiteralValidator)
}

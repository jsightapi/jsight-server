package validator

import (
	"fmt"

	"github.com/jsightapi/jsight-schema-go-library/errors"
	"github.com/jsightapi/jsight-schema-go-library/internal/lexeme"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema/internal/schema"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema/internal/schema/constraint"
)

// Validates json according to jSchema's ArrayNode.

type arrayValidator struct {
	// node_ an array or mixed node.
	node_   schema.Node
	parent_ validator

	// rootSchema the scheme from which it is possible to receive type by their
	// name.
	rootSchema   schema.Schema
	itemsCounter uint
}

func newArrayValidator(node schema.Node, parent validator, rootSchema schema.Schema) *arrayValidator {
	switch node.(type) {
	case *schema.ArrayNode, *schema.MixedNode, *schema.MixedValueNode:
		v := arrayValidator{
			node_:      node,
			parent_:    parent,
			rootSchema: rootSchema,
		}
		return &v
	default:
		panic(errors.ErrValidator)
	}
}

func (v arrayValidator) node() schema.Node {
	return v.node_
}

func (v arrayValidator) parent() validator {
	return v.parent_
}

func (v *arrayValidator) setParent(parent validator) {
	v.parent_ = parent
}

// return array (pointers to validators, or nil if not found) and bool (true if validator is done)
func (v *arrayValidator) feed(jsonLexeme lexeme.LexEvent) ([]validator, bool) {
	defer lexeme.CatchLexEventError(jsonLexeme)

	switch jsonLexeme.Type() { //nolint:exhaustive // We will throw a panic in over cases.
	case lexeme.ArrayBegin, lexeme.ArrayItemEnd:
		return nil, false

	case lexeme.ArrayItemBegin:
		if arrayNode, ok := v.node_.(*schema.ArrayNode); ok {
			childNode := arrayNode.Child(v.itemsCounter) // can panic
			v.itemsCounter++
			return NodeValidatorList(childNode, v.rootSchema, v), false
		} else { // mixed node
			panic(errors.ErrElementNotFoundInArray)
		}

	case lexeme.ArrayEnd:
		if arrayNode, ok := v.node_.(*schema.ArrayNode); ok {
			arrayNode.ConstraintMap().EachSafe(func(_ constraint.Type, av constraint.Constraint) {
				if arrayValidator, ok := av.(constraint.ArrayValidator); ok {
					arrayValidator.ValidateTheArray(v.itemsCounter)
				}
			})
		}
		return nil, true
	}

	panic(errors.ErrUnexpectedLexInArrayValidator)
}

func (v arrayValidator) log() string {
	return fmt.Sprintf("array [%p]", v.node_)
}

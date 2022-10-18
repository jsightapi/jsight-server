package checker

import (
	"github.com/jsightapi/jsight-schema-go-library/errors"
	"github.com/jsightapi/jsight-schema-go-library/internal/lexeme"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema/internal/schema"
)

type nodeChecker interface {
	Check(lexeme.LexEvent) errors.Error
}

func newNodeChecker(node schema.Node) (nodeChecker, error) {
	switch node.(type) {
	case *schema.LiteralNode:
		return newLiteralChecker(node), nil

	case *schema.ObjectNode:
		return newObjectChecker(), nil

	case *schema.ArrayNode:
		return newArrayChecker(), nil

	case *schema.MixedNode:
		return newMixedChecker(node), nil

	default:
		return nil, errors.ErrImpossible
	}
}

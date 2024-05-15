package checker

import (
	"github.com/jsightapi/jsight-schema-core/errs"
	"github.com/jsightapi/jsight-schema-core/kit"
	"github.com/jsightapi/jsight-schema-core/lexeme"
	"github.com/jsightapi/jsight-schema-core/notations/jschema/ischema"
)

type nodeChecker interface {
	Check(lexeme.LexEvent) kit.Error
}

func newNodeChecker(node ischema.Node) (nodeChecker, error) {
	switch node.(type) {
	case *ischema.LiteralNode:
		return newLiteralChecker(node), nil

	case *ischema.ObjectNode:
		return newObjectChecker(), nil

	case *ischema.ArrayNode:
		return newArrayChecker(), nil

	case *ischema.MixedNode:
		return newMixedChecker(node), nil

	default:
		return nil, errs.ErrRuntimeFailure.F()
	}
}

package checker

import (
	"github.com/jsightapi/jsight-schema-core/errs"
	"github.com/jsightapi/jsight-schema-core/kit"
	"github.com/jsightapi/jsight-schema-core/lexeme"
	"github.com/jsightapi/jsight-schema-core/notations/jschema/ischema"
)

type literalChecker struct {
	node ischema.Node
}

func newLiteralChecker(node ischema.Node) literalChecker {
	return literalChecker{
		node: node,
	}
}

func (c literalChecker) Check(nodeLex lexeme.LexEvent) (err kit.Error) {
	defer func() {
		if r := recover(); r != nil {
			err = lexeme.ConvertError(nodeLex, r)
		}
	}()

	if nodeLex.Type() != lexeme.LiteralEnd {
		return lexeme.NewError(nodeLex, errs.ErrChecker.F())
	}

	ValidateLiteralValue(c.node, nodeLex.Value()) // can panic

	return nil
}

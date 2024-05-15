package checker

import (
	"github.com/jsightapi/jsight-schema-core/kit"
	"github.com/jsightapi/jsight-schema-core/lexeme"
	"github.com/jsightapi/jsight-schema-core/notations/jschema/ischema"
)

type mixedChecker struct {
	node ischema.Node
}

func newMixedChecker(node ischema.Node) mixedChecker {
	return mixedChecker{
		node: node,
	}
}

func (c mixedChecker) Check(nodeLex lexeme.LexEvent) (err kit.Error) {
	defer func() {
		if r := recover(); r != nil {
			err = lexeme.ConvertError(nodeLex, r)
		}
	}()

	if nodeLex.Type() == lexeme.LiteralEnd {
		ValidateLiteralValue(c.node, nodeLex.Value()) // can panic
	}

	return nil
}

package checker

import (
	"github.com/jsightapi/jsight-schema-go-library/errors"
	"github.com/jsightapi/jsight-schema-go-library/internal/lexeme"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema/internal/schema"
)

type objectChecker struct {
	node schema.Node
}

func newObjectChecker(node schema.Node) objectChecker {
	return objectChecker{
		node: node,
	}
}

func (objectChecker) check(nodeLex lexeme.LexEvent) (err errors.Error) {
	if nodeLex.Type() != lexeme.ObjectEnd {
		return lexeme.NewLexEventError(nodeLex, errors.ErrChecker)
	}

	return nil
}

func (c objectChecker) indentedString(depth int) string {
	return c.node.IndentedNodeString(depth)
}

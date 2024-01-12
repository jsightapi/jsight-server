package ischema

import (
	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/errs"
	"github.com/jsightapi/jsight-schema-core/json"
	"github.com/jsightapi/jsight-schema-core/lexeme"
)

type LiteralNode struct {
	baseNode
}

var _ Node = &LiteralNode{}

func newLiteralNode(lex lexeme.LexEvent) *LiteralNode {
	n := LiteralNode{
		baseNode: newBaseNode(lex),
	}
	return &n
}

func (n *LiteralNode) Grow(lex lexeme.LexEvent) (Node, bool) {
	switch lex.Type() {
	case lexeme.LiteralBegin:

	case lexeme.LiteralEnd:
		n.schemaLexEvent = lex
		t := json.Guess(lex.Value()).LiteralJsonType()
		n.setJsonType(t)
		return n.parent, false

	default:
		panic(errs.ErrUnexpectedLexicalEvent.F(lex.Type().String(), "in literal node"))
	}

	return n, false
}

func (n *LiteralNode) ASTNode() (schema.ASTNode, error) {
	an := astNodeFromNode(n)
	an.Value = n.Value().Unquote().String()
	return an, nil
}

func (n *LiteralNode) Copy() Node {
	nn := *n
	nn.baseNode = n.baseNode.Copy()
	return &nn
}

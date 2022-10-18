package schema

import (
	jschema "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/internal/json"
	"github.com/jsightapi/jsight-schema-go-library/internal/lexeme"
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
		panic(`Unexpected lexical event "` + lex.Type().String() + `" in literal node`)
	}

	return n, false
}

func (n *LiteralNode) ASTNode() (jschema.ASTNode, error) {
	an := astNodeFromNode(n)
	an.Value = n.Value().Unquote().String()
	return an, nil
}

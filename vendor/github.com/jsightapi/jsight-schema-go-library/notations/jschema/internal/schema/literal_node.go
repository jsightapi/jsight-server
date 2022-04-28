package schema

import (
	"strings"

	jschema "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/internal/json"
	"github.com/jsightapi/jsight-schema-go-library/internal/lexeme"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema/internal/schema/constraint"
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

func (n LiteralNode) IndentedTreeString(depth int) string {
	return n.IndentedNodeString(depth)
}

func (n LiteralNode) IndentedNodeString(depth int) string {
	indent := strings.Repeat("\t", depth)

	var str strings.Builder
	str.WriteString(indent + "* " + n.Type().String() + "\n")

	n.constraints.EachSafe(func(_ constraint.Type, v constraint.Constraint) {
		str.WriteString(indent + "* " + v.String() + "\n")
	})

	return str.String()
}

func (n *LiteralNode) ASTNode() (jschema.ASTNode, error) {
	an := astNodeFromNode(n)
	an.Value = n.Value().Unquote().String()
	return an, nil
}

package schema

import (
	jschema "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/errors"
	"github.com/jsightapi/jsight-schema-go-library/internal/json"
	"github.com/jsightapi/jsight-schema-go-library/internal/lexeme"
)

type ArrayNode struct {
	// children a children node list.
	children []Node

	baseNode

	// waitingForChild indicates that we should add children.
	// The Grow method will create a child node by getting the next lexical event.
	waitingForChild bool
}

var _ Node = &ArrayNode{}

func newArrayNode(lex lexeme.LexEvent) *ArrayNode {
	n := ArrayNode{
		baseNode: newBaseNode(lex),
		children: make([]Node, 0, 10),
	}
	n.setJsonType(json.TypeArray)
	return &n
}

func (n *ArrayNode) Grow(lex lexeme.LexEvent) (Node, bool) {
	if n.waitingForChild {
		n.waitingForChild = false
		child := NewNode(lex)
		n.addChild(child)
		return child, true
	}

	switch lex.Type() {
	case lexeme.ArrayBegin, lexeme.ArrayItemEnd:

	case lexeme.ArrayItemBegin:
		n.waitingForChild = true

	case lexeme.ArrayEnd:
		return n.parent, false

	default:
		panic(`Unexpected lexical event "` + lex.Type().String() + `" in array node`)
	}

	return n, false
}

func (n *ArrayNode) addChild(child Node) {
	child.SetParent(n)
	n.children = append(n.children, child)
}

func (n ArrayNode) Children() []Node {
	return n.children
}

func (n ArrayNode) Len() int {
	return len(n.children)
}

func (n ArrayNode) Child(i uint) Node {
	length := uint(len(n.children))
	if length == 0 {
		panic(errors.ErrElementNotFoundInArray)
	} else if i >= length {
		i = length - 1
	}
	return n.children[i]
}

func (n *ArrayNode) ASTNode() (jschema.ASTNode, error) {
	an := astNodeFromNode(n)
	l := len(n.children)

	if l > 0 {
		an.Children = make([]jschema.ASTNode, 0, l)
	}

	for _, c := range n.children {
		cn, err := c.ASTNode()
		if err != nil {
			return jschema.ASTNode{}, err
		}
		an.Children = append(an.Children, cn)
	}

	return an, nil
}

package ischema

import (
	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/errs"

	"github.com/jsightapi/jsight-schema-core/json"
	"github.com/jsightapi/jsight-schema-core/lexeme"
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
		panic(errs.ErrUnexpectedLexicalEvent.F(lex.Type().String(), "in array node"))
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
		panic(errs.ErrElementNotFoundInArray.F())
	} else if i >= length {
		i = length - 1
	}
	return n.children[i]
}

func (n *ArrayNode) ASTNode() (schema.ASTNode, error) {
	an := astNodeFromNode(n)
	l := len(n.children)

	if l > 0 {
		an.Children = make([]schema.ASTNode, 0, l)
	}

	for _, c := range n.children {
		cn, err := c.ASTNode()
		if err != nil {
			return schema.ASTNode{}, err
		}
		an.Children = append(an.Children, cn)
	}

	return an, nil
}

func (n *ArrayNode) Copy() Node {
	nn := *n
	nn.baseNode = n.baseNode.Copy()

	nn.children = make([]Node, 0, len(n.children))
	for _, v := range n.children {
		nn.children = append(nn.children, v.Copy())
	}

	return &nn
}

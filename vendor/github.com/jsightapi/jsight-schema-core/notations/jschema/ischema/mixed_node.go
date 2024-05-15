package ischema

import (
	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/errs"
	"github.com/jsightapi/jsight-schema-core/json"
	"github.com/jsightapi/jsight-schema-core/lexeme"
)

type MixedNode struct {
	baseNode
}

var _ Node = &MixedNode{}

func NewMixedNode(lex lexeme.LexEvent) *MixedNode {
	n := MixedNode{
		baseNode: newBaseNode(lex),
	}
	n.setJsonType(json.Guess(lex.Value()).JsonType())
	return &n
}

func (*MixedNode) SetRealType(string) bool {
	// Mixed value node is always have mixed type.
	return true
}

// SetJsonType for mixed node n.baseNode.jsonType is an EXAMPLE type
func (n *MixedNode) SetJsonType(t json.Type) {
	n.setJsonType(t)
}

func (*MixedNode) Grow(lexeme.LexEvent) (Node, bool) {
	panic(errs.ErrNodeGrow.F())
}

func (n MixedNode) ASTNode() (schema.ASTNode, error) {
	an := newASTNode()

	an.SchemaType = n.Type().String()
	an.Value = n.Value().Unquote().String()
	an.Rules = collectASTRules(n.constraints)
	an.Comment = n.comment

	return an, nil
}

func (n *MixedNode) Copy() Node {
	nn := *n
	nn.baseNode = n.baseNode.Copy()
	return &nn
}

package schema

import (
	jschema "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/errors"
	"github.com/jsightapi/jsight-schema-go-library/internal/json"
	"github.com/jsightapi/jsight-schema-go-library/internal/lexeme"
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

// SetJsonType for mixed node n.baseNode.jsonType is an EXAMPLE type
func (n *MixedNode) SetJsonType(t json.Type) {
	n.setJsonType(t)
}

func (*MixedNode) Grow(lexeme.LexEvent) (Node, bool) {
	panic(errors.ErrNodeGrow)
}

func (n MixedNode) ASTNode() (jschema.ASTNode, error) {
	an := newASTNode()

	an.SchemaType = n.Type().String()
	an.Value = n.Value().Unquote().String()
	an.Rules = collectASTRules(n.constraints)
	an.Comment = n.comment

	return an, nil
}

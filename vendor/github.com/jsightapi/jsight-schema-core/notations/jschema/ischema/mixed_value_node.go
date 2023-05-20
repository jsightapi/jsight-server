package ischema

import (
	"strings"

	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/bytes"
	"github.com/jsightapi/jsight-schema-core/errs"
	"github.com/jsightapi/jsight-schema-core/json"
	"github.com/jsightapi/jsight-schema-core/lexeme"
	"github.com/jsightapi/jsight-schema-core/notations/jschema/ischema/constraint"
)

type MixedValueNode struct {
	schemaType string
	value      string

	types []string

	baseNode
}

var _ Node = (*MixedValueNode)(nil)

func NewMixedValueNode(lex lexeme.LexEvent) *MixedValueNode {
	n := MixedValueNode{
		baseNode: newBaseNode(lex),
	}
	n.setJsonType(json.TypeMixed)
	n.realType = json.TypeMixed.String()
	return &n
}

func (*MixedValueNode) SetRealType(string) bool {
	// Mixed value node is always have mixed type.
	return true
}

func (n *MixedValueNode) AddConstraint(c constraint.Constraint) {
	switch t := c.(type) {
	case *constraint.TypeConstraint:
		n.addTypeConstraint(t)
		n.types = []string{t.Bytes().String()}

	case *constraint.Or:
		n.addOrConstraint(t)

	case *constraint.TypesList:
		n.types = t.Names()
		n.baseNode.AddConstraint(t)

	default:
		n.baseNode.AddConstraint(t)
	}
}

func (n *MixedValueNode) addTypeConstraint(c *constraint.TypeConstraint) {
	exists, ok := n.constraints.Get(constraint.TypeConstraintType)
	if !ok {
		n.baseNode.AddConstraint(c)
		n.schemaType = c.Bytes().Unquote().String()
		return
	}

	newVal := c.Bytes().Unquote().String()
	existsVal := exists.(constraint.BytesKeeper).Bytes().Unquote().String()
	if newVal != existsVal && newVal != "mixed" {
		panic(errs.ErrDuplicateRule.F(c.Type().String()))
	}
	n.constraints.Set(c.Type(), c)
	n.schemaType = "mixed"
}

func (n *MixedValueNode) addOrConstraint(c *constraint.Or) {
	if tc, ok := n.constraints.Get(constraint.TypeConstraintType); ok {
		n.addTypeConstraint(constraint.NewType(
			bytes.NewBytes(`"mixed"`),
			tc.(*constraint.TypeConstraint).Source(),
		))
	}
	n.baseNode.AddConstraint(c)
}

func (n *MixedValueNode) Grow(lex lexeme.LexEvent) (Node, bool) {
	switch lex.Type() {
	case lexeme.MixedValueBegin:

	case lexeme.MixedValueEnd:
		n.schemaLexEvent = lex
		n.value = lex.Value().TrimSpaces().String()
		n.schemaType = n.value
		return n.parent, false

	default:
		panic(errs.ErrUnexpectedLexicalEvent.F(lex.Type().String(), "in mixed value node"))
	}

	return n, false
}

func (n *MixedValueNode) ASTNode() (schema.ASTNode, error) {
	an := astNodeFromNode(n)

	an.SchemaType = n.schemaType
	if strings.ContainsRune(n.value, '|') {
		an.SchemaType = json.TypeMixed.String()
	}
	an.Value = n.value
	return an, nil
}

func (n *MixedValueNode) GetTypes() []string {
	return n.types
}

func (n *MixedValueNode) Copy() Node {
	nn := *n
	nn.baseNode = n.baseNode.Copy()
	return &nn
}

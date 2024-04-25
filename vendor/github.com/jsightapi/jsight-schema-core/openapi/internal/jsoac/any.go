package jsoac

import (
	schema "github.com/jsightapi/jsight-schema-core"
)

type Any struct {
	Example     *Example     `json:"example,omitempty"`
	Nullable    *Nullable    `json:"nullable,omitempty"`
	Description *Description `json:"description,omitempty"`
}

var _ Node = (*Any)(nil)

func newAny(astNode schema.ASTNode) *Any {
	a := Any{
		Nullable:    newNullable(astNode),
		Description: newDescription(astNode),
	}

	switch astNode.TokenType {
	case schema.TokenTypeString:
		if astNode.Value == "" {
			a.Example = nil
		} else {
			a.Example = newExample(astNode.Value, true)
		}
	case schema.TokenTypeNumber, schema.TokenTypeBoolean, schema.TokenTypeNull:
		a.Example = newExample(astNode.Value, false)
	default:
		a.Example = nil
	}

	return &a
}

func (n *Any) SetNodeDescription(s string) {
	n.Description = newDescriptionFromString(s)
}

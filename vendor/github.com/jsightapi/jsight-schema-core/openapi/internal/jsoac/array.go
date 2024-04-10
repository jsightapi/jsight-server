package jsoac

import (
	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/openapi/internal"
)

type Array struct {
	OADType     OADType      `json:"type"`
	Items       ArrayItems   `json:"items"`
	MinItems    *int64       `json:"minItems,omitempty"`
	MaxItems    *int64       `json:"maxItems,omitempty"`
	Nullable    *Nullable    `json:"nullable,omitempty"`
	Description *Description `json:"description,omitempty"`
}

var _ Node = (*Array)(nil)

func newArray(astNode schema.ASTNode) *Array {
	maxItems := newMaxItems(astNode)
	if len(astNode.Children) == 0 {
		maxItems = internal.Int64Ref(0)
	}
	a := Array{
		OADType:     OADTypeArray,
		Items:       newArrayItems(len(astNode.Children)),
		MinItems:    newMinItems(astNode),
		MaxItems:    maxItems,
		Nullable:    newNullable(astNode),
		Description: newDescription(astNode),
	}
	for _, an := range astNode.Children {
		a.appendItem(an)
	}
	return &a
}

func (a *Array) appendItem(astNode schema.ASTNode) {
	a.Items.append(newNode(astNode))
}

func (a *Array) SetNodeDescription(s string) {
	a.Description = newDescriptionFromString(s)
}

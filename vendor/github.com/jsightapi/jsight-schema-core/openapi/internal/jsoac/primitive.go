package jsoac

import (
	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/openapi/internal"
)

type Primitive struct {
	OADType          *OADType     `json:"type,omitempty"`
	Example          *Example     `json:"example,omitempty"`
	Pattern          *Pattern     `json:"pattern,omitempty"`
	Format           *string      `json:"format,omitempty"`
	Enum             *Enum        `json:"enum,omitempty"`
	Minimum          *Number      `json:"minimum,omitempty"`
	Maximum          *Number      `json:"maximum,omitempty"`
	ExclusiveMinimum *bool        `json:"exclusiveMinimum,omitempty"`
	ExclusiveMaximum *bool        `json:"exclusiveMaximum,omitempty"`
	MinLength        *int64       `json:"minLength,omitempty"`
	MaxLength        *int64       `json:"maxLength,omitempty"`
	MultipleOf       *float64     `json:"multipleOf,omitempty"`
	Nullable         *Nullable    `json:"nullable,omitempty"`
	Description      *Description `json:"description,omitempty"`
}

var _ Node = (*Primitive)(nil)

func oadType(schemaType string, t OADType) *OADType {
	if schemaType == internal.StringEnum {
		return nil
	}
	return &t
}

func newPrimitive(astNode schema.ASTNode) Node {
	if astNode.Rules.Has("or") {
		return newOr(astNode)
	}

	if rule, ok := astNode.Rules.Get("type"); ok && rule.TokenType == schema.TokenTypeShortcut {
		ref := newRef(astNode)
		return ref
	}

	t := oadTypeFromASTNode(astNode)
	var p = Primitive{
		OADType:          oadType(astNode.SchemaType, t),
		Example:          newExample(astNode.Value, t == OADTypeString),
		Pattern:          newPattern(astNode),
		Format:           newFormat(astNode),
		Enum:             newEnum(astNode),
		Minimum:          newMinimum(astNode),
		Maximum:          newMaximum(astNode),
		ExclusiveMinimum: newExclusiveMinimum(astNode),
		ExclusiveMaximum: newExclusiveMaximum(astNode),
		MinLength:        newMinLength(astNode),
		MaxLength:        newMaxLength(astNode),
		MultipleOf:       newMultipleOf(astNode),
		Nullable:         newNullable(astNode),
		Description:      newDescription(astNode),
	}
	return &p
}

func (p *Primitive) SetNodeDescription(s string) {
	p.Description = newDescriptionFromString(s)
}

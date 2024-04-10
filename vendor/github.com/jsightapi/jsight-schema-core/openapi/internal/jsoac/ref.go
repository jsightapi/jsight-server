package jsoac

import (
	"encoding/json"

	"github.com/jsightapi/jsight-schema-core/openapi/internal"

	schema "github.com/jsightapi/jsight-schema-core"
)

type Ref struct {
	UserType    UserTypeArray `json:"allOf"`
	Example     *Example      `json:"example,omitempty"`
	Nullable    *Nullable     `json:"nullable,omitempty"`
	Description *Description  `json:"description,omitempty"`
}

var _ Node = (*Ref)(nil)

func newRef(astNode schema.ASTNode) Node {
	if astNode.Rules.Has("or") {
		return newOr(astNode)
	}

	return &Ref{
		UserType:    UserTypeArray{UserType{name: astNode.SchemaType}},
		Nullable:    newNullableFromBool(internal.IsNullable(astNode)),
		Example:     refExample(astNode),
		Description: newDescription(astNode),
	}
}

func newRefFromUserTypeName(name string, nullable bool) *Ref {
	return &Ref{
		UserType: UserTypeArray{UserType{name: name}},
		Nullable: newNullableFromBool(nullable),
	}
}

func refExample(astNode schema.ASTNode) *Example {
	if astNode.TokenType == schema.TokenTypeShortcut {
		return nil
	}
	return newExample(astNode.Value, internal.IsString(astNode))
}

func (r *Ref) SetNodeDescription(s string) {
	r.Description = newDescriptionFromString(s)
}

func (r Ref) MarshalJSON() ([]byte, error) {
	if r.Example == nil && r.Nullable == nil && r.Description == nil {
		return r.UserType.UserType.MarshalJSON()
	}

	type Alias Ref
	var data = struct {
		Alias
	}{
		Alias: Alias(r),
	}

	b, err := json.Marshal(data)
	return b, err
}

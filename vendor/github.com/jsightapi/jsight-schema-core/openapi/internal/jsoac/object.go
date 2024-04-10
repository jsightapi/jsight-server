package jsoac

import (
	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/openapi/internal"
)

type Object struct {
	cap                  int
	OADType              OADType               `json:"type"`
	Properties           ObjectProperties      `json:"properties"`
	Required             []string              `json:"required,omitempty"`
	AllOf                *AllOf                `json:"allOf,omitempty"`
	AdditionalProperties *AdditionalProperties `json:"additionalProperties,omitempty"`
	Nullable             *Nullable             `json:"nullable,omitempty"`
	Description          *Description          `json:"description,omitempty"`
}

var _ Node = (*Object)(nil)

func newObject(astNode schema.ASTNode) *Object {
	o := Object{
		cap:                  len(astNode.Children),
		OADType:              OADTypeObject,
		Properties:           newObjectProperties(len(astNode.Children)),
		AdditionalProperties: newAdditionalProperties(astNode),
		Required:             nil,
		AllOf:                newAllOf(astNode),
		Nullable:             newNullable(astNode),
		Description:          newDescription(astNode),
	}

	for _, an := range astNode.Children {
		if !an.IsKeyShortcut {
			o.appendProperty(an)
		}
	}

	return &o
}

func (o *Object) appendProperty(astNode schema.ASTNode) {
	key := astNode.Key
	value := newNode(astNode)

	o.Properties.append(key, value)

	if !astNode.Rules.Has("optional") || astNode.Rules.GetValue("optional").Value == internal.StringFalse {
		o.appendToRequired(key)
	}
}

func (o *Object) appendToRequired(key string) {
	o.initRequiredIfNecessary()
	o.Required = append(o.Required, key)
}

func (o *Object) initRequiredIfNecessary() {
	if o.Required == nil {
		o.Required = make([]string, 0, o.cap)
	}
}

func (o *Object) SetNodeDescription(s string) {
	o.Description = newDescriptionFromString(s)
}

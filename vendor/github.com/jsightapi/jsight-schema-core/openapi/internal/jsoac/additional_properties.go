package jsoac

import (
	"encoding/json"

	"github.com/jsightapi/jsight-schema-core/openapi/internal"

	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/errs"
)

type additionalPropertiesMode int

const (
	additionalPropertiesNull additionalPropertiesMode = iota
	additionalPropertiesFalse
	additionalPropertiesArray
	additionalPropertiesPrimitive
	additionalPropertiesFormat
	additionalPropertiesUserType
)

type AdditionalProperties struct {
	mode         additionalPropertiesMode
	oadType      *OADType
	format       string
	userTypeName string
}

var _ json.Marshaler = AdditionalProperties{}
var _ json.Marshaler = &AdditionalProperties{}

func newAdditionalProperties(astNode schema.ASTNode) *AdditionalProperties {
	if astNode.Rules.Has("additionalProperties") {
		r := astNode.Rules.GetValue("additionalProperties")
		switch r.TokenType {
		case schema.TokenTypeBoolean:
			return newBooleanAdditionalProperties(r)
		case schema.TokenTypeString:
			return newStringAdditionalProperties(r)
		default:
			panic(errs.ErrRuntimeFailure.F())
		}
	}

	// The additionalProperties JSight rule is missing
	return newFalseAdditionalProperties()
}

func newStringAdditionalProperties(r schema.RuleASTNode) *AdditionalProperties {
	if r.Value == stringNull {
		return &AdditionalProperties{mode: additionalPropertiesNull}
	}

	if r.Value == stringArray {
		return &AdditionalProperties{mode: additionalPropertiesArray}
	}

	if r.Value == stringAny {
		return nil
	}

	if r.Value[0] == '@' {
		return &AdditionalProperties{mode: additionalPropertiesUserType, userTypeName: r.Value}
	}

	t := oadTypeFromSchemaType(r.Value)
	f := internal.FormatFromSchemaType(r.Value)

	if f == nil {
		return &AdditionalProperties{
			mode:    additionalPropertiesPrimitive,
			oadType: &t,
		}
	}

	return &AdditionalProperties{
		mode:    additionalPropertiesFormat,
		oadType: &t,
		format:  *f,
	}
}

func newBooleanAdditionalProperties(r schema.RuleASTNode) *AdditionalProperties {
	if r.Value == stringFalse {
		return newFalseAdditionalProperties()
	}
	return nil // JSight additionalProperties: true
}

func newFalseAdditionalProperties() *AdditionalProperties {
	return &AdditionalProperties{
		mode: additionalPropertiesFalse,
	}
}

func (a AdditionalProperties) MarshalJSON() ([]byte, error) {
	switch a.mode {
	case additionalPropertiesFalse:
		return a.booleanJSON()
	case additionalPropertiesNull:
		return a.nullJSON()
	case additionalPropertiesArray:
		return a.arrayJSON()
	case additionalPropertiesFormat:
		return a.formatJSON()
	case additionalPropertiesPrimitive:
		return a.primitiveJSON()
	case additionalPropertiesUserType:
		return a.userTypeJSON()
	default:
		panic(errs.ErrRuntimeFailure.F())
	}
}

func (a AdditionalProperties) arrayJSON() ([]byte, error) {
	data := struct {
		OADType OADType        `json:"type"`
		Items   map[string]any `json:"items"`
	}{
		OADType: OADTypeArray,
		Items:   map[string]any{},
	}
	return json.Marshal(data)
}

func (a AdditionalProperties) booleanJSON() ([]byte, error) {
	return []byte(stringFalse), nil
}

func (a AdditionalProperties) nullJSON() ([]byte, error) {
	return []byte(`{ "enum": [null] }`), nil
}

func (a AdditionalProperties) primitiveJSON() ([]byte, error) {
	data := struct {
		OADType OADType `json:"type"`
	}{
		OADType: *a.oadType,
	}
	return json.Marshal(data)
}

func (a AdditionalProperties) formatJSON() ([]byte, error) {
	data := struct {
		OADType OADType `json:"type"`
		Format  string  `json:"format"`
	}{
		OADType: *a.oadType,
		Format:  a.format,
	}
	return json.Marshal(data)
}

func (a AdditionalProperties) userTypeJSON() ([]byte, error) {
	ref := newRefFromUserTypeName(a.userTypeName, false)
	return ref.MarshalJSON()
}

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
	additionalPropertiesAnyOf
	additionalPropertiesObject
)

type AdditionalProperties struct {
	mode         additionalPropertiesMode
	oadType      *OADType
	format       string
	userTypeName string
	node         schema.ASTNode
}

var _ json.Marshaler = AdditionalProperties{}
var _ json.Marshaler = &AdditionalProperties{}

func newAdditionalProperties(astNode schema.ASTNode) *AdditionalProperties {
	if hasKeyShortcutChild(astNode) {
		return newKeyShortcutAdditionalProperties(astNode)
	}

	if ap, ok := astNode.Rules.Get(internal.StringAdditionalProperties); ok {
		return newBasicAdditionalProperties(ap)
	}

	// The additionalProperties JSight rule is missing
	return newFalseAdditionalProperties()
}

func newBasicAdditionalProperties(ap schema.RuleASTNode) *AdditionalProperties {
	switch ap.TokenType {
	case schema.TokenTypeBoolean:
		return newBooleanAdditionalProperties(ap)
	case schema.TokenTypeString:
		return newStringAdditionalProperties(ap)
	default:
		panic(errs.ErrRuntimeFailure.F())
	}
}

// check is astNode have children with some key as shortcut and with additional properties
func hasKeyShortcutChild(astNode schema.ASTNode) bool {
	for _, an := range astNode.Children {
		if an.IsKeyShortcut {
			return true
		}
	}
	return false
}

func newKeyShortcutAdditionalProperties(astNode schema.ASTNode) *AdditionalProperties {
	ap, hasAdditionalPropertiesRule := astNode.Rules.Get(internal.StringAdditionalProperties)

	for _, an := range astNode.Children {
		if an.IsKeyShortcut {
			if hasAdditionalPropertiesRule {
				if (ap.TokenType == schema.TokenTypeBoolean && ap.Value == internal.StringTrue) ||
					(ap.TokenType == schema.TokenTypeString && ap.Value == internal.StringAny) {
					return nil
				}
			}
			return newAnyOfAdditionalProperties(astNode)
		}
	}

	return nil
}

func newAnyOfAdditionalProperties(node schema.ASTNode) *AdditionalProperties {
	t := OADTypeObject
	return &AdditionalProperties{
		mode:    additionalPropertiesAnyOf,
		oadType: &t,
		node:    node,
	}
}

func newStringAdditionalProperties(r schema.RuleASTNode) *AdditionalProperties {
	if r.Value == internal.StringNull {
		return &AdditionalProperties{mode: additionalPropertiesNull}
	}

	if r.Value == internal.StringArray {
		return &AdditionalProperties{mode: additionalPropertiesArray}
	}

	if r.Value == internal.StringObject {
		return &AdditionalProperties{mode: additionalPropertiesObject}
	}

	if r.Value == internal.StringAny {
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
	if r.Value == internal.StringFalse {
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
	case additionalPropertiesObject:
		return a.objectJSON()
	case additionalPropertiesFormat:
		return a.formatJSON()
	case additionalPropertiesPrimitive:
		return a.primitiveJSON()
	case additionalPropertiesUserType:
		return a.userTypeJSON()
	case additionalPropertiesAnyOf:
		return a.anyOfJSON(a.node)
	default:
		panic(errs.ErrRuntimeFailure.F())
	}
}

func (a *AdditionalProperties) anyOfJSON(node schema.ASTNode) ([]byte, error) {
	var items []any

	if ap, ok := node.Rules.Get(internal.StringAdditionalProperties); ok {
		if ap.TokenType == schema.TokenTypeString {
			items = append(items, makeAdditionalAnyJSONObjects(ap))
		}
	}

	for _, astNode := range node.Children {
		if astNode.Key != "" && astNode.Key[0] == '@' {
			items = append(items, newNode(astNode))
		}
	}

	data := struct {
		Items []any `json:"anyOf"`
	}{
		Items: items,
	}
	m, err := json.Marshal(data)
	return m, err
}

func (a *AdditionalProperties) arrayJSON() ([]byte, error) {
	data := struct {
		OADType OADType        `json:"type"`
		Items   map[string]any `json:"items"`
	}{
		OADType: OADTypeArray,
		Items:   map[string]any{},
	}
	return json.Marshal(data)
}

func (a *AdditionalProperties) objectJSON() ([]byte, error) {
	data := struct {
		OADType              OADType        `json:"type"`
		Properties           map[string]any `json:"properties"`
		AdditionalProperties bool           `json:"additionalProperties"`
	}{
		OADType:              OADTypeObject,
		Properties:           map[string]any{},
		AdditionalProperties: false,
	}
	return json.Marshal(data)
}

func (a *AdditionalProperties) booleanJSON() ([]byte, error) {
	return []byte(internal.StringFalse), nil
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

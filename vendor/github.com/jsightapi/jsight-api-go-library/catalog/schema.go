package catalog

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"

	jschemaLib "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/bytes"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema"

	"github.com/jsightapi/jsight-api-go-library/notation"
)

// Schema represent a user defined schema.
type Schema struct {
	// Notation used notation for this schema.
	Notation notation.SchemaNotation

	// JSight notation specific fields.

	// ContentJSight a JSight schema.
	ContentJSight *SchemaContentJSight
	// UsedUserTypes a list of used user types.
	UsedUserTypes *StringSet
	// UserUserTypes a list of used user enums.
	UsedUserEnums *StringSet

	// Regexp notation specific fields.

	// ContentRegexp a regular expression.
	ContentRegexp string
}

func NewRegexSchema(regexStr bytes.Bytes) Schema {
	s := NewSchema(notation.SchemaNotationRegex)
	s.ContentRegexp = strings.TrimPrefix(regexStr.String(), "/")
	s.ContentRegexp = strings.TrimSuffix(s.ContentRegexp, "/")
	return s
}

func NewSchema(n notation.SchemaNotation) Schema {
	return Schema{
		Notation:      n,
		UsedUserTypes: &StringSet{},
		UsedUserEnums: &StringSet{},
	}
}

// StringSet a set of strings.
// gen:Set
type StringSet struct {
	data  map[string]struct{}
	order []string
	mx    sync.RWMutex
}

// UnmarshalSchema unmarshal a schema from the given slice of bytes.
// Marshaled schema shouldn't contain any trailing symbols.
func UnmarshalSchema(name string, b []byte, userTypes *UserSchemas) (_ Schema, err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			} else {
				panic(r)
			}
		}
	}()

	s := jschema.New(name, b)

	err = userTypes.Each(func(k string, v jschemaLib.Schema) error {
		return s.AddType(k, v)
	})
	if err != nil {
		return Schema{}, err
	}

	n, err := s.GetAST()
	if err != nil {
		return Schema{}, err
	}

	return unmarshalSchema(n), nil
}

func unmarshalSchema(n jschemaLib.ASTNode) Schema {
	s := NewSchema(notation.SchemaNotationJSight)
	s.ContentJSight = astNodeToJsightContent(n, s.UsedUserTypes, s.UsedUserEnums)
	return s
}

func astNodeToJsightContent(
	node jschemaLib.ASTNode,
	usedUserTypes, usedUserEnums *StringSet,
) *SchemaContentJSight {
	rules := collectJSightContentRules(node, usedUserTypes)

	var isOptional bool
	if r, ok := rules.Get("optional"); ok {
		var err error
		isOptional, err = strconv.ParseBool(r.ScalarValue)
		if err != nil {
			// Normally this shouldn't happen.
			panic(err)
		}
	}

	return &SchemaContentJSight{
		IsKeyShortcut: node.IsKeyShortcut,
		JsonType:      node.JSONType,
		Type:          node.SchemaType,
		Optional:      isOptional,
		ScalarValue:   node.Value,
		InheritedFrom: "", // Will be filled during catalog compilation.
		Note:          Annotation(node.Comment),
		Rules:         rules,
		Properties:    collectJSightContentProperties(node, usedUserTypes, usedUserEnums),
		Items:         collectJSightContentItems(node, usedUserTypes, usedUserEnums),
	}
}

func collectJSightContentRules(node jschemaLib.ASTNode, usedUserTypes *StringSet) *Rules {
	rr := &Rules{}

	if node.Rules.Len() == 0 {
		return rr
	}

	node.Rules.EachSafe(func(k string, v jschemaLib.RuleASTNode) {
		switch k {
		case "type":
			if v.Value[0] == '@' {
				usedUserTypes.Add(v.Value)
			}
			if v.Source == jschemaLib.RuleASTNodeSourceGenerated {
				return
			}

		case "allOf":
			if v.Value != "" {
				usedUserTypes.Add(v.Value)
			}

			for _, i := range v.Items {
				usedUserTypes.Add(i.Value)
			}

		case "or":
			for _, i := range v.Items {
				var userType string
				if i.Value != "" {
					userType = i.Value
				} else {
					v, ok := i.Properties.Get("type")
					if ok {
						userType = v.Value
					} else {
						userType = node.SchemaType
					}
				}

				if userType[0] != '@' {
					continue
				}

				usedUserTypes.Add(userType)
			}

			if v.Source == jschemaLib.RuleASTNodeSourceGenerated {
				return
			}
		}
		rr.Set(k, astNodeToSchemaRule(v))
	})

	return rr
}

func astNodeToSchemaRule(node jschemaLib.RuleASTNode) Rule {
	properties := &Rules{}
	if node.Properties.Len() > 0 {
		node.Properties.EachSafe(func(k string, v jschemaLib.RuleASTNode) {
			properties.Set(k, astNodeToSchemaRule(v))
		})
	}

	var items []Rule
	if len(node.Items) > 0 {
		items = make([]Rule, 0, len(node.Items))
		for _, n := range node.Items {
			items = append(items, astNodeToSchemaRule(n))
		}
	}

	return Rule{
		JsonType:    node.JSONType,
		ScalarValue: node.Value,
		Note:        node.Comment,
		Properties:  properties,
		Items:       items,
	}
}

func collectJSightContentProperties(
	node jschemaLib.ASTNode,
	usedUserTypes, usedUserEnums *StringSet,
) *Properties {
	pp := &Properties{}
	if node.Properties.Len() > 0 {
		node.Properties.EachSafe(func(k string, v jschemaLib.ASTNode) {
			pp.Set(k, astNodeToJsightContent(v, usedUserTypes, usedUserEnums))

			if v.IsKeyShortcut {
				usedUserTypes.Add(k)
			}
		})
	}
	return pp
}

func collectJSightContentItems(
	node jschemaLib.ASTNode,
	usedUserTypes, usedUserEnums *StringSet,
) []*SchemaContentJSight {
	var ii []*SchemaContentJSight
	if len(node.Items) > 0 {
		ii = make([]*SchemaContentJSight, 0, len(node.Items))
		for _, n := range node.Items {
			an := astNodeToJsightContent(n, usedUserTypes, usedUserEnums)
			an.Optional = true
			ii = append(ii, an)
		}
	}
	return ii
}

func (schema Schema) MarshalJSON() ([]byte, error) {
	data := struct {
		Content       interface{}             `json:"content,omitempty"`
		Example       string                  `json:"example,omitempty"`
		Notation      notation.SchemaNotation `json:"notation"`
		UsedUserTypes []string                `json:"usedUserTypes,omitempty"`
		UsedUserEnums []string                `json:"usedUserEnums,omitempty"`
	}{
		Notation: schema.Notation,
	}

	switch schema.Notation {
	case notation.SchemaNotationJSight:
		data.Content = schema.ContentJSight
		if schema.UsedUserTypes != nil && schema.UsedUserTypes.Len() > 0 {
			data.UsedUserTypes = schema.UsedUserTypes.Data()
		}
		if schema.UsedUserEnums != nil && schema.UsedUserEnums.Len() > 0 {
			data.UsedUserEnums = schema.UsedUserEnums.Data()
		}
		// TODO data.Example = ...

	case notation.SchemaNotationRegex:
		data.Content = schema.ContentRegexp

	case notation.SchemaNotationAny, notation.SchemaNotationEmpty:
		// nothing

	default:
		return []byte{}, fmt.Errorf(`invalid schema notation "%s"`, schema.Notation)
	}

	return json.Marshal(data)
}

type SchemaContentJSight struct {
	// JsonType a JSON type.
	JsonType string

	// Type a JSight type.
	Type string

	// ScalarValue contains scalar value from the example.
	// Make sense only for scalar types like string, integer, and etc.
	ScalarValue string

	// InheritedFrom a user defined type from which this property is inherited.
	InheritedFrom string

	// Note a user note.
	Note string

	// Rules a list of attached rules.
	Rules *Rules

	// Properties represent available object properties.
	// Make sense only when Type is "object".
	Properties *Properties

	// Items represent available array items.
	// Make sense only when Type is "array".
	Items []*SchemaContentJSight

	// IsKeyShortcut indicates that this is an object property which is described
	// by user defined type.
	IsKeyShortcut bool

	// Optional indicates that this schema item is option or not.
	Optional bool
}

var (
	_ json.Marshaler = SchemaContentJSight{}
	_ json.Marshaler = &SchemaContentJSight{}
)

func (c SchemaContentJSight) MarshalJSON() (b []byte, err error) {
	switch c.JsonType {
	case jschemaLib.JSONTypeObject:
		b, err = c.marshalJSONObject()

	case jschemaLib.JSONTypeArray:
		b, err = c.marshalJSONArray()

	default:
		b, err = c.marshalJSONLiteral()
	}
	return b, err
}

func (c SchemaContentJSight) marshalJSONObject() ([]byte, error) {
	var data struct {
		Rules         *Rules      `json:"rules,omitempty"`
		Properties    *Properties `json:"properties,omitempty"`
		JsonType      string      `json:"jsonType,omitempty"`
		Type          string      `json:"type,omitempty"`
		InheritedFrom string      `json:"inheritedFrom,omitempty"`
		Note          string      `json:"note,omitempty"`
		IsKeyShortcut bool        `json:"isKeyShortcut,omitempty"`
		Optional      bool        `json:"optional"`
	}

	data.IsKeyShortcut = c.IsKeyShortcut
	data.JsonType = c.JsonType
	data.Type = c.Type
	data.Optional = c.Optional
	data.InheritedFrom = c.InheritedFrom
	data.Note = c.Note
	if c.Rules != nil && c.Rules.Len() > 0 {
		data.Rules = c.Rules
	}
	if c.Properties != nil && c.Properties.Len() > 0 {
		data.Properties = c.Properties
	}

	return json.Marshal(data)
}

func (c SchemaContentJSight) marshalJSONArray() ([]byte, error) {
	var data struct {
		Rules         *Rules                 `json:"rules,omitempty"`
		JsonType      string                 `json:"jsonType,omitempty"`
		Type          string                 `json:"type,omitempty"`
		InheritedFrom string                 `json:"inheritedFrom,omitempty"`
		Note          string                 `json:"note,omitempty"`
		Items         []*SchemaContentJSight `json:"items,omitempty"`
		IsKeyShortcut bool                   `json:"isKeyShortcut,omitempty"`
		Optional      bool                   `json:"optional"`
	}

	data.IsKeyShortcut = c.IsKeyShortcut
	data.JsonType = c.JsonType
	data.Type = c.Type
	data.Optional = c.Optional
	data.InheritedFrom = c.InheritedFrom
	data.Note = c.Note
	if c.Rules != nil && c.Rules.Len() > 0 {
		data.Rules = c.Rules
	}
	data.Items = c.Items

	return json.Marshal(data)
}

func (c SchemaContentJSight) marshalJSONLiteral() ([]byte, error) {
	var data struct {
		Rules         *Rules `json:"rules,omitempty"`
		JsonType      string `json:"jsonType,omitempty"`
		Type          string `json:"type,omitempty"`
		ScalarValue   string `json:"scalarValue"`
		InheritedFrom string `json:"inheritedFrom,omitempty"`
		Note          string `json:"note,omitempty"`
		IsKeyShortcut bool   `json:"isKeyShortcut,omitempty"`
		Optional      bool   `json:"optional"`
	}

	data.IsKeyShortcut = c.IsKeyShortcut
	data.JsonType = c.JsonType
	data.Type = c.Type
	data.Optional = c.Optional
	data.ScalarValue = c.ScalarValue
	data.InheritedFrom = c.InheritedFrom
	data.Note = c.Note
	if c.Rules != nil && c.Rules.Len() > 0 {
		data.Rules = c.Rules
	}

	return json.Marshal(data)
}

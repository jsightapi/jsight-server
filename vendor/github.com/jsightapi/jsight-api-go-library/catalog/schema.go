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
func UnmarshalSchema(
	name string,
	b []byte,
	userTypes *UserSchemas,
	enumRules map[string]jschemaLib.Rule,
) (_ Schema, err error) {
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

	for n, v := range enumRules {
		if err := s.AddRule(n, v); err != nil {
			return Schema{}, err
		}
	}

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

	if rules.Len() == 0 {
		rules = nil
	}

	c := &SchemaContentJSight{
		IsKeyUserTypeRef: node.IsKeyShortcut,
		TokenType:        node.JSONType,
		Type:             node.SchemaType,
		Optional:         isOptional,
		ScalarValue:      node.Value,
		InheritedFrom:    "", // Will be filled during catalog compilation.
		Note:             Annotation(node.Comment),
		Rules:            rules,
	}

	switch node.JSONType {
	case jschemaLib.JSONTypeObject:
		c.collectJSightContentObjectProperties(node, usedUserTypes, usedUserEnums)
	case jschemaLib.JSONTypeArray:
		c.collectJSightContentArrayItems(node, usedUserTypes, usedUserEnums)
	}

	return c
}

func collectJSightContentRules(node jschemaLib.ASTNode, usedUserTypes *StringSet) *Rules {
	if node.Rules.Len() == 0 {
		return &Rules{}
	}

	rr := newRulesBuilder(node.Rules.Len())

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

	return rr.Rules()
}

func astNodeToSchemaRule(node jschemaLib.RuleASTNode) Rule {
	rr := newRulesBuilder(node.Properties.Len() + len(node.Items))

	if node.Properties.Len() != 0 {
		node.Properties.EachSafe(func(k string, v jschemaLib.RuleASTNode) {
			rr.Set(k, astNodeToSchemaRule(v))
		})
	}

	if len(node.Items) != 0 {
		for _, n := range node.Items {
			rr.Append(astNodeToSchemaRule(n))
		}
	}

	var children []Rule
	if rr.Rules().Len() != 0 {
		children = rr.Rules().data
	}

	return Rule{
		TokenType:   node.JSONType,
		ScalarValue: node.Value,
		Note:        node.Comment,
		Children:    children,
	}
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
	// Key is key of object element.
	Key string

	// TokenType a JSON type.
	TokenType string

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

	// Children represent available object properties or array items.
	Children []*SchemaContentJSight

	// IsKeyUserTypeRef indicates that this is an object property which is described
	// by user defined type.
	IsKeyUserTypeRef bool

	// Optional indicates that this schema item is option or not.
	Optional bool
}

var (
	_ json.Marshaler = SchemaContentJSight{}
	_ json.Marshaler = &SchemaContentJSight{}
)

func (c *SchemaContentJSight) IsObjectHaveProperty(k string) bool {
	return c.ObjectProperty(k) != nil
}

func (c *SchemaContentJSight) ObjectProperty(k string) *SchemaContentJSight {
	for _, v := range c.Children {
		if v.Key == k {
			return v
		}
	}
	return nil
}

func (c *SchemaContentJSight) Unshift(v *SchemaContentJSight) {
	c.Children = append([]*SchemaContentJSight{v}, c.Children...)
}

func (c *SchemaContentJSight) collectJSightContentObjectProperties(
	node jschemaLib.ASTNode,
	usedUserTypes, usedUserEnums *StringSet,
) {
	if len(node.Children) > 0 {
		if c.Children == nil {
			c.Children = make([]*SchemaContentJSight, 0, len(node.Children))
		}
		for _, v := range node.Children {
			an := astNodeToJsightContent(v, usedUserTypes, usedUserEnums)
			an.Key = v.Key

			c.Children = append(c.Children, an)

			if v.IsKeyShortcut {
				usedUserTypes.Add(v.Key)
			}
		}
	}
}

func (c *SchemaContentJSight) collectJSightContentArrayItems(
	node jschemaLib.ASTNode,
	usedUserTypes, usedUserEnums *StringSet,
) {
	if len(node.Children) > 0 {
		if c.Children == nil {
			c.Children = make([]*SchemaContentJSight, 0, len(node.Children))
		}
		for _, n := range node.Children {
			an := astNodeToJsightContent(n, usedUserTypes, usedUserEnums)
			an.Optional = true
			c.Children = append(c.Children, an)
		}
	}
}

func (c SchemaContentJSight) MarshalJSON() (b []byte, err error) {
	switch c.TokenType {
	case jschemaLib.JSONTypeObject:
		b, err = c.marshalJSONObject()

	case jschemaLib.JSONTypeArray:
		b, err = c.marshalJSONArray()

	default:
		b, err = c.marshalJSONLiteral()
	}
	return b, err
}

func (c SchemaContentJSight) marshalJSONObject() ([]byte, error) { //nolint:dupl
	var data struct {
		Rules            []Rule                 `json:"rules,omitempty"`
		Key              string                 `json:"key,omitempty"`
		TokenType        string                 `json:"tokenType,omitempty"`
		Type             string                 `json:"type,omitempty"`
		InheritedFrom    string                 `json:"inheritedFrom,omitempty"`
		Note             string                 `json:"note,omitempty"`
		Children         []*SchemaContentJSight `json:"children"`
		IsKeyUserTypeRef bool                   `json:"isKeyUserTypeRef,omitempty"`
		Optional         bool                   `json:"optional"`
	}

	data.Key = c.Key
	data.IsKeyUserTypeRef = c.IsKeyUserTypeRef
	data.TokenType = c.TokenType
	data.Type = c.Type
	data.Optional = c.Optional
	data.InheritedFrom = c.InheritedFrom
	data.Note = c.Note
	if c.Rules != nil && c.Rules.Len() != 0 {
		data.Rules = c.Rules.data
	}
	if len(c.Children) == 0 {
		data.Children = make([]*SchemaContentJSight, 0)
	} else {
		data.Children = c.Children
	}

	return json.Marshal(data)
}

func (c SchemaContentJSight) marshalJSONArray() ([]byte, error) { //nolint:dupl
	var data struct {
		Rules            []Rule                 `json:"rules,omitempty"`
		Key              string                 `json:"key,omitempty"`
		TokenType        string                 `json:"tokenType,omitempty"`
		Type             string                 `json:"type,omitempty"`
		InheritedFrom    string                 `json:"inheritedFrom,omitempty"`
		Note             string                 `json:"note,omitempty"`
		Children         []*SchemaContentJSight `json:"children"`
		IsKeyUserTypeRef bool                   `json:"isKeyUserTypeRef,omitempty"`
		Optional         bool                   `json:"optional"`
	}

	data.Key = c.Key
	data.IsKeyUserTypeRef = c.IsKeyUserTypeRef
	data.TokenType = c.TokenType
	data.Type = c.Type
	data.Optional = c.Optional
	data.InheritedFrom = c.InheritedFrom
	data.Note = c.Note
	if c.Rules != nil && c.Rules.Len() != 0 {
		data.Rules = c.Rules.data
	}
	if len(c.Children) == 0 {
		data.Children = make([]*SchemaContentJSight, 0)
	} else {
		data.Children = c.Children
	}

	return json.Marshal(data)
}

func (c SchemaContentJSight) marshalJSONLiteral() ([]byte, error) {
	var data struct {
		Note             string `json:"note,omitempty"`
		Key              string `json:"key,omitempty"`
		TokenType        string `json:"tokenType,omitempty"`
		Type             string `json:"type,omitempty"`
		ScalarValue      string `json:"scalarValue"`
		InheritedFrom    string `json:"inheritedFrom,omitempty"`
		Rules            []Rule `json:"rules,omitempty"`
		IsKeyUserTypeRef bool   `json:"isKeyUserTypeRef,omitempty"`
		Optional         bool   `json:"optional"`
	}

	data.Key = c.Key
	data.IsKeyUserTypeRef = c.IsKeyUserTypeRef
	data.TokenType = c.TokenType
	data.Type = c.Type
	data.Optional = c.Optional
	data.ScalarValue = c.ScalarValue
	data.InheritedFrom = c.InheritedFrom
	data.Note = c.Note
	if c.Rules != nil && c.Rules.Len() != 0 {
		data.Rules = c.Rules.data
	}

	return json.Marshal(data)
}

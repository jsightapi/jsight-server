package catalog

import (
	"encoding/json"
	"strconv"

	jschemaLib "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema"

	"github.com/jsightapi/jsight-api-go-library/notation"
)

// UnmarshalJSightSchema unmarshal a schema from the given slice of bytes.
// Marshaled schema shouldn't contain any trailing symbols.
func UnmarshalJSightSchema(
	name string,
	b []byte,
	userTypes *UserSchemas,
	enumRules map[string]jschemaLib.Rule,
) (schema Schema, err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			} else {
				panic(r)
			}
		}
	}()

	s, err := prepareJSightSchema(name, b, userTypes, enumRules)
	if err != nil {
		return Schema{}, err
	}

	return unmarshalJSightSchema(s)
}

func prepareJSightSchema(
	name string,
	b []byte,
	userTypes *UserSchemas,
	enumRules map[string]jschemaLib.Rule,
) (jschemaLib.Schema, error) {
	s := jschema.New(name, b)

	for n, v := range enumRules {
		if err := s.AddRule(n, v); err != nil {
			return nil, err
		}
	}

	err := userTypes.Each(func(k string, v jschemaLib.Schema) error {
		return s.AddType(k, v)
	})
	if err != nil {
		return nil, err
	}
	return s, nil
}

func unmarshalJSightSchema(s jschemaLib.Schema) (Schema, error) {
	n, err := s.GetAST()
	if err != nil {
		return Schema{}, err
	}

	example, err := s.Example()
	if err != nil {
		return Schema{}, err
	}

	ret := NewSchema(notation.SchemaNotationJSight)
	ret.ContentJSight = astNodeToJsightContent(n, ret.UsedUserTypes, ret.UsedUserEnums)
	ret.Example = string(example)
	return ret, nil
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

type SchemaContentJSight struct {
	// Key is key of object element.
	Key *string

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
		if *(v.Key) == k {
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
			an.Key = SrtPtr(v.Key)

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
	case jschemaLib.JSONTypeObject, jschemaLib.JSONTypeArray:
		b, err = c.marshalJSONObjectOrArray()

	default:
		b, err = c.marshalJSONLiteral()
	}
	return b, err
}

func (c SchemaContentJSight) marshalJSONObjectOrArray() ([]byte, error) {
	var data struct {
		Rules            []Rule                 `json:"rules,omitempty"`
		Key              *string                `json:"key,omitempty"`
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
		Note             string  `json:"note,omitempty"`
		Key              *string `json:"key,omitempty"`
		TokenType        string  `json:"tokenType,omitempty"`
		Type             string  `json:"type,omitempty"`
		ScalarValue      string  `json:"scalarValue"`
		InheritedFrom    string  `json:"inheritedFrom,omitempty"`
		Rules            []Rule  `json:"rules,omitempty"`
		IsKeyUserTypeRef bool    `json:"isKeyUserTypeRef,omitempty"`
		Optional         bool    `json:"optional"`
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

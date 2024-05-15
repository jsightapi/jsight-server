package catalog

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/bytes"

	"github.com/jsightapi/jsight-api-core/jerr"
)

type ExchangeContent struct {
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
	Children []*ExchangeContent

	// IsKeyUserTypeRef indicates that this is an object property which is described
	// by user defined type.
	IsKeyUserTypeRef bool

	// Optional indicates that this schema item is option or not.
	Optional bool
}

var _ json.Marshaler = &ExchangeContent{}

func (c *ExchangeContent) processAllOf(uut *StringSet, catalogUserTypes *UserTypes) error {
	if c.TokenType != schema.TokenTypeObject {
		return nil
	}

	for _, cc := range c.Children {
		if err := cc.processAllOf(uut, catalogUserTypes); err != nil {
			return err
		}
	}

	rule, ok := c.Rules.Get("allOf")
	if !ok {
		return nil
	}

	switch rule.TokenType { //nolint:exhaustive // We expects only this types.
	case RuleTokenTypeArray:
		for i := len(rule.Children) - 1; i >= 0; i-- {
			r := rule.Children[i]
			if err := c.inheritPropertiesFromUserType(uut, r.ScalarValue, catalogUserTypes); err != nil {
				return err
			}
		}
	case RuleTokenTypeReference:
		if err := c.inheritPropertiesFromUserType(uut, rule.ScalarValue, catalogUserTypes); err != nil {
			return err
		}
	}
	return nil
}

func (c *ExchangeContent) inheritPropertiesFromUserType(
	uut *StringSet,
	userTypeName string,
	catalogUserTypes *UserTypes,
) error {
	ut, ok := catalogUserTypes.Get(userTypeName)
	if !ok {
		return fmt.Errorf(`%s (%s)`, jerr.UserTypeNotFound, userTypeName)
	}

	uts, ok := ut.Schema.(*ExchangeJSightSchema)
	if !ok {
		return errors.New(jerr.RuntimeFailure)
	}

	err := uts.Compile()
	if err != nil {
		return err
	}

	if uts.exchangeContent.TokenType != schema.TokenTypeObject {
		return fmt.Errorf("%s %q", jerr.UserTypeIsNotAnObject, userTypeName)
	}

	if c.Children == nil {
		c.Children = make([]*ExchangeContent, 0, 10)
	}

	for i := len(uts.exchangeContent.Children) - 1; i >= 0; i-- {
		cc := uts.exchangeContent.Children[i]

		if cc.Key == nil {
			return errors.New(jerr.RuntimeFailure)
		}

		p := c.ObjectProperty(*(cc.Key))
		if p != nil && p.InheritedFrom == "" {
			// Don't allow to override original properties.
			return fmt.Errorf(jerr.NotAllowedToOverrideTheProperty,
				*(cc.Key),
				userTypeName,
			)
		}

		if p != nil && p.InheritedFrom != "" {
			// This property already defined, skip.
			continue
		}

		dup := *cc
		dup.ToUsedUserTypes(uut)
		dup.InheritedFrom = userTypeName
		c.Unshift(&dup)
	}

	return nil
}

func (c *ExchangeContent) ToUsedUserTypes(uut *StringSet) {
	if c.TokenType == schema.TokenTypeShortcut {
		if c.Type == "mixed" {
			for _, ut := range strings.Split(c.ScalarValue, "|") {
				s := strings.TrimSpace(ut)
				if bytes.NewBytes(s).IsUserTypeName() {
					uut.Add(s)
				}
			}
		} else {
			uut.Add(c.Type)
		}
	}
}

func (c *ExchangeContent) IsObjectHaveProperty(k string) bool {
	return c.ObjectProperty(k) != nil
}

func (c *ExchangeContent) ObjectProperty(k string) *ExchangeContent {
	for _, v := range c.Children {
		if *(v.Key) == k {
			return v
		}
	}
	return nil
}

func (c *ExchangeContent) Unshift(v *ExchangeContent) {
	c.Children = append([]*ExchangeContent{v}, c.Children...)
}

func (c *ExchangeContent) MarshalJSON() (b []byte, err error) {
	switch c.TokenType {
	case schema.TokenTypeObject, schema.TokenTypeArray:
		b, err = c.marshalJSONObjectOrArray()

	default:
		b, err = c.marshalJSONLiteral()
	}
	return b, err
}

func (c *ExchangeContent) marshalJSONObjectOrArray() ([]byte, error) {
	var data struct {
		Rules            []Rule             `json:"rules,omitempty"`
		Key              *string            `json:"key,omitempty"`
		TokenType        string             `json:"tokenType,omitempty"`
		Type             string             `json:"type,omitempty"`
		InheritedFrom    string             `json:"inheritedFrom,omitempty"`
		Note             string             `json:"note,omitempty"`
		Children         []*ExchangeContent `json:"children"`
		IsKeyUserTypeRef bool               `json:"isKeyUserTypeRef,omitempty"`
		Optional         bool               `json:"optional"`
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
		data.Children = make([]*ExchangeContent, 0)
	} else {
		data.Children = c.Children
	}

	return json.Marshal(data)
}

func (c *ExchangeContent) marshalJSONLiteral() ([]byte, error) {
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

func astNodeToJsightContent(
	node schema.ASTNode,
	usedUserTypes, usedUserEnums *StringSet,
) *ExchangeContent {
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

	c := &ExchangeContent{
		IsKeyUserTypeRef: node.IsKeyShortcut,
		TokenType:        node.TokenType,
		Type:             node.SchemaType,
		Optional:         isOptional,
		ScalarValue:      node.Value,
		InheritedFrom:    node.InheritedFrom,
		Note:             Annotation(node.Comment),
		Rules:            rules,
	}

	switch node.TokenType {
	case schema.TokenTypeObject:
		c.collectJSightContentObjectProperties(node, usedUserTypes, usedUserEnums)
	case schema.TokenTypeArray:
		c.collectJSightContentArrayItems(node, usedUserTypes, usedUserEnums)
	}

	return c
}

func collectJSightContentRules(node schema.ASTNode, usedUserTypes *StringSet) *Rules {
	if node.Rules.Len() == 0 {
		return &Rules{}
	}

	rr := newRulesBuilder(node.Rules.Len())

	node.Rules.EachSafe(func(k string, v schema.RuleASTNode) {
		switch k {
		case "type":
			if v.Value[0] == '@' {
				usedUserTypes.Add(v.Value)
			}
			if v.Source == schema.RuleASTNodeSourceGenerated {
				return
			}

		case "allOf":
			if v.Value != "" {
				usedUserTypes.Add(v.Value)
			}

			for _, i := range v.Items {
				usedUserTypes.Add(i.Value)
			}

		case "additionalProperties":
			if v.Value[0] == '@' {
				usedUserTypes.Add(v.Value)
			}

		case "or":
			for _, i := range v.Items {
				var userType string
				if i.Value != "" {
					userType = i.Value
				} else {
					vv, ok := i.Properties.Get("type")
					if ok {
						userType = vv.Value
					} else {
						userType = node.SchemaType
					}
				}

				if userType[0] != '@' {
					continue
				}

				usedUserTypes.Add(userType)
			}

			if v.Source == schema.RuleASTNodeSourceGenerated {
				return
			}
		}
		rr.Set(k, astNodeToSchemaRule(v))
	})

	return rr.Rules()
}

func (c *ExchangeContent) collectJSightContentObjectProperties(
	node schema.ASTNode,
	usedUserTypes, usedUserEnums *StringSet,
) {
	if len(node.Children) > 0 {
		if c.Children == nil {
			c.Children = make([]*ExchangeContent, 0, len(node.Children))
		}
		for _, v := range node.Children {
			an := astNodeToJsightContent(v, usedUserTypes, usedUserEnums)
			an.Key = StrPtr(v.Key)

			c.Children = append(c.Children, an)

			if v.IsKeyShortcut {
				usedUserTypes.Add(v.Key)
			}
		}
	}
}

func (c *ExchangeContent) collectJSightContentArrayItems(
	node schema.ASTNode,
	usedUserTypes, usedUserEnums *StringSet,
) {
	if len(node.Children) > 0 {
		if c.Children == nil {
			c.Children = make([]*ExchangeContent, 0, len(node.Children))
		}
		for _, n := range node.Children {
			an := astNodeToJsightContent(n, usedUserTypes, usedUserEnums)
			an.Optional = true
			c.Children = append(c.Children, an)
		}
	}
}

func astNodeToSchemaRule(node schema.RuleASTNode) Rule {
	rr := newRulesBuilder(node.Properties.Len() + len(node.Items))

	if node.Properties.Len() != 0 {
		node.Properties.EachSafe(func(k string, v schema.RuleASTNode) {
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
		TokenType:   RuleTokenType(node.TokenType),
		ScalarValue: node.Value,
		Note:        node.Comment,
		Children:    children,
	}
}

func StrPtr(s string) *string {
	return &s
}

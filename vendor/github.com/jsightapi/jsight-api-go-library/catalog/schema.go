package catalog

import (
	"encoding/json"
	"fmt"
	"sync"

	jschemaLib "github.com/jsightapi/jsight-schema-go-library"

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

	// Example of schema.
	Example string
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
		TokenType:   RuleTokenType(node.JSONType),
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

	case notation.SchemaNotationRegex:
		data.Content = schema.ContentRegexp

	case notation.SchemaNotationAny, notation.SchemaNotationEmpty:
		// nothing

	default:
		return nil, fmt.Errorf(`invalid schema notation "%s"`, schema.Notation)
	}

	data.Example = schema.Example

	return json.Marshal(data)
}

func SrtPtr(s string) *string {
	return &s
}

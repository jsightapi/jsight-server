package catalog

import (
	"encoding/json"
)

type Rule struct {
	Key       string
	TokenType RuleTokenType

	// ScalarValue specified only if "tokenType" specifies "string", "number",
	// "boolean", "annotation" or "null".
	ScalarValue string

	// Note may contain line breaks.
	Note string

	// Children specified only if tokenType: "array" or "object".
	Children []Rule
}

type RuleTokenType string

const (
	RuleTokenTypeObject     RuleTokenType = "object"
	RuleTokenTypeArray      RuleTokenType = "array"
	RuleTokenTypeString     RuleTokenType = "string"
	RuleTokenTypeNumber     RuleTokenType = "number"
	RuleTokenTypeBoolean    RuleTokenType = "boolean"
	RuleTokenTypeNull       RuleTokenType = "null"
	RuleTokenTypeAnnotation RuleTokenType = "annotation"
	RuleTokenTypeReference  RuleTokenType = "reference"
)

var (
	_ json.Marshaler = Rule{}
	_ json.Marshaler = &Rule{}
)

func (r Rule) MarshalJSON() ([]byte, error) {
	var data struct {
		Key         string        `json:"key,omitempty"`
		TokenType   RuleTokenType `json:"tokenType"`
		ScalarValue string        `json:"scalarValue,omitempty"`
		Note        string        `json:"note,omitempty"`
		Children    []Rule        `json:"children,omitempty"`
	}

	data.Key = r.Key
	data.TokenType = r.TokenType
	data.ScalarValue = r.ScalarValue
	data.Note = r.Note
	data.Children = r.Children

	return json.Marshal(data)
}

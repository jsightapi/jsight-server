package catalog

import (
	"encoding/json"
)

type Rule struct {
	Key       string
	TokenType RuleTokenType

	// ScalarValue specified only if "tokenType" specifies "string", "number",
	// "boolean", "annotation", "null" or "reference".
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

func (r Rule) MarshalJSON() (b []byte, err error) {
	switch r.TokenType {
	case RuleTokenTypeObject, RuleTokenTypeArray:
		b, err = r.marshalJSONObjectOrArray()

	default:
		b, err = r.marshalJSONLiteral()
	}
	return b, err
}

func (r Rule) marshalJSONObjectOrArray() ([]byte, error) {
	var data struct {
		Key       string        `json:"key,omitempty"`
		TokenType RuleTokenType `json:"tokenType"`
		Note      string        `json:"note,omitempty"`
		Children  []Rule        `json:"children,omitempty"`
	}

	data.Key = r.Key
	data.TokenType = r.TokenType
	data.Note = r.Note
	data.Children = r.Children

	return json.Marshal(data)
}

func (r Rule) marshalJSONLiteral() ([]byte, error) {
	var data struct {
		Key         string        `json:"key,omitempty"`
		TokenType   RuleTokenType `json:"tokenType"`
		Note        string        `json:"note,omitempty"`
		ScalarValue string        `json:"scalarValue"`
	}

	data.Key = r.Key
	data.TokenType = r.TokenType
	data.Note = r.Note
	data.ScalarValue = r.ScalarValue

	return json.Marshal(data)
}

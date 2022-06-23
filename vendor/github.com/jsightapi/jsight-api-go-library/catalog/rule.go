package catalog

import (
	"encoding/json"
)

type Rule struct {
	Key         string
	TokenType   string
	ScalarValue string
	Note        string
	Children    []Rule
}

var (
	_ json.Marshaler = Rule{}
	_ json.Marshaler = &Rule{}
)

func (r Rule) MarshalJSON() ([]byte, error) {
	var data struct {
		Key         string `json:"key,omitempty"`
		TokenType   string `json:"tokenType"`
		ScalarValue string `json:"scalarValue,omitempty"`
		Note        string `json:"note,omitempty"`
		Children    []Rule `json:"children,omitempty"`
	}

	data.Key = r.Key
	data.TokenType = r.TokenType
	data.ScalarValue = r.ScalarValue
	data.Note = r.Note
	data.Children = r.Children

	return json.Marshal(data)
}

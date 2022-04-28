package catalog

import (
	"encoding/json"
)

type Rule struct {
	JsonType    string
	ScalarValue string
	Note        string
	Properties  *Rules
	Items       []Rule
}

var (
	_ json.Marshaler = Rule{}
	_ json.Marshaler = &Rule{}
)

func (r Rule) MarshalJSON() ([]byte, error) {
	var data struct {
		JsonType    string `json:"jsonType"`
		ScalarValue string `json:"scalarValue,omitempty"`
		Note        string `json:"note,omitempty"`
		Properties  *Rules `json:"properties,omitempty"`
		Items       []Rule `json:"items,omitempty"`
	}

	data.JsonType = r.JsonType
	data.ScalarValue = r.ScalarValue
	data.Note = r.Note
	if r.Properties != nil && r.Properties.Len() > 0 {
		data.Properties = r.Properties
	}
	data.Items = r.Items

	return json.Marshal(data)
}

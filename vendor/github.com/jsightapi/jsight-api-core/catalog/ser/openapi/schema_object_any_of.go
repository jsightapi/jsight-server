package openapi

import (
	"encoding/json"
)

func schemaObjectForAnyOf(anyOf []schemaObject) schemaObject {
	return &schemaObjectAnyOf{anyOf, ""}
}

type schemaObjectAnyOf struct {
	AnyOf       []schemaObject `json:"anyOf"`
	Description string         `json:"description,omitempty"`
}

func (s *schemaObjectAnyOf) SetDescription(d string) {
	s.Description = d
}

func (s schemaObjectAnyOf) MarshalJSON() (b []byte, err error) {
	type Alias schemaObjectAnyOf
	return json.Marshal(&struct {
		Alias
	}{
		Alias(s),
	})
}

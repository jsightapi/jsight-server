package openapi

import "encoding/json"

type schemaObjectAny struct {
	Description string `json:"description,omitempty"`
}

func (s *schemaObjectAny) SetDescription(d string) {
	s.Description = d
}

func (s schemaObjectAny) MarshalJSON() (b []byte, err error) {
	type Alias schemaObjectAny
	return json.Marshal(&struct {
		Alias
	}{
		Alias: Alias(s),
	})
}

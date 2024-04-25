package openapi

// Common fields for HeaderObject and ParameterObject
// In OAS Header Object is officially a subset of Parameter Object
type ParameterBase struct {
	Description string         `json:"description,omitempty"`
	Required    bool           `json:"required,omitempty"`
	Style       parameterStyle `json:"style,omitempty"`
	Explode     bool           `json:"explode,omitempty"`
	Schema      schemaObject   `json:"schema"`
}

package openapi

type ParameterObject struct {
	ParameterBase
	Name string            `json:"name"`
	In   parameterLocation `json:"in"`
}

func newParameterObject(
	in parameterLocation,
	name string,
	description string,
	required bool,
	schema schemaObject,
) *ParameterObject {
	return &ParameterObject{
		In:   in,
		Name: name,
		ParameterBase: ParameterBase{
			Description: description,
			Required:    required,
			Schema:      schema,
		},
	}
}

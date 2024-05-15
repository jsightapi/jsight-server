package openapi

type HeaderObject struct {
	ParameterBase
}

func newHeaderObject(required bool, description string, schema schemaObject) *HeaderObject {
	return &HeaderObject{
		ParameterBase{
			Required:    required,
			Description: description,
			Schema:      schema,
		}}
}

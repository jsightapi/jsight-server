package openapi

type parameterStyle string

const (
	ParameterStyleMatrix     parameterStyle = "matrix"
	ParameterStyleLabel      parameterStyle = "label"
	ParameterStyleForm       parameterStyle = "form"
	ParameterStyleSimple     parameterStyle = "simple"
	ParameterStyleDeepObject parameterStyle = "deepObject"
)

package openapi

type parameterLocation string

const (
	ParameterLocationPath   parameterLocation = "path"
	ParameterLocationQuery  parameterLocation = "query"
	ParameterLocationHeader parameterLocation = "header"
	ParameterLocationCookie parameterLocation = "cookie"
)

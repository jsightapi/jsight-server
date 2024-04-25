package openapi

type SchemaInfoType int

const (
	SchemaInfoTypeRegex SchemaInfoType = iota
	SchemaInfoTypeObject
	SchemaInfoTypeArray
	SchemaInfoTypeScalar // string, integer, boolean, null
	SchemaInfoTypeAny
	SchemaInfoTypeReference
)

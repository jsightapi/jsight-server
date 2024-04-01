package openapi

type SchemaInformer interface {
	Type() SchemaInfoType

	// SchemaObject returns an OpenAPI Schema Object based on the JSight schema
	SchemaObject() SchemaObject

	// Annotation returns the JSight schema annotation
	Annotation() string
}

type ObjectInformer interface {
	SchemaInformer

	// PropertiesInfos returns properties of the object
	PropertiesInfos() []PropertyInformer
}

type PropertyInformer interface {
	SchemaInformer

	// Key returns the object key name
	Key() string

	// Optional returns the value for the "optional" rule of the object property
	Optional() bool
}

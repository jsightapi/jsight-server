package openapi

type parameterInfo interface {
	name() string
	optional() bool
	schemaObject() schemaObject
	annotation() string
}

type paramInfo struct {
	name_         string
	optional_     bool
	schemaObject_ schemaObject
	annotation_   string
}

func (p paramInfo) name() string {
	return p.name_
}

func (p paramInfo) optional() bool {
	return p.optional_
}

func (p paramInfo) schemaObject() schemaObject {
	return p.schemaObject_
}

func (p paramInfo) annotation() string {
	return p.annotation_
}

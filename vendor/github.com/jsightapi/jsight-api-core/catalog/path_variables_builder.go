package catalog

import (
	"github.com/jsightapi/jsight-schema-core/notations/jschema/ischema"
)

type PathVariablesBuilder struct {
	catalogUserTypes *UserTypes
	objectBuilder    ObjectBuilder
}

func NewPathVariablesBuilder(catalogUserTypes *UserTypes) PathVariablesBuilder {
	return PathVariablesBuilder{
		catalogUserTypes: catalogUserTypes,
		objectBuilder:    NewObjectBuilder(),
	}
}

func (b PathVariablesBuilder) AddProperty(key string, node ischema.Node, types map[string]ischema.Type) {
	b.objectBuilder.AddProperty(key, node, types)
}

func (b PathVariablesBuilder) Len() int {
	return b.objectBuilder.Len()
}

func (b PathVariablesBuilder) Build() *PathVariables {
	uutNames := b.objectBuilder.UserTypeNames()
	for _, name := range uutNames {
		if ut, ok := b.catalogUserTypes.Get(name); ok {
			switch es := ut.Schema.(type) {
			case *ExchangeJSightSchema:
				b.objectBuilder.AddType(name, es.JSchema)
			case *ExchangeRegexSchema:
				b.objectBuilder.AddType(name, es.RSchema)
			}
		}
	}

	s := b.objectBuilder.Build()

	es := newExchangeJSightSchema(s)
	es.disableExchangeExample = true
	es.catalogUserTypes = b.catalogUserTypes

	return &PathVariables{
		Schema: es,
	}
}

package openapi

import (
	"github.com/jsightapi/jsight-api-core/catalog"
)

type ComponentsSchemas map[string]schemaObject

func newSchemas(tt *catalog.UserTypes) ComponentsSchemas {
	if tt.Len() == 0 {
		return nil
	}

	ss := make(ComponentsSchemas, tt.Len())
	_ = tt.Each(func(name string, ut *catalog.UserType) error {
		typeSchemaObject := schemaObjectFromExchangeSchema(ut.Schema)
		typeSchemaObject.SetDescription(ut.Annotation)

		ss[typeNameToSchemaName(name)] = typeSchemaObject

		return nil
	})

	return ss
}

// all names in JSight start with `@`
func typeNameToSchemaName(n string) string {
	return n[1:]
}

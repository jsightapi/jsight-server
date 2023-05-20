package loader

import "github.com/jsightapi/jsight-schema-core/notations/jschema/ischema"

func AddUnnamedTypes(rootSchema *ischema.ISchema) {
	for _, typ := range rootSchema.TypesList() {
		for unnamed, unnamedTyp := range typ.Schema.TypesList() {
			rootSchema.AddType(unnamed, unnamedTyp)
		}
	}
}

package loader

import "github.com/jsightapi/jsight-schema-go-library/notations/jschema/internal/schema"

func AddUnnamedTypes(rootSchema *schema.Schema) {
	for _, typ := range rootSchema.TypesList() {
		for unnamed, unnamedTyp := range typ.Schema().TypesList() {
			rootSchema.AddType(unnamed, unnamedTyp)
		}
	}
}

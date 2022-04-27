package loader

import (
	"github.com/jsightapi/jsight-schema-go-library/fs"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema/internal/scanner"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema/internal/schema"
)

// NewSchemaForSdk reads the Schema from a file without adding to the collection.
// Does not compile allOf, in order that before there was a possibility to add
// additional TYPES.
func NewSchemaForSdk(file *fs.File, areKeysOptionalByDefault bool) *schema.Schema {
	return LoadSchema(scanner.NewSchemaScanner(file, false), nil, areKeysOptionalByDefault)
}

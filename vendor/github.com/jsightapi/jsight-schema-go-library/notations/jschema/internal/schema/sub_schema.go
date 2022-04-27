package schema

import (
	"github.com/jsightapi/jsight-schema-go-library/bytes"
	"github.com/jsightapi/jsight-schema-go-library/fs"
)

type Type struct {
	schema   *Schema
	rootFile *fs.File
	begin    bytes.Index
}

func (s *Type) Schema() *Schema {
	return s.schema
}

func (s *Type) RootFile() *fs.File {
	return s.rootFile
}

func (s *Type) Begin() bytes.Index {
	return s.begin
}

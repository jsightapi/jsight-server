package scanner

import (
	"github.com/jsightapi/jsight-schema-core/bytes"
	"github.com/jsightapi/jsight-schema-core/fs"
	"github.com/jsightapi/jsight-schema-core/kit"
	"github.com/jsightapi/jsight-schema-core/notations/jschema"

	"github.com/jsightapi/jsight-api-core/jerr"
)

// Pass rest of the file to jsc scanner to find out where jschema ends
func stateJSchema(s *Scanner, _ byte) *jerr.JApiError {
	s.found(SchemaBegin)
	schemaLength, je := s.readSchemaWithJsc()
	if je != nil {
		return je
	}
	if schemaLength > 0 {
		s.curIndex += bytes.Index(schemaLength - 1)
	}
	s.step = stateSchemaClosed
	return nil
}

func (s *Scanner) readSchemaWithJsc() (uint, *jerr.JApiError) {
	fc := s.file.Content()
	file := fs.NewFile("", fc.Sub(s.curIndex, fc.LenIndex()))

	l, err := jschema.FromFile(file).Len()
	if err != nil {
		err := kit.ConvertError(file, err)
		return 0, s.japiError(err.Message(), s.curIndex+bytes.Index(err.Index()))
	}
	return l, nil
}

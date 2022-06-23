package scanner

import (
	"github.com/jsightapi/jsight-schema-go-library/bytes"
	"github.com/jsightapi/jsight-schema-go-library/fs"
	"github.com/jsightapi/jsight-schema-go-library/kit"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema"

	"github.com/jsightapi/jsight-api-go-library/jerr"
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
	file := fs.NewFile("", fc.Slice(s.curIndex, bytes.Index(fc.Len()-1)))

	l, err := jschema.FromFile(file).Len()
	if err != nil {
		err := kit.ConvertError(file, err)
		return 0, s.japiError(err.Message(), s.curIndex+bytes.Index(err.Position()))
	}
	return l, nil
}

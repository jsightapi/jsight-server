package catalog

import (
	"encoding/json"

	"github.com/jsightapi/jsight-schema-core/bytes"
	"github.com/jsightapi/jsight-schema-core/notations/regex"

	"github.com/jsightapi/jsight-api-core/notation"
)

type ExchangeRegexSchema struct {
	*regex.RSchema
}

func (e ExchangeRegexSchema) MarshalJSON() ([]byte, error) {
	data := struct {
		Content  interface{}             `json:"content,omitempty"`
		Example  string                  `json:"example,omitempty"`
		Notation notation.SchemaNotation `json:"notation"`
	}{
		Notation: notation.SchemaNotationRegex,
	}

	var err error

	data.Content, err = e.Pattern()
	if err != nil {
		return []byte{}, err
	}

	example, err := e.Example()
	if err != nil {
		return []byte{}, err
	}

	data.Example = string(example)

	return json.Marshal(data)
}

func (e ExchangeRegexSchema) Notation() notation.SchemaNotation {
	return notation.SchemaNotationRegex
}

func NewExchangeRegexSchema(regexStr bytes.Bytes) (*ExchangeRegexSchema, error) {
	s := regex.New("", regexStr)
	return &ExchangeRegexSchema{RSchema: s}, nil
}

func newExchangeRegexSchema(s *regex.RSchema) *ExchangeRegexSchema {
	return &ExchangeRegexSchema{RSchema: s}
}

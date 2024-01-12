package catalog

import (
	"encoding/json"
	"errors"

	schema "github.com/jsightapi/jsight-schema-core"

	"github.com/jsightapi/jsight-api-core/jerr"

	"github.com/jsightapi/jsight-api-core/notation"
)

type ExchangePseudoSchema struct {
	schema.Schema
	Notation notation.SchemaNotation
}

func NewExchangePseudoSchema(n notation.SchemaNotation) *ExchangePseudoSchema {
	return &ExchangePseudoSchema{
		Notation: n,
	}
}

func (e ExchangePseudoSchema) MarshalJSON() ([]byte, error) {
	if e.Notation != notation.SchemaNotationAny && e.Notation != notation.SchemaNotationEmpty {
		return nil, errors.New(jerr.RuntimeFailure)
	}

	data := struct {
		Notation notation.SchemaNotation `json:"notation"`
	}{
		Notation: e.Notation,
	}

	return json.Marshal(data)
}

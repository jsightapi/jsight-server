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
	notation notation.SchemaNotation
}

func (e ExchangePseudoSchema) Notation() notation.SchemaNotation {
	return e.notation
}

func NewExchangePseudoSchema(n notation.SchemaNotation) *ExchangePseudoSchema {
	return &ExchangePseudoSchema{
		notation: n,
	}
}

func (e ExchangePseudoSchema) MarshalJSON() ([]byte, error) {
	if e.notation != notation.SchemaNotationAny && e.notation != notation.SchemaNotationEmpty {
		return nil, errors.New(jerr.RuntimeFailure)
	}

	data := struct {
		Notation notation.SchemaNotation `json:"notation"`
	}{
		Notation: e.notation,
	}

	return json.Marshal(data)
}

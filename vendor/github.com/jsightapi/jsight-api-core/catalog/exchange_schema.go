package catalog

import (
	"github.com/jsightapi/jsight-api-core/notation"

	schema "github.com/jsightapi/jsight-schema-core"
)

type ExchangeSchema interface {
	schema.Schema
	MarshalJSON() ([]byte, error)
	Notation() notation.SchemaNotation
}

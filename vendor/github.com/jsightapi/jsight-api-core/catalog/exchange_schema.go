package catalog

import (
	schema "github.com/jsightapi/jsight-schema-core"
)

type ExchangeSchema interface {
	schema.Schema
	MarshalJSON() ([]byte, error)
}

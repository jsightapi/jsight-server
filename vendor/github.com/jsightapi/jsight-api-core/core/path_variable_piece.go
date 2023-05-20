package core

import (
	"github.com/jsightapi/jsight-schema-core/notations/jschema/ischema"
)

type PieceOfPathVariable struct {
	node  ischema.Node
	types map[string]ischema.Type
}

package core

import (
	"github.com/jsightapi/jsight-schema-core/notations/jschema/ischema"
)

type PieceOfPathVariable struct {
	node  ischema.Node
	types map[string]ischema.Type

	// temp workaround. true means that this was not gathered from Path directive,
	// but from URL or Method-directive, imitating real rawPathVariable
	imitated bool
}

// imitates piece for param, gatherd from URL or Method-directive, not from Path schema.
func PieceOfPathVariableImitation() PieceOfPathVariable {
	return PieceOfPathVariable{
		node:     ischema.VirtualNodeForAny(),
		imitated: true,
	}
}

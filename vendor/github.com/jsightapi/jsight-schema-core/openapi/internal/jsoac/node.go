package jsoac

import (
	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/errs"
)

type Node interface {
	SetNodeDescription(string)
}

func newNode(astNode schema.ASTNode) Node {
	if astNode.SchemaType == stringAny {
		return newAny(astNode)
	}

	switch astNode.TokenType {
	case schema.TokenTypeNumber, schema.TokenTypeString, schema.TokenTypeBoolean:
		return newPrimitive(astNode)
	case schema.TokenTypeArray:
		return newArray(astNode)
	case schema.TokenTypeObject:
		return newObject(astNode)
	case schema.TokenTypeNull:
		return newNull(astNode)
	case schema.TokenTypeShortcut:
		return newRef(astNode)
	default:
		panic(errs.ErrRuntimeFailure.F())
	}
}

package jsoac

import (
	"encoding/json"

	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/errs"
)

type OADType int // OpenAPI Data Types

//go:generate stringer -type=OADType -linecomment
const (
	OADTypeString  OADType = iota // string
	OADTypeInteger                // integer
	OADTypeNumber                 // number
	OADTypeBoolean                // boolean
	OADTypeArray                  // array
	OADTypeObject                 // object
)

func (t OADType) MarshalJSON() (b []byte, err error) {
	return json.Marshal(t.String())
}

func oadTypeFromASTNode(astNode schema.ASTNode) OADType {
	if astNode.SchemaType == "integer" {
		return OADTypeInteger
	}
	switch astNode.TokenType {
	case schema.TokenTypeNumber:
		return OADTypeNumber
	case schema.TokenTypeString:
		return OADTypeString
	case schema.TokenTypeBoolean:
		return OADTypeBoolean
	case schema.TokenTypeArray:
		return OADTypeArray
	case schema.TokenTypeObject:
		return OADTypeObject
	default:
		// schema.TokenTypeShortcut:
		panic(errs.ErrRuntimeFailure.F())
	}
}

func oadTypeFromSchemaType(s string) OADType {
	switch s {
	case string(schema.SchemaTypeString), string(schema.SchemaTypeEmail), string(schema.SchemaTypeURI),
		string(schema.SchemaTypeUUID), string(schema.SchemaTypeDate), string(schema.SchemaTypeDateTime):
		return OADTypeString
	case string(schema.SchemaTypeInteger):
		return OADTypeInteger
	case string(schema.SchemaTypeFloat):
		return OADTypeNumber
	case string(schema.SchemaTypeBoolean):
		return OADTypeBoolean
	case string(schema.SchemaTypeObject):
		return OADTypeObject
	default:
		panic(errs.ErrRuntimeFailure.F())
	}
}

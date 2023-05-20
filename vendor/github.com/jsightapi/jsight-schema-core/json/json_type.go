//nolint:goconst // Not important here.
package json

import (
	"github.com/jsightapi/jsight-schema-core/bytes"
	"github.com/jsightapi/jsight-schema-core/errs"
)

type Type uint8

const (
	// TypeUndefined default value for literal and mixed nodes.
	TypeUndefined Type = iota
	TypeObject
	TypeArray
	TypeString
	TypeInteger
	// TypeFloat to be precise, there is no separate "Integer" and "Float" in JSON,
	// there is a single "Number" type. But in our case, we will assume that there is.
	TypeFloat
	TypeBoolean
	TypeNull

	// TypeMixed indicates that here can be anything.
	TypeMixed
)

func NewJsonType(b bytes.Bytes) Type {
	switch b.String() {
	case "object":
		return TypeObject
	case "array":
		return TypeArray
	case "string":
		return TypeString
	case "integer":
		return TypeInteger
	case "float":
		return TypeFloat
	case "boolean":
		return TypeBoolean
	case nullStr:
		return TypeNull
	}
	panic(errs.ErrUnknownValueOfTheTypeRule.F(b.String()))
}

var AllTypes = []Type{
	TypeObject,
	TypeArray,
	TypeString,
	TypeInteger,
	TypeFloat,
	TypeBoolean,
	TypeNull,
	TypeMixed,
}

func (t Type) String() string {
	switch t {
	case TypeObject:
		return "object"
	case TypeArray:
		return "array"
	case TypeString:
		return "string"
	case TypeInteger:
		return "integer"
	case TypeFloat:
		return "float"
	case TypeBoolean:
		return "boolean"
	case TypeNull:
		return nullStr
	case TypeMixed:
		return "mixed"
	default:
		return "unknown"
	}
}

func (t Type) IsLiteralType() bool {
	switch t { //nolint:exhaustive // It's okay.
	case TypeString, TypeBoolean, TypeInteger, TypeFloat, TypeNull, TypeMixed:
		return true
	}
	return false
}

func (t Type) ToTokenType() string {
	switch t { //nolint:exhaustive // We return an empty string.
	case TypeObject:
		return "object"
	case TypeArray:
		return "array"
	case TypeString:
		return "string"
	case TypeInteger, TypeFloat:
		return "number"
	case TypeBoolean:
		return "boolean"
	case TypeNull:
		return nullStr
	case TypeMixed:
		return "reference"
	}
	return ""
}

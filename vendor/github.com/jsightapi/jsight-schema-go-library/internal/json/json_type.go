//nolint:goconst // Not important here.
package json

import (
	"github.com/jsightapi/jsight-schema-go-library/bytes"
	"github.com/jsightapi/jsight-schema-go-library/errors"
)

type Type uint8

const (
	TypeUndefined Type = iota // default value for literal and mixed nodes
	TypeObject
	TypeArray
	TypeString
	TypeInteger
	TypeFloat // To be precise, there is no separate "Integer" and "Float" in JSON, there is a single "Number" type. But in our case, we will assume that there is.
	TypeBoolean
	TypeNull

	// TypeMixed indicates that here can be anything.
	TypeMixed
)

func NewJsonType(bytes bytes.Bytes) Type {
	switch string(bytes) {
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
	case "null":
		return TypeNull
	}
	panic(errors.Format(errors.ErrUnknownType, string(bytes)))
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
		return "null"
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
		return "null"
	case TypeMixed:
		return "reference"
	}
	return ""
}

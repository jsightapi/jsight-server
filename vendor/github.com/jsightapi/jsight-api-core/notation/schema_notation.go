package notation

import (
	"errors"

	"github.com/jsightapi/jsight-api-core/jerr"
)

// SchemaNotation represent available schema notations.
type SchemaNotation string

const (
	SchemaNotationJSight SchemaNotation = "jsight"
	SchemaNotationRegex  SchemaNotation = "regex"
	SchemaNotationAny    SchemaNotation = "any"
	SchemaNotationEmpty  SchemaNotation = "empty"
)

func NewSchemaNotation(sn string) (SchemaNotation, error) {
	switch sn {
	case "jsight", "":
		return SchemaNotationJSight, nil
	case "regex":
		return SchemaNotationRegex, nil
	case "any":
		return SchemaNotationAny, nil
	case "empty":
		return SchemaNotationEmpty, nil
	default:
		return "", errors.New(jerr.UnknownNotation)
	}
}

func (sn SchemaNotation) IsAnyOrEmpty() bool {
	return sn == SchemaNotationAny || sn == SchemaNotationEmpty
}

package catalog

import (
	"fmt"

	"github.com/jsightapi/jsight-api-go-library/notation"
)

type SerializeFormat string

const (
	SerializeFormatJSON        SerializeFormat = "json"
	SerializeFormatPlainString SerializeFormat = "plainString"
	SerializeFormatBinary      SerializeFormat = "binary"
	// SerializeFormatHtmlFormEncoded SerializeFormat = "htmlFormEncoded"
	// SerializeFormatNoFormat        SerializeFormat = "noFormat"
)

func SchemaSerializeFormat(n notation.SchemaNotation) (SerializeFormat, error) {
	switch n {
	case notation.SchemaNotationJSight:
		return SerializeFormatJSON, nil
	case notation.SchemaNotationRegex:
		return SerializeFormatPlainString, nil
	case notation.SchemaNotationAny:
		return SerializeFormatBinary, nil
	case notation.SchemaNotationEmpty:
		return SerializeFormatBinary, nil
	default:
		return "", fmt.Errorf(`unknown serialize format for schema notation "%s"`, n)
	}
}

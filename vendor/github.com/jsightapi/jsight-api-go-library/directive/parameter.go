package directive

import (
	stdBytes "bytes"
	"fmt"

	"github.com/jsightapi/jsight-schema-go-library/bytes"

	"github.com/jsightapi/jsight-api-go-library/jerr"
	"github.com/jsightapi/jsight-api-go-library/notation"
)

func unescapeParameter(b bytes.Bytes) bytes.Bytes {
	c := b.Unquote()
	if len(c) != 0 && len(c) != len(b) {
		c = stdBytes.ReplaceAll(c, []byte(`\"`), []byte(`"`))
		c = stdBytes.ReplaceAll(c, []byte(`\\`), []byte(`\`))
	}
	return c
}

func IsArrayOfTypes(b bytes.Bytes) bool {
	l := len(b)
	if l >= 4 && b[0] == '[' && b[l-1] == ']' {
		c := b[1 : l-1]
		if c.IsUserTypeName() {
			return true
		}
	}
	return false
}

func (d *Directive) AppendParameter(b bytes.Bytes) error {
	b = unescapeParameter(b)
	s := b.String()

	switch d.Type() { //nolint:exhaustive // We catch all uncovered enumeration.
	case Url, Get, Post, Put, Patch, Delete:
		return d.SetParameter("Path", s)

	case Request, HTTPResponseCode, Body:
		switch {
		case isSchemaNotation(s):
			return d.SetParameter("SchemaNotation", s)
		case IsArrayOfTypes(b):
			return d.SetParameter("Type", s)
		case b.IsUserTypeName():
			return d.SetParameter("Type", s)
		}

	case Type:
		switch {
		case isSchemaNotation(s):
			return d.SetParameter("SchemaNotation", s)
		case IsArrayOfTypes(b):
			return d.SetParameter("Name", s)
		case b.IsUserTypeName():
			return d.SetParameter("Name", s)
		}

	case Query:
		switch s {
		case "htmlFormEncoded", "noFormat":
			return d.SetParameter("Format", s)
		default:
			return d.SetParameter("QueryExample", s)
		}

	case Jsight, Version:
		return d.SetParameter("Version", s)

	case Title:
		return d.SetParameter("Title", s)

	case BaseUrl:
		return d.SetParameter("Path", s)

	case Server, Enum, Macro, Paste:
		if b.IsUserTypeName() {
			return d.SetParameter("Name", s)
		}

	case Protocol:
		return d.SetParameter("Protocol", s)

	case Method:
		return d.SetParameter("Method", s)
	}

	return fmt.Errorf("%s %q", jerr.IncorrectParameter, s)
}

func isSchemaNotation(s string) bool {
	_, e := notation.NewSchemaNotation(s)
	return e == nil
}

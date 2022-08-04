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
		return d.SetNamedParameter("Path", s)

	case Request, HTTPResponseCode, Body:
		switch {
		case isSchemaNotation(s):
			return d.SetNamedParameter("SchemaNotation", s)
		case IsArrayOfTypes(b):
			return d.SetNamedParameter("Type", s)
		case b.IsUserTypeName():
			return d.SetNamedParameter("Type", s)
		}

	case Type:
		switch {
		case isSchemaNotation(s):
			return d.SetNamedParameter("SchemaNotation", s)
		case IsArrayOfTypes(b):
			return d.SetNamedParameter("Name", s)
		case b.IsUserTypeName():
			return d.SetNamedParameter("Name", s)
		}

	case Query:
		switch s {
		case "htmlFormEncoded", "noFormat":
			return d.SetNamedParameter("Format", s)
		default:
			return d.SetNamedParameter("QueryExample", s)
		}

	case Jsight, Version:
		return d.SetNamedParameter("Version", s)

	case Title:
		return d.SetNamedParameter("Title", s)

	case BaseUrl:
		return d.SetNamedParameter("Path", s)

	case Server, Enum, Macro, Paste:
		if b.IsUserTypeName() {
			return d.SetNamedParameter("Name", s)
		}

	case Protocol:
		return d.SetNamedParameter("ProtocolName", s)

	case Method:
		return d.SetNamedParameter("MethodName", s)

	case TAG:
		if b.IsUserTypeName() {
			return d.SetNamedParameter("TagName", s)
		}

	case Tags:
		if b.IsUserTypeName() {
			d.AppendUnnamedParameter(s)
			return nil
		}
	}

	return fmt.Errorf("%s %q", jerr.IncorrectParameter, s)
}

func isSchemaNotation(s string) bool {
	_, e := notation.NewSchemaNotation(s)
	return e == nil
}

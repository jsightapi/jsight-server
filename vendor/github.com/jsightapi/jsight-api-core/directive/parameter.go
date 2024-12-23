package directive

import (
	"fmt"

	"github.com/jsightapi/jsight-schema-core/bytes"

	"github.com/jsightapi/jsight-api-core/jerr"
	"github.com/jsightapi/jsight-api-core/notation"
)

func IsArrayOfTypes(b bytes.Bytes) bool {
	l := b.Len()
	if l >= 4 && b.FirstByte() == '[' && b.LastByte() == ']' {
		c := b.Sub(1, l-1)
		if c.IsUserTypeName() {
			return true
		}
	}
	return false
}

func (d *Directive) AppendParameter(b bytes.Bytes) error {
	b = b.Unquote()
	s := b.String()

	switch d.Type() { //nolint:exhaustive // We catch all uncovered enumeration.
	case URL, Get, Post, Put, Patch, Delete:
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

	case BaseURL:
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

	case OperationID:
		return d.SetNamedParameter("OperationId", s)
	}

	return fmt.Errorf("%s %q", jerr.IncorrectParameter, s)
}

func isSchemaNotation(s string) bool {
	_, e := notation.NewSchemaNotation(s)
	return e == nil
}

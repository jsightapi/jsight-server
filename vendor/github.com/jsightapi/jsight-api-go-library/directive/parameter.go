package directive

import (
	by "bytes"
	"fmt"

	"github.com/jsightapi/jsight-schema-go-library/bytes"

	"github.com/jsightapi/jsight-api-go-library/jerr"
	"github.com/jsightapi/jsight-api-go-library/notation"
)

func unescapeParameter(b bytes.Bytes) bytes.Bytes {
	c := b.Unquote()
	if len(c) != 0 && len(c) != len(b) {
		c = by.ReplaceAll(c, []byte(`\"`), []byte(`"`))
		c = by.ReplaceAll(c, []byte(`\\`), []byte(`\`))
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
	var err error

	b = unescapeParameter(b)
	s := b.String()

	switch d.Type() {
	case Url, Get, Post, Put, Patch, Delete:
		err = d.SetParameter("Path", s)

	case Request, HTTPResponseCode, Body:
		if _, e := notation.NewSchemaNotation(s); e == nil { //nolint:gocritic // ifElseChain: rewrite if-else to switch statement
			err = d.SetParameter("SchemaNotation", s)
		} else if IsArrayOfTypes(b) {
			err = d.SetParameter("Type", b.String())
		} else if b.IsUserTypeName() {
			err = d.SetParameter("Type", s)
		} else {
			err = fmt.Errorf("%s %q", jerr.IncorrectParameter, s)
		}

	case Type:
		if _, e := notation.NewSchemaNotation(s); e == nil { //nolint:gocritic // ifElseChain: rewrite if-else to switch statement
			err = d.SetParameter("SchemaNotation", s)
		} else if IsArrayOfTypes(b) {
			err = d.SetParameter("Name", b.String())
		} else if b.IsUserTypeName() {
			err = d.SetParameter("Name", s)
		} else {
			err = fmt.Errorf("%s %q", jerr.IncorrectParameter, s)
		}

	case Query:
		switch s {
		case "htmlFormEncoded", "noFormat":
			err = d.SetParameter("Format", s)
		default:
			err = d.SetParameter("QueryExample", s)
		}

	case Jsight:
		err = d.SetParameter("Version", s)

	case Title:
		err = d.SetParameter("Title", s)

	case Version:
		err = d.SetParameter("Version", s)

	case BaseUrl:
		err = d.SetParameter("Path", s)

	case Server, Enum, Macro, Paste:
		if b.IsUserTypeName() {
			err = d.SetParameter("Name", s)
		} else {
			err = fmt.Errorf("%s %q", jerr.IncorrectParameter, s)
		}

	default:
		err = fmt.Errorf("%s %q", jerr.IncorrectParameter, s)
	}

	return err
}

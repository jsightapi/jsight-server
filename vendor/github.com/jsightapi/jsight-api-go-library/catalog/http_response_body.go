package catalog

import (
	"errors"

	jschemaLib "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/bytes"
	"github.com/jsightapi/jsight-schema-go-library/kit"

	"github.com/jsightapi/jsight-api-go-library/directive"
	"github.com/jsightapi/jsight-api-go-library/jerr"
	"github.com/jsightapi/jsight-api-go-library/notation"
)

type HTTPResponseBody struct {
	Format    SerializeFormat     `json:"format"`
	Schema    *Schema             `json:"schema"`
	Directive directive.Directive `json:"-"`
}

func NewHTTPResponseBody(
	b bytes.Bytes,
	f SerializeFormat,
	sn notation.SchemaNotation,
	d directive.Directive,
	tt *UserSchemas,
	rr map[string]jschemaLib.Rule,
) (HTTPResponseBody, *jerr.JApiError) {
	body := HTTPResponseBody{
		Format:    f,
		Schema:    nil,
		Directive: d,
	}

	var s Schema
	switch f {
	case SerializeFormatJSON:
		var err error
		s, err = UnmarshalJSightSchema("", b, tt, rr)
		if err != nil {
			return HTTPResponseBody{}, adoptErrorForResponseBody(d, err)
		}
	case SerializeFormatPlainString:
		var err error
		s, err = UnmarshalRegexSchema("", b)
		if err != nil {
			return HTTPResponseBody{}, adoptErrorForResponseBody(d, err)
		}
	default:
		s = NewSchema(sn)
	}

	body.Schema = &s

	return body, nil
}

func adoptErrorForResponseBody(d directive.Directive, err error) *jerr.JApiError {
	var e kit.Error
	if errors.As(err, &e) {
		if d.BodyCoords.IsSet() {
			return d.BodyErrorIndex(e.Message(), e.Position())
		}
		return d.ParameterError(e.Message())
	}
	return d.KeywordError(err.Error())
}

package catalog

import (
	"errors"

	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/bytes"
	"github.com/jsightapi/jsight-schema-core/kit"

	"github.com/jsightapi/jsight-api-core/directive"
	"github.com/jsightapi/jsight-api-core/jerr"
	"github.com/jsightapi/jsight-api-core/notation"
)

type HTTPResponseBody struct {
	Format    SerializeFormat     `json:"format"`
	Schema    ExchangeSchema      `json:"schema"`
	Directive directive.Directive `json:"-"`
}

func NewHTTPResponseBody(
	b bytes.Bytes,
	f SerializeFormat,
	sn notation.SchemaNotation,
	d directive.Directive,
	tt *UserSchemas,
	rr map[string]schema.Rule,
	catalogUserTypes *UserTypes,
) (HTTPResponseBody, *jerr.JApiError) {
	body := HTTPResponseBody{
		Format:    f,
		Schema:    nil,
		Directive: d,
	}

	var s ExchangeSchema
	var err error

	switch f {
	case SerializeFormatJSON:
		s, err = NewExchangeJSightSchema(b.Data(), tt, rr, catalogUserTypes)
		if err != nil {
			return HTTPResponseBody{}, adoptErrorForResponseBody(d, err)
		}
	case SerializeFormatPlainString:
		s, err = NewExchangeRegexSchema(b)
		if err != nil {
			return HTTPResponseBody{}, adoptErrorForResponseBody(d, err)
		}
	default:
		s = NewExchangePseudoSchema(sn)
	}

	body.Schema = s

	return body, nil
}

func adoptErrorForResponseBody(d directive.Directive, err error) *jerr.JApiError {
	var e kit.Error
	if errors.As(err, &e) {
		if d.BodyCoords.IsSet() {
			return d.BodyErrorIndex(e.Message(), e.Index())
		}
		return d.ParameterError(e.Message())
	}
	return d.KeywordError(err.Error())
}

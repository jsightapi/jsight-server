package directive

import (
	"github.com/jsightapi/jsight-api-go-library/jerr"

	"github.com/jsightapi/jsight-schema-go-library/bytes"
)

// TODO error interface as parameter functions. Wrap error.

func (d Directive) KeywordError(msg string) *jerr.JAPIError {
	return jerr.NewJAPIError(msg, d.keywordCoords.File(), d.keywordCoords.b)
}

func (d Directive) BodyError(msg string) *jerr.JAPIError {
	if d.BodyCoords.IsSet() {
		return jerr.NewJAPIError(msg, d.BodyCoords.File(), d.BodyCoords.b)
	}
	return d.KeywordError(msg)
}

func (d Directive) BodyErrorIndex(msg string, i uint) *jerr.JAPIError {
	return jerr.NewJAPIError(msg, d.BodyCoords.File(), d.BodyCoords.b+bytes.Index(i))
}

func (d Directive) ParameterError(msg string) *jerr.JAPIError {
	return jerr.NewJAPIError(msg, d.keywordCoords.File(), d.keywordCoords.b)
}

package directive

import (
	"github.com/jsightapi/jsight-schema-core/fs"

	"github.com/jsightapi/jsight-api-core/jerr"

	"github.com/jsightapi/jsight-schema-core/bytes"
)

func (d Directive) KeywordError(msg string) *jerr.JApiError {
	return d.makeError(msg, d.keywordCoords.File(), d.keywordCoords.begin)
}

func (d Directive) BodyError(msg string) *jerr.JApiError {
	if d.BodyCoords.IsSet() {
		return d.makeError(msg, d.BodyCoords.File(), d.BodyCoords.begin)
	}
	return d.KeywordError(msg)
}

func (d Directive) BodyErrorIndex(msg string, i uint) *jerr.JApiError {
	return d.makeError(msg, d.BodyCoords.File(), d.BodyCoords.begin+bytes.Index(i))
}

func (d Directive) ParameterError(msg string) *jerr.JApiError {
	return d.KeywordError(msg)
}

func (d Directive) makeError(msg string, file *fs.File, begin bytes.Index) *jerr.JApiError {
	je := jerr.NewJApiError(msg, file, begin)
	d.includeTracer.AddIncludeTraceToError(je)
	return je
}

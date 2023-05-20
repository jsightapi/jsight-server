package core

import (
	"github.com/jsightapi/jsight-schema-core/bytes"

	"github.com/jsightapi/jsight-api-core/jerr"
)

func (core *JApiCore) japiError(msg string, i bytes.Index) *jerr.JApiError {
	return jerr.NewJApiError(msg, core.scanner.File(), i)
}

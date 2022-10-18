package core

import (
	"github.com/jsightapi/jsight-schema-go-library/bytes"

	"github.com/jsightapi/jsight-api-go-library/jerr"
)

func (core *JApiCore) japiError(msg string, i bytes.Index) *jerr.JApiError {
	return jerr.NewJApiError(msg, core.scanner.File(), i)
}

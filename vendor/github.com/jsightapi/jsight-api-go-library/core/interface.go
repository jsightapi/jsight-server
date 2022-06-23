package core

import (
	"github.com/jsightapi/jsight-api-go-library/catalog"
	"github.com/jsightapi/jsight-api-go-library/jerr"
)

// ValidateJAPI should be used to check if .jst file is valid according to specification
func (core *JApiCore) ValidateJAPI() *jerr.JApiError {
	return core.processJApiProject()
}

func (core *JApiCore) Catalog() *catalog.Catalog {
	return core.catalog
}

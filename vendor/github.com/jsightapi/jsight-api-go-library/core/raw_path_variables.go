package core

import (
	"github.com/jsightapi/jsight-api-go-library/catalog"
	"github.com/jsightapi/jsight-api-go-library/directive"
)

type rawPathVariable struct {
	schema          catalog.Schema
	parameters      []PathParameter
	pathDirective   directive.Directive // to detect and display an error
	parentDirective directive.Directive
}

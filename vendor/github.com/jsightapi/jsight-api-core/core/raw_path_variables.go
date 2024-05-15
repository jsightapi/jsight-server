package core

import (
	"github.com/jsightapi/jsight-schema-core/notations/jschema"

	"github.com/jsightapi/jsight-api-core/directive"
)

type rawPathVariable struct {
	schema          *jschema.JSchema
	parameters      []PathParameter
	pathDirective   directive.Directive // to detect and display an error
	parentDirective directive.Directive

	// temp workaround. true means that this was not gathered from Path directive,
	// but from URL or Method-directive, imitating real rawPathVariable
	imitated bool
}

package catalog

import (
	"github.com/jsightapi/jsight-api-go-library/directive"
)

type UserType struct {
	Annotation  string              `json:"annotation,omitempty"`
	Description string              `json:"description,omitempty"`
	Schema      Schema              `json:"schema"`
	Directive   directive.Directive `json:"-"`
}

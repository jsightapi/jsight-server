package catalog

import (
	"github.com/jsightapi/jsight-api-go-library/directive"
)

type UserRule struct {
	Annotation  string               `json:"annotation"`
	Description string               `json:"description"`
	Directive   *directive.Directive `json:"-"`
	Value       Rule                 `json:"value"`
}

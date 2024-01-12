package catalog

import (
	"github.com/jsightapi/jsight-api-core/directive"
)

type UserRule struct {
	Annotation  string               `json:"annotation"`
	Description string               `json:"description"`
	Directive   *directive.Directive `json:"-"`
	Value       Rule                 `json:"value"`
}

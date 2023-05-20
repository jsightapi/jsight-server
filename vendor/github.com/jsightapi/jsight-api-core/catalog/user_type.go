package catalog

import (
	"github.com/jsightapi/jsight-api-core/directive"
)

type UserType struct {
	Annotation  string              `json:"annotation,omitempty"`
	Description string              `json:"description,omitempty"`
	Schema      ExchangeSchema      `json:"schema"`
	Directive   directive.Directive `json:"-"`
}

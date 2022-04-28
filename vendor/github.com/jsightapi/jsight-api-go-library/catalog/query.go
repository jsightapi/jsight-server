package catalog

import (
	"github.com/jsightapi/jsight-api-go-library/directive"
)

type Query struct {
	Example   string              `json:"example,omitempty"`
	Format    string              `json:"format"`
	Schema    *Schema             `json:"schema"`
	Directive directive.Directive `json:"-"`
}

func NewQuery(d directive.Directive) Query {
	return Query{
		Directive: d,
	}
}

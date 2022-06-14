package catalog

import (
	"github.com/jsightapi/jsight-api-go-library/directive"
)

type HTTPRequest struct {
	*HTTPRequestHeaders `json:"headers,omitempty"`
	*HTTPRequestBody    `json:"body,omitempty"`
	Directive           directive.Directive `json:"-"`
}

type HTTPRequestHeaders struct {
	Schema    *Schema             `json:"schema"`
	Directive directive.Directive `json:"-"`
}

type HTTPRequestBody struct {
	Format    SerializeFormat     `json:"format"`
	Schema    *Schema             `json:"schema"`
	Directive directive.Directive `json:"-"`
}

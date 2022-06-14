package catalog

import (
	"github.com/jsightapi/jsight-api-go-library/directive"
)

type HTTPResponse struct {
	Code       string               `json:"code"`
	Annotation string               `json:"annotation,omitempty"`
	Headers    *HTTPResponseHeaders `json:"headers,omitempty"`
	Body       *HTTPResponseBody    `json:"body"`
	Directive  directive.Directive  `json:"-"`
}

type HTTPResponseHeaders struct {
	Schema    *Schema             `json:"schema"`
	Directive directive.Directive `json:"-"`
}

package catalog

import (
	"github.com/jsightapi/jsight-api-core/directive"
)

type HTTPRequest struct {
	*HTTPRequestHeaders `json:"headers,omitempty"`
	*HTTPRequestBody    `json:"body,omitempty"`
	Directive           directive.Directive `json:"-"`
}

type HTTPRequestHeaders struct {
	Schema    *ExchangeJSightSchema `json:"schema"`
	Directive directive.Directive   `json:"-"`
}

type HTTPRequestBody struct {
	Format    SerializeFormat     `json:"format"`
	Schema    ExchangeSchema      `json:"schema"`
	Directive directive.Directive `json:"-"`
}

// func (h *HTTPRequestBody) Validate(value []byte) error {
// 	return h.Schema.Validate(json.New("", value))
// }

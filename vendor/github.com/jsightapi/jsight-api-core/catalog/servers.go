package catalog

import (
	"sync"

	"github.com/jsightapi/jsight-api-core/directive"
)

type Server struct {
	BaseUrlVariables *baseURLVariables `json:"baseUrlVariables,omitempty"`
	Annotation       string            `json:"annotation,omitempty"`
	BaseUrl          string            `json:"baseUrl"`
}

type baseURLVariables struct {
	Schema    *ExchangeJSightSchema `json:"schema"`
	Directive directive.Directive   `json:"-"`
}

// Servers represent available servers.
// gen:OrderedMap
type Servers struct {
	data  map[string]*Server
	order []string
	mx    sync.RWMutex
}

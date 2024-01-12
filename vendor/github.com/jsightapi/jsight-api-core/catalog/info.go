package catalog

import (
	"github.com/jsightapi/jsight-api-core/directive"
)

// Info general info about api
type Info struct {
	Title       string              `json:"title,omitempty"`
	Version     string              `json:"version,omitempty"`
	Description *string             `json:"description,omitempty"`
	Directive   directive.Directive `json:"-"`
}

package openapi

import (
	"github.com/jsightapi/jsight-api-core/catalog"
)

type Components struct {
	Schemas ComponentsSchemas `json:"schemas,omitempty"`
}

func newComponents(c *catalog.Catalog) *Components {
	if hasComponents(c) {
		return &Components{
			Schemas: newSchemas(c.UserTypes),
		}
	}
	return nil
}

// Currently only UserTypes participate in components
func hasComponents(c *catalog.Catalog) bool {
	return c.UserTypes.Len() > 0
}

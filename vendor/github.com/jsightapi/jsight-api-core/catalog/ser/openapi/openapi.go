package openapi

import (
	"errors"
	"fmt"

	"github.com/jsightapi/jsight-api-core/catalog"
)

type OpenAPI struct {
	catalog *catalog.Catalog

	OpenAPI    string      `json:"openapi"`
	Info       *Info       `json:"info"`
	Servers    []Server    `json:"servers,omitempty"`
	Paths      Paths       `json:"paths"`
	Components *Components `json:"components,omitempty"`
}

func NewOpenAPI(c *catalog.Catalog) (oa *OpenAPI, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(panicToString(r))
		}
	}()

	oa = &OpenAPI{
		catalog:    c,
		OpenAPI:    "3.0.3",
		Info:       newInfo(c.Info),
		Servers:    newServers(c.Servers),
		Paths:      newPaths(c),
		Components: newComponents(c),
	}

	return oa, err
}

func panicToString(r any) string {
	return fmt.Sprintf("%v", r)
}

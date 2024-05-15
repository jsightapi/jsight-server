package openapi

import (
	"github.com/jsightapi/jsight-api-core/catalog"
)

// Other properties of OA MediaTypeObject are not used in JSigh
type MediaTypeObject struct {
	Schema schemaObject `json:"schema,omitempty"` // TODO: empty?
}

func mediaTypeObjectForSchema(es catalog.ExchangeSchema) *MediaTypeObject {
	return &MediaTypeObject{
		Schema: schemaObjectFromExchangeSchema(es),
	}
}

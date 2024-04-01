package openapi

import (
	"github.com/jsightapi/jsight-api-core/catalog"
)

// Content is used in Responses and Requests
type Content map[mediaType]*MediaTypeObject

func defaultContent() *Content {
	return contentForAny()
}

func contentForAny() *Content {
	c := make(Content, 1)
	c[MediaTypeRangeAny] = &MediaTypeObject{
		Schema: schemaObjectForAny(),
	}
	return &c
}

// JSight's pseudo-notation empty is expressed via empty OA's content object.
func contentForEmpty() *Content {
	c := make(Content, 0)
	return &c
}

func contentForVariousMediaTypes(schemaObjectsMap map[mediaType][]schemaObject) *Content {
	c := make(Content, 1)
	for mt, schemaObjects := range schemaObjectsMap {
		if len(schemaObjects) == 1 {
			c[mt] = &MediaTypeObject{schemaObjects[0]}
		} else {
			c[mt] = &MediaTypeObject{schemaObjectForAnyOf(schemaObjects)}
		}
	}
	return &c
}

func contentForSchema(f catalog.SerializeFormat, es catalog.ExchangeSchema) *Content {
	c := make(Content)
	mt := formatToMediaType(f)
	c[mt] = mediaTypeObjectForSchema(es)
	return &c
}

// func ContentWithMediaTypeObject(mt mediaType, o *MediaTypeObject) *Content {
// 	c := make(Content)
// 	c[mt] = o
// 	return &c
// }

// func NewContent(f catalog.SerializeFormat, s catalog.ExchangeSchema) *Content {
// 	c := make(Content)
// 	mt := FormatToMediaType(f)
// 	c[mt] = NewMediaTypeObjectFromExchangeSchema(s)
// 	return &c
// }

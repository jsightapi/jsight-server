package openapi

import (
	"github.com/jsightapi/jsight-api-core/catalog"
)

// Only for Response objects. For request headers refer to "parameters"
type ResponseHeaders map[string]*HeaderObject

type headerInfo struct {
	parameterInfo
	contextAnnotation string
}

func makeResponseHeaders(headersArr ...*catalog.HTTPResponseHeaders) ResponseHeaders {
	r := make(ResponseHeaders, 0)

	sortedHeaders := make(map[string][]headerInfo)
	for _, headers := range headersArr {
		if headers == nil {
			continue
		}

		headersRootAnnotation := getSchemaObjectInfo(headers.Schema.JSchema).Annotation()
		headersInfos := getParamInfo(headers.Schema.JSchema)
		for _, hi := range headersInfos {
			hName := hi.name()
			sortedHeaders[hName] = append(sortedHeaders[hName],
				headerInfo{
					hi,
					headersRootAnnotation,
				})
		}
	}

	for name, headerInfos := range sortedHeaders {
		if len(headerInfos) == 1 {
			i := headerInfos[0]
			r[name] = newHeaderObject(
				!i.optional(),
				i.annotation(),
				i.schemaObject(),
			)
		} else {
			r[name] = headerObjectForAnyOf(headerInfos)
		}
	}
	return r
}

func headerObjectForAnyOf(headersInfos []headerInfo) *HeaderObject {
	schemaObjects := make([]schemaObject, 0)
	required := true

	for _, i := range headersInfos {
		if i.optional() {
			required = false
		}
		so := i.schemaObject()
		so.SetDescription(concatenateDescription(i.contextAnnotation, i.annotation()))
		schemaObjects = append(schemaObjects, so)
	}

	return newHeaderObject(
		required, "",
		&schemaObjectAnyOf{schemaObjects, ""},
	)
}

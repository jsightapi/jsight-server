package openapi

import (
	"github.com/jsightapi/jsight-api-core/catalog"
)

func paramsFromJSchema(es *catalog.ExchangeJSightSchema, in parameterLocation) []*ParameterObject {
	r := make([]*ParameterObject, 0)
	if es == nil {
		return r
	}

	paramInfos := getParamInfo(es.JSchema)
	for _, pi := range paramInfos { // TODO:in future may have multiple infos with same name
		po := newParameterObject(
			in,
			pi.name(),
			pi.annotation(),
			!pi.optional(),
			pi.schemaObject(),
		)
		r = append(r, po)
	}
	return r
}

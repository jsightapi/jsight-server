package openapi

import (
	"github.com/jsightapi/jsight-api-core/catalog"
)

func paramsFromJSchema(es *catalog.ExchangeJSightSchema, in parameterLocation) ([]*ParameterObject, Error) {
	r := make([]*ParameterObject, 0)
	if es == nil {
		return r, newErr("cannot convert nil schema into parameters")
	}

	paramInfos, err := getParamInfos(es.JSchema)
	if err != nil {
		return nil, err.wrapWith("unable to convert schema to parameters")
	} else {
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
	}
	return r, nil
}

package openapi

import (
	"github.com/jsightapi/jsight-api-core/catalog"
)

type Operation struct {
	Summary     string             `json:"summary,omitempty"`
	Description string             `json:"description,omitempty"`
	Parameters  []*ParameterObject `json:"parameters,omitempty"`
	RequestBody *RequestBody       `json:"requestBody,omitempty"`
	Responses   *Responses         `json:"responses"`
	Tags        []string           `json:"tags,omitempty"`
}

func newOperation(i *catalog.HTTPInteraction, tags []string) (*Operation, Error) {
	var requestBody *RequestBody
	if i.HttpMethod == catalog.GET || i.HttpMethod == catalog.DELETE {
		requestBody = nil
	} else {
		requestBody = newRequestBody(i.Request)
	}

	params, err := getOperationParams(i)
	if err != nil {
		return nil, err
	}

	responses, err := newResponses(i)
	if err != nil {
		return nil, err
	}

	return &Operation{
		Summary:     processMethodSummary(i.Annotation),
		Description: processMethodDescription(i.Description),
		Parameters:  params,
		RequestBody: requestBody,
		Responses:   responses,
		Tags:        tags,
	}, nil
}

func processMethodDescription(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}

func processMethodSummary(s *string) string {
	if s == nil {
		return ""
	}

	r := *s
	return r
}

func getOperationParams(i *catalog.HTTPInteraction) ([]*ParameterObject, Error) {
	op := make([]*ParameterObject, 0)

	if querySchemaDefined(i) {
		qp, err := paramsFromJSchema(i.Query.Schema, ParameterLocationQuery)
		if err != nil {
			return op, err.wrapWithf(
				"error converting query schema to OpenAPI parameters for interaction (%s %s)",
				i.HttpMethod.String(), i.Path())
		}
		for _, par := range qp {
			par.Style = ParameterStyleDeepObject
			par.Explode = true
			op = append(op, par)
		}
	}

	if headerSchemaDefined(i) {
		hp, err := paramsFromJSchema(i.Request.HTTPRequestHeaders.Schema, ParameterLocationHeader)
		if err != nil {
			return op, err.wrapWithf("error converting headers for interaction (%s %s)", i.HttpMethod.String(), i.Path())
		} else {
			op = append(op, hp...)
		}
	}

	return op, nil
}

func querySchemaDefined(i *catalog.HTTPInteraction) bool {
	return i.Query != nil && i.Query.Schema != nil
}

func headerSchemaDefined(i *catalog.HTTPInteraction) bool {
	return i.Request != nil &&
		i.Request.HTTPRequestHeaders != nil &&
		i.Request.HTTPRequestHeaders.Schema != nil
}

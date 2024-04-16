package openapi

import "github.com/jsightapi/jsight-api-core/catalog"

type PathItem struct {
	Parameters []*ParameterObject `json:"parameters,omitempty"`
	Get        *Operation         `json:"get,omitempty"`
	Put        *Operation         `json:"put,omitempty"`
	Post       *Operation         `json:"post,omitempty"`
	Patch      *Operation         `json:"patch,omitempty"`
	Delete     *Operation         `json:"delete,omitempty"`
}

func newPathItem(i *catalog.HTTPInteraction) (*PathItem, Error) {
	pp, err := getPathParams(i)
	if err != nil {
		return nil, err
	}
	pi := PathItem{
		Parameters: pp,
	}
	return &pi, nil
}

func getPathParams(i *catalog.HTTPInteraction) ([]*ParameterObject, Error) {
	r := make([]*ParameterObject, 0)
	if pathSchemaDefined(i) {
		params, err := paramsFromJSchema(i.PathVariables.Schema, ParameterLocationPath)
		if err != nil {
			return r, err.wrapWithf(
				"error converting path parameters to OpenaAPI parameters for interaction: %s %s",
				i.HttpMethod.String(), i.Path())
		}
		for _, par := range params {
			par.Required = true
			r = append(r, par)
		}
	}
	return r, nil
}

func pathSchemaDefined(i *catalog.HTTPInteraction) bool {
	return i.PathVariables != nil &&
		i.PathVariables.Schema != nil
}

func (pi *PathItem) assignOperation(method catalog.HTTPMethod, o *Operation) {
	switch method {
	case catalog.GET:
		pi.Get = o
	case catalog.PUT:
		pi.Put = o
	case catalog.POST:
		pi.Post = o
	case catalog.PATCH:
		pi.Patch = o
	case catalog.DELETE:
		pi.Delete = o
	default:
		panic("Unsupported method")
	}
}

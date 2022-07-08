package catalog

import "github.com/jsightapi/jsight-api-go-library/directive"

type JsonRpcInteraction struct { //nolint:govet
	Id       string   `json:"id"`
	Protocol Protocol `json:"protocol"`
	PathVal  Path     `json:"path"`
	// PathVariables *PathVariables `json:"pathVariables,omitempty"`
	Method      string         `json:"method"`
	Tags        []TagName      `json:"tags"`
	Description *string        `json:"annotation,omitempty"`
	Annotation  *string        `json:"description,omitempty"`
	Params      *jsonRpcParams `json:"params,omitempty"`
	Result      *jsonRpcResult `json:"result,omitempty"`
}

type jsonRpcParams struct {
	Schema    *Schema             `json:"schema"`
	Directive directive.Directive `json:"-"`
}

type jsonRpcResult struct {
	Schema    *Schema             `json:"schema"`
	Directive directive.Directive `json:"-"`
}

func (j JsonRpcInteraction) Path() Path {
	return j.PathVal
}

func newJsonRpcInteraction(id JsonRpcInteractionId, method string, annotation string, tn TagName) *JsonRpcInteraction {
	j := &JsonRpcInteraction{
		Id:          id.String(),
		Protocol:    JsonRpc,
		PathVal:     id.path,
		Method:      method,
		Tags:        []TagName{tn},
		Annotation:  nil,
		Description: nil,
		Params:      nil,
		Result:      nil,
	}
	if annotation != "" {
		j.Annotation = &annotation
	}
	return j
}

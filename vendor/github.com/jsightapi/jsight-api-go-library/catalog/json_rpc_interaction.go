package catalog

import "github.com/jsightapi/jsight-api-go-library/directive"

type JsonRpcInteraction struct { //nolint:govet
	Id          string         `json:"id"`
	Protocol    Protocol       `json:"protocol"`
	PathVal     Path           `json:"path"`
	Method      string         `json:"method"`
	Tags        []TagName      `json:"tags"`
	Annotation  *string        `json:"annotation,omitempty"`
	Description *string        `json:"description,omitempty"`
	Params      *jsonRpcParams `json:"params,omitempty"`
	Result      *jsonRpcResult `json:"result,omitempty"`
	// PathVariables *PathVariables `json:"pathVariables,omitempty"`
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

func (j *JsonRpcInteraction) appendTagName(tn TagName) {
	j.Tags = append(j.Tags, tn)
}

func newJsonRpcInteraction(id JsonRpcInteractionId, method string, annotation string) *JsonRpcInteraction {
	j := &JsonRpcInteraction{
		Id:          id.String(),
		Protocol:    JsonRpc,
		PathVal:     id.path,
		Method:      method,
		Tags:        make([]TagName, 0, 3),
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

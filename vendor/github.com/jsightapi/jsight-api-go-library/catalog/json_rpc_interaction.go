package catalog

import "github.com/jsightapi/jsight-api-go-library/directive"

type JsonRpcInteraction struct {
	Description   *string
	Params        *jsonRpcParams
	Result        *jsonRpcResult
	PathVariables *PathVariables
	Annotation    *string
	Protocol      Protocol
	Path          Path
	Method        string
	Tags          []TagName
}

func (ji JsonRpcInteraction) path() Path {
	return ji.Path
}

type jsonRpcParams struct {
	Schema    *Schema             `json:"schema"`
	Directive directive.Directive `json:"-"`
}

type jsonRpcResult struct {
	Schema    *Schema             `json:"schema"`
	Directive directive.Directive `json:"-"`
}

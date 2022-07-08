package catalog

type HttpInteraction struct { //nolint:govet
	Id            string         `json:"id"`
	Protocol      Protocol       `json:"protocol"`
	HttpMethod    HttpMethod     `json:"httpMethod"`
	PathVal       Path           `json:"path"`
	PathVariables *PathVariables `json:"pathVariables,omitempty"`
	Tags          []TagName      `json:"tags"`
	Annotation    *string        `json:"annotation,omitempty"`
	Description   *string        `json:"description,omitempty"`
	Query         *Query         `json:"query,omitempty"`
	Request       *HTTPRequest   `json:"request,omitempty"`
	Responses     []HTTPResponse `json:"responses,omitempty"`
}

func (h HttpInteraction) Path() Path {
	return h.PathVal
}

func (h *HttpInteraction) SetPathVariables(p *PathVariables) {
	h.PathVariables = p
}

func newHttpInteraction(id HttpInteractionId, annotation string, tn TagName) *HttpInteraction {
	h := &HttpInteraction{
		Id:            id.String(),
		Protocol:      HTTP,
		HttpMethod:    id.method,
		PathVal:       id.path,
		Tags:          []TagName{tn},
		PathVariables: nil,
		Annotation:    nil,
		Description:   nil,
		Query:         nil,
		Request:       nil,
		Responses:     []HTTPResponse{},
	}
	if annotation != "" {
		h.Annotation = &annotation
	}
	return h
}

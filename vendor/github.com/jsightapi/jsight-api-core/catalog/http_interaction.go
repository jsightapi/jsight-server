package catalog

type HTTPInteraction struct {
	Id            string         `json:"id"`
	Protocol      Protocol       `json:"protocol"`
	HttpMethod    HTTPMethod     `json:"httpMethod"`
	PathVal       Path           `json:"path"`
	PathVariables *PathVariables `json:"pathVariables,omitempty"`
	Tags          []TagName      `json:"tags"`
	Annotation    *string        `json:"annotation,omitempty"`
	Description   *string        `json:"description,omitempty"`
	Query         *Query         `json:"query,omitempty"`
	Request       *HTTPRequest   `json:"request,omitempty"`
	Responses     []HTTPResponse `json:"responses,omitempty"`
}

func (h HTTPInteraction) Path() Path {
	return h.PathVal
}

func (h *HTTPInteraction) appendTagName(tn TagName) {
	h.Tags = append(h.Tags, tn)
}

func (h *HTTPInteraction) SetPathVariables(p *PathVariables) {
	h.PathVariables = p
}

func newHTTPInteraction(id HTTPInteractionID, annotation string) *HTTPInteraction {
	h := &HTTPInteraction{
		Id:            id.String(),
		Protocol:      HTTP,
		HttpMethod:    id.method,
		PathVal:       id.path,
		Tags:          make([]TagName, 0, 3),
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

package openapi

import "github.com/jsightapi/jsight-api-core/catalog"

type Responses map[responseCode]*ResponseObject

type responseCode string

func defaultResponses() *Responses {
	r := make(Responses, 1)
	r["default"] = defaultResponse()
	return &r
}

func newResponses(i *catalog.HTTPInteraction) (*Responses, Error) {
	if len(i.Responses) == 0 {
		return defaultResponses(), nil
	}

	sortedResponses := make(map[responseCode][]*catalog.HTTPResponse)
	for idx, resp := range i.Responses {
		rCode := responseCode(resp.Code)
		sortedResponses[rCode] = append(sortedResponses[rCode], &i.Responses[idx])
	}

	r := make(Responses, 1)
	for rc, respArr := range sortedResponses {
		var err Error
		var resp *ResponseObject

		if len(respArr) == 1 {
			resp, err = newResponse(respArr[0])
		} else {
			resp, err = newResponseAnyOf(respArr)
		}

		if err != nil {
			return nil, err
		}
		r[rc] = resp
	}
	return &r, nil
}

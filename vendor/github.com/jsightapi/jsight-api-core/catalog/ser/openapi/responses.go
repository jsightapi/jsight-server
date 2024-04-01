package openapi

import "github.com/jsightapi/jsight-api-core/catalog"

type Responses map[responseCode]*ResponseObject

type responseCode string

func defaultResponses() *Responses {
	r := make(Responses, 1)
	r["default"] = defaultResponse()
	return &r
}

func newResponses(i *catalog.HTTPInteraction) *Responses {
	if len(i.Responses) == 0 {
		return defaultResponses()
	}

	sortedResponses := make(map[responseCode][]*catalog.HTTPResponse)
	for idx, resp := range i.Responses {
		rCode := responseCode(resp.Code)
		sortedResponses[rCode] = append(sortedResponses[rCode], &i.Responses[idx])
	}

	r := make(Responses, 1)
	for rc, respArr := range sortedResponses {
		if len(respArr) == 1 {
			r[rc] = newResponse(respArr[0])
		} else {
			r[rc] = newResponseAnyOf(respArr)
		}
	}
	return &r
}

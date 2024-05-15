package core

import (
	"fmt"

	"github.com/jsightapi/jsight-api-core/catalog"
	"github.com/jsightapi/jsight-api-core/jerr"
)

func (core *JApiCore) validateCatalog() *jerr.JApiError {
	if je := core.validateInfo(); je != nil {
		return je
	}

	if je := core.validateRequestBody(); je != nil {
		return je
	}

	if je := core.validateResponseBody(); je != nil {
		return je
	}

	if je := core.validateHeaders(); je != nil {
		return je
	}

	return nil
}

func (core *JApiCore) validateInfo() *jerr.JApiError {
	isEmpty := core.catalog.Info != nil &&
		core.catalog.Info.Title == "" &&
		core.catalog.Info.Version == "" &&
		core.catalog.Info.Description == nil

	if isEmpty {
		return core.catalog.Info.Directive.KeywordError(jerr.InfoIsEmpty)
	}
	return nil
}

func (core *JApiCore) validateRequestBody() *jerr.JApiError {
	return adoptError(core.catalog.Interactions.Each(func(k catalog.InteractionID, v catalog.Interaction) error {
		if hi, ok := v.(*catalog.HTTPInteraction); ok {
			r := hi.Request
			if r != nil && r.HTTPRequestBody == nil {
				return r.Directive.KeywordError(
					fmt.Sprintf(`%s %q`, jerr.UndefinedRequestBodyForResource, k.String()),
				)
			}
		}
		return nil
	}))
}

func (core *JApiCore) validateResponseBody() *jerr.JApiError {
	return adoptError(core.catalog.Interactions.Each(func(k catalog.InteractionID, v catalog.Interaction) error {
		if hi, ok := v.(*catalog.HTTPInteraction); ok {
			for _, response := range hi.Responses {
				if response.Body == nil {
					return response.Directive.KeywordError(fmt.Sprintf(
						"undefined response body for resource %q, HTTP-code %q",
						k.String(),
						response.Code,
					))
				}
			}
		}
		return nil
	}))
}

func (core *JApiCore) isJsightCastToObject(schema catalog.ExchangeSchema) bool {
	if schema == nil {
		return false
	}

	s, ok := schema.(*catalog.ExchangeJSightSchema)
	if !ok {
		return false
	}

	return s.CastToObject() != nil
}

func (core *JApiCore) validateHeaders() *jerr.JApiError {
	return adoptError(core.catalog.Interactions.Each(func(_ catalog.InteractionID, v catalog.Interaction) error {
		if hi, ok := v.(*catalog.HTTPInteraction); ok {
			request := hi.Request
			isNotAnObject := request != nil &&
				request.HTTPRequestHeaders != nil &&
				!core.isJsightCastToObject(request.HTTPRequestHeaders.Schema)
			if isNotAnObject {
				return request.HTTPRequestHeaders.Directive.BodyError(jerr.BodyMustBeObject)
			}
			for _, response := range hi.Responses {
				if response.Headers != nil && !core.isJsightCastToObject(response.Headers.Schema) {
					return response.Headers.Directive.BodyError(jerr.BodyMustBeObject)
				}
			}
		}
		return nil
	}))
}

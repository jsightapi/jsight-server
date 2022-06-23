package core

import (
	"fmt"

	"github.com/jsightapi/jsight-api-go-library/catalog"
	"github.com/jsightapi/jsight-api-go-library/jerr"
	"github.com/jsightapi/jsight-api-go-library/notation"
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

	return core.validateUsedUserTypes()
}

func (core *JApiCore) validateInfo() *jerr.JApiError {
	if core.catalog.Info != nil && core.catalog.Info.Title == "" && core.catalog.Info.Version == "" && core.catalog.Info.Description == nil {
		return core.catalog.Info.Directive.KeywordError("empty info")
	}
	return nil
}

func (core *JApiCore) validateUsedUserTypes() *jerr.JApiError {
	err := core.catalog.UserTypes.Each(func(k string, v *catalog.UserType) error {
		if err := core.findUserTypes(v.Schema.UsedUserTypes); err != nil {
			return v.Directive.BodyError(err.Error())
		}
		return nil
	})
	if err != nil {
		return adoptError(err)
	}

	err = core.catalog.Servers.Each(func(k string, v *catalog.Server) error {
		s := v.BaseUrlVariables
		if s != nil && s.Schema != nil {
			if err := core.findUserTypes(s.Schema.UsedUserTypes); err != nil {
				return s.Directive.BodyError(err.Error())
			}
		}
		return nil
	})
	if err != nil {
		return adoptError(err)
	}

	return adoptError(core.catalog.HttpInteractions.Each(func(k catalog.HttpInteractionId, v *catalog.HttpInteraction) error {
		if v.Query != nil && v.Query.Schema != nil {
			if err := core.findUserTypes(v.Query.Schema.UsedUserTypes); err != nil {
				return v.Query.Directive.BodyError(err.Error())
			}
		}

		if v.Request != nil {
			if v.Request.HTTPRequestHeaders != nil && v.Request.HTTPRequestHeaders.Schema != nil {
				if err := core.findUserTypes(v.Request.HTTPRequestHeaders.Schema.UsedUserTypes); err != nil {
					return v.Request.HTTPRequestHeaders.Directive.BodyError(err.Error())
				}
			}

			if v.Request.HTTPRequestBody != nil && v.Request.HTTPRequestBody.Schema != nil {
				if err := core.findUserTypes(v.Request.HTTPRequestBody.Schema.UsedUserTypes); err != nil {
					return v.Request.HTTPRequestBody.Directive.BodyError(err.Error())
				}
			}
		}

		for _, r := range v.Responses {
			if r.Headers != nil && r.Headers.Schema != nil {
				if err := core.findUserTypes(r.Headers.Schema.UsedUserTypes); err != nil {
					return r.Headers.Directive.BodyError(err.Error())
				}
			}

			if r.Body != nil && r.Body.Schema != nil {
				if err := core.findUserTypes(r.Body.Schema.UsedUserTypes); err != nil {
					return r.Body.Directive.BodyError(err.Error())
				}
			}
		}
		return nil
	}))
}

// findUserTypes returns an error if a user type cannot be found
func (core *JApiCore) findUserTypes(uu *catalog.StringSet) error {
	for _, u := range uu.Data() {
		if !core.catalog.UserTypes.Has(u) {
			return fmt.Errorf("user type %q wasn't found", u)
		}
	}
	return nil
}

func (core *JApiCore) validateRequestBody() *jerr.JApiError {
	return adoptError(core.catalog.HttpInteractions.Each(func(k catalog.HttpInteractionId, v *catalog.HttpInteraction) error {
		r := v.Request
		if r != nil && r.HTTPRequestBody == nil {
			return r.Directive.KeywordError(fmt.Sprintf(`undefined request body for resource "%s"`, k.String()))
		}
		return nil
	}))
}

func (core *JApiCore) validateResponseBody() *jerr.JApiError {
	return adoptError(core.catalog.HttpInteractions.Each(func(k catalog.HttpInteractionId, v *catalog.HttpInteraction) error {
		for _, response := range v.Responses {
			if response.Body == nil {
				return response.Directive.KeywordError(fmt.Sprintf(`undefined response body for resource "%s", HTTP-code "%s"`, k.String(), response.Code))
			}
		}
		return nil
	}))
}

func (core *JApiCore) isJsightCastToObject(schema *catalog.Schema) bool {
	if schema != nil && schema.ContentJSight != nil && schema.Notation == notation.SchemaNotationJSight {
		switch schema.ContentJSight.TokenType {
		case "object":
			return true
		case "reference":
			if userType, ok := core.catalog.UserTypes.Get(schema.ContentJSight.ScalarValue); ok {
				return core.isJsightCastToObject(&userType.Schema)
			}
		}
	}
	return false
}

func (core *JApiCore) validateHeaders() *jerr.JApiError {
	return adoptError(core.catalog.HttpInteractions.Each(func(_ catalog.HttpInteractionId, v *catalog.HttpInteraction) error {
		request := v.Request
		if request != nil && request.HTTPRequestHeaders != nil && !core.isJsightCastToObject(request.HTTPRequestHeaders.Schema) {
			return v.Request.HTTPRequestHeaders.Directive.BodyError(jerr.BodyMustBeObject)
		}
		for _, response := range v.Responses {
			if response.Headers != nil && !core.isJsightCastToObject(response.Headers.Schema) {
				return response.Headers.Directive.BodyError(jerr.BodyMustBeObject)
			}
		}
		return nil
	}))
}

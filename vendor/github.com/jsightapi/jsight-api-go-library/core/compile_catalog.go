package core

import (
	"errors"
	"fmt"
	"strings"

	jschema "github.com/jsightapi/jsight-schema-go-library"

	"github.com/jsightapi/jsight-api-go-library/catalog"
	"github.com/jsightapi/jsight-api-go-library/jerr"
	"github.com/jsightapi/jsight-api-go-library/notation"
)

func (core *JApiCore) compileCatalog() *jerr.JApiError {
	if je := core.ProcessAllOf(); je != nil {
		return je
	}

	if je := core.ExpandRawPathVariableShortcuts(); je != nil {
		return je
	}

	if je := core.CheckRawPathVariableSchemas(); je != nil {
		return je
	}

	return core.BuildResourceMethodsPathVariables()
}

func (core *JApiCore) ExpandRawPathVariableShortcuts() *jerr.JApiError {
	for i := 0; i < len(core.rawPathVariables); i++ {
		r := &core.rawPathVariables[i]

		for r.schema.ContentJSight.TokenType == jschema.TokenTypeShortcut {
			typeName := r.schema.ContentJSight.Type
			if typeName == "mixed" {
				return r.pathDirective.KeywordError("The root schema object cannot have an OR rule")
			}

			ut, ok := core.catalog.UserTypes.Get(typeName)
			if !ok {
				return r.pathDirective.KeywordError(fmt.Sprintf(`User type "%s" not found`, typeName))
			}

			r.schema = ut.Schema // copy schema
		}

		if err := checkPathSchema(r.schema); err != nil {
			return r.pathDirective.KeywordError(err.Error())
		}
	}

	return nil
}

func (core *JApiCore) CheckRawPathVariableSchemas() *jerr.JApiError {
	for i := 0; i < len(core.rawPathVariables); i++ {
		if err := checkPathSchema(core.rawPathVariables[i].schema); err != nil {
			return core.rawPathVariables[i].pathDirective.KeywordError(err.Error())
		}
	}
	return nil
}

func checkPathSchema(s catalog.Schema) error {
	if s.ContentJSight.TokenType != jschema.TokenTypeObject {
		return errors.New("the body of the Path DIRECTIVE must be an object")
	}

	if s.ContentJSight.Rules.Has("additionalProperties") {
		return errors.New(`the "additionalProperties" rule is invalid in the Path directive`)
	}

	if s.ContentJSight.Rules.Has("nullable") {
		return errors.New(`the "nullable" rule is invalid in the Path directive`)
	}

	if s.ContentJSight.Rules.Has("or") {
		return errors.New(`the "or" rule is invalid in the Path directive`)
	}

	if s.ContentJSight.Children == nil || len(s.ContentJSight.Children) == 0 {
		return errors.New("an empty object in the Path directive")
	}

	for _, v := range s.ContentJSight.Children {
		if v.TokenType == jschema.TokenTypeObject || v.TokenType == jschema.TokenTypeArray {
			return fmt.Errorf("the multi-level property %q is not allowed in the Path directive", *(v.Key))
		}
	}

	return nil
}

func (core *JApiCore) BuildResourceMethodsPathVariables() *jerr.JApiError {
	allProjectProperties := make(map[catalog.Path]prop)
	for _, v := range core.rawPathVariables {
		pp := core.propertiesToMap(v.schema.ContentJSight.Children)

		for _, p := range v.parameters {
			if sc, ok := pp[p.parameter]; ok {
				if _, ok := allProjectProperties[p.path]; ok {
					return v.pathDirective.KeywordError(fmt.Sprintf(
						"The parameter %q has already been defined earlier",
						p.parameter,
					))
				}

				allProjectProperties[p.path] = prop{
					schemaContentJSight: sc,
					directive:           v.pathDirective,
				}

				delete(pp, p.parameter)
			}
		}

		// Check that all path properties in schema is exists in the path.
		if len(pp) > 0 {
			ss := core.getPropertiesNames(pp)
			return v.pathDirective.KeywordError(fmt.Sprintf("Has unused parameters %q in schema", ss))
		}
	}

	// set PathVariables
	err := core.catalog.Interactions.Map(
		func(_ catalog.InteractionID, v catalog.Interaction) (catalog.Interaction, error) {
			if hi, ok := v.(*catalog.HTTPInteraction); ok {
				pp := pathParameters(v.Path().String())
				properties := make([]prop, 0, len(pp))

				for _, p := range pp {
					if pr, ok := allProjectProperties[p.path]; ok {
						pr.parameter = p.parameter
						properties = append(properties, pr)
					}
				}

				if len(properties) != 0 {
					pv, err := core.newPathVariables(properties)
					if err != nil {
						return nil, err
					}
					hi.SetPathVariables(pv)
				}
			}
			return v, nil
		},
	)
	if err != nil {
		return err.(*jerr.JApiError) //nolint:errorlint
	}

	return nil
}

func (*JApiCore) propertiesToMap(pp []*catalog.SchemaContentJSight) map[string]*catalog.SchemaContentJSight {
	if len(pp) == 0 {
		return nil
	}

	res := make(map[string]*catalog.SchemaContentJSight, len(pp))
	for _, v := range pp {
		res[*(v.Key)] = v
	}
	return res
}

func (*JApiCore) getPropertiesNames(pp map[string]*catalog.SchemaContentJSight) string {
	if len(pp) == 0 {
		return ""
	}

	buf := strings.Builder{}
	for k := range pp {
		buf.WriteString(k)
		buf.WriteString(", ")
	}
	return strings.TrimSuffix(buf.String(), ", ")
}

func (core *JApiCore) ProcessAllOf() *jerr.JApiError {
	if je := core.processUserTypes(); je != nil {
		return je
	}

	if je := core.processBaseUrlAllOf(); je != nil {
		return je
	}

	if je := core.processRawPathVariablesAllOf(); je != nil {
		return je
	}

	if je := core.processQueryAllOf(); je != nil {
		return je
	}

	if je := core.processRequestHeaderAllOf(); je != nil {
		return je
	}

	if je := core.processRequestAllOf(); je != nil {
		return je
	}

	if je := core.processResponseHeaderAllOf(); je != nil {
		return je
	}

	return core.processResponseAllOf()
}

func (core *JApiCore) processUserTypes() *jerr.JApiError {
	return adoptError(core.catalog.UserTypes.Each(func(k string, v *catalog.UserType) error {
		if v.Schema.Notation == notation.SchemaNotationJSight {
			if err := core.processSchemaContentJSightAllOf(v.Schema.ContentJSight, v.Schema.UsedUserTypes); err != nil {
				return v.Directive.BodyError(err.Error())
			}
		}
		return nil
	}))
}

func (core *JApiCore) processBaseUrlAllOf() *jerr.JApiError {
	return adoptError(core.catalog.Servers.Each(func(k string, v *catalog.Server) error {
		s := v.BaseUrlVariables
		if s != nil && s.Schema != nil && s.Schema.Notation == notation.SchemaNotationJSight {
			if err := core.processSchemaContentJSightAllOf(s.Schema.ContentJSight, s.Schema.UsedUserTypes); err != nil {
				return s.Directive.BodyError(err.Error())
			}
		}
		return nil
	}))
}

func (core *JApiCore) processRawPathVariablesAllOf() *jerr.JApiError {
	for _, r := range core.rawPathVariables {
		if r.schema.Notation == notation.SchemaNotationJSight {
			if err := core.processSchemaContentJSightAllOf(r.schema.ContentJSight, r.schema.UsedUserTypes); err != nil {
				return r.pathDirective.BodyError(err.Error())
			}
		}
	}
	return nil
}

func (core *JApiCore) processQueryAllOf() *jerr.JApiError {
	return adoptError(core.catalog.Interactions.Each(func(_ catalog.InteractionID, v catalog.Interaction) error {
		if hi, ok := v.(*catalog.HTTPInteraction); ok {
			q := hi.Query
			if q != nil && q.Schema != nil && q.Schema.Notation == notation.SchemaNotationJSight {
				err := core.processSchemaContentJSightAllOf(q.Schema.ContentJSight, q.Schema.UsedUserTypes)
				if err != nil {
					return q.Directive.BodyError(err.Error())
				}
			}
		}
		return nil
	}))
}

func (core *JApiCore) processRequestHeaderAllOf() *jerr.JApiError {
	return adoptError(core.catalog.Interactions.Each(func(_ catalog.InteractionID, v catalog.Interaction) error {
		if hi, ok := v.(*catalog.HTTPInteraction); ok {
			r := hi.Request
			isJSight := r != nil &&
				r.HTTPRequestHeaders != nil &&
				r.HTTPRequestHeaders.Schema != nil &&
				r.HTTPRequestHeaders.Schema.Notation == notation.SchemaNotationJSight
			if isJSight {
				h := r.HTTPRequestHeaders
				err := core.processSchemaContentJSightAllOf(h.Schema.ContentJSight, h.Schema.UsedUserTypes)
				if err != nil {
					return r.HTTPRequestHeaders.Directive.BodyError(err.Error())
				}
			}
		}
		return nil
	}))
}

func (core *JApiCore) processRequestAllOf() *jerr.JApiError {
	return adoptError(core.catalog.Interactions.Each(func(_ catalog.InteractionID, v catalog.Interaction) error {
		if hi, ok := v.(*catalog.HTTPInteraction); ok {
			r := hi.Request
			isJSight := r != nil &&
				r.HTTPRequestBody != nil &&
				r.HTTPRequestBody.Schema != nil &&
				r.HTTPRequestBody.Schema.Notation == notation.SchemaNotationJSight
			if isJSight {
				b := r.HTTPRequestBody
				err := core.processSchemaContentJSightAllOf(b.Schema.ContentJSight, b.Schema.UsedUserTypes)
				if err != nil {
					return r.HTTPRequestBody.Directive.BodyError(err.Error())
				}
			}
		}
		return nil
	}))
}

func (core *JApiCore) processResponseHeaderAllOf() *jerr.JApiError {
	return adoptError(core.catalog.Interactions.Each(func(_ catalog.InteractionID, v catalog.Interaction) error {
		if hi, ok := v.(*catalog.HTTPInteraction); ok {
			for _, resp := range hi.Responses {
				h := resp.Headers
				if h != nil && h.Schema != nil && h.Schema.Notation == notation.SchemaNotationJSight {
					err := core.processSchemaContentJSightAllOf(h.Schema.ContentJSight, h.Schema.UsedUserTypes)
					if err != nil {
						return resp.Headers.Directive.BodyError(err.Error())
					}
				}
			}
		}
		return nil
	}))
}

func (core *JApiCore) processResponseAllOf() *jerr.JApiError {
	return adoptError(core.catalog.Interactions.Each(func(_ catalog.InteractionID, v catalog.Interaction) error {
		if hi, ok := v.(*catalog.HTTPInteraction); ok {
			for _, resp := range hi.Responses {
				b := resp.Body
				isJSight := b != nil &&
					b.Schema != nil &&
					b.Schema.Notation == notation.SchemaNotationJSight
				if isJSight {
					err := core.processSchemaContentJSightAllOf(b.Schema.ContentJSight, b.Schema.UsedUserTypes)
					if err != nil {
						return resp.Body.Directive.BodyError(err.Error())
					}
				}
			}
		}
		return nil
	}))
}

func (core *JApiCore) processSchemaContentJSightAllOf(sc *catalog.SchemaContentJSight, uut *catalog.StringSet) error {
	if sc.TokenType != jschema.TokenTypeObject {
		return nil
	}

	for _, v := range sc.Children {
		if err := core.processSchemaContentJSightAllOf(v, uut); err != nil {
			return err
		}
	}

	rule, ok := sc.Rules.Get("allOf")
	if !ok {
		return nil
	}

	switch rule.TokenType { //nolint:exhaustive // We expects only this types.
	case catalog.RuleTokenTypeArray:
		for i := len(rule.Children) - 1; i >= 0; i-- {
			r := rule.Children[i]
			if err := core.inheritPropertiesFromUserType(sc, uut, r.ScalarValue); err != nil {
				return err
			}
		}
	case catalog.RuleTokenTypeReference:
		if err := core.inheritPropertiesFromUserType(sc, uut, rule.ScalarValue); err != nil {
			return err
		}
	}
	return nil
}

func (core *JApiCore) inheritPropertiesFromUserType(
	sc *catalog.SchemaContentJSight,
	uut *catalog.StringSet,
	userTypeName string,
) error {
	ut, ok := core.catalog.UserTypes.Get(userTypeName)
	if !ok {
		return fmt.Errorf(`the user type %q not found`, userTypeName)
	}

	if ut.Schema.ContentJSight.TokenType != jschema.TokenTypeObject {
		return fmt.Errorf(`the user type %q is not an object`, userTypeName)
	}

	if _, ok := core.processedByAllOf[userTypeName]; !ok {
		core.processedByAllOf[userTypeName] = struct{}{}
		if err := core.processSchemaContentJSightAllOf(ut.Schema.ContentJSight, uut); err != nil {
			return err
		}
	}

	if sc.Children == nil {
		sc.Children = make([]*catalog.SchemaContentJSight, 0, 10)
	}

	for i := len(ut.Schema.ContentJSight.Children) - 1; i >= 0; i-- {
		v := ut.Schema.ContentJSight.Children[i]

		if v.Key == nil {
			return fmt.Errorf(jerr.InternalServerError)
		}

		p := sc.ObjectProperty(*(v.Key))
		if p != nil && p.InheritedFrom == "" {
			// Don't allow to override original properties.
			return fmt.Errorf(
				"it is not allowed to override the %q property from the user type %q",
				*(v.Key),
				userTypeName,
			)
		}

		if p != nil && p.InheritedFrom != "" {
			// This property already defined, skip.
			continue
		}

		vv := *v
		if vv.InheritedFrom == "" {
			uut.Add(userTypeName)
		}
		vv.InheritedFrom = userTypeName
		sc.Unshift(&vv)
	}

	return nil
}

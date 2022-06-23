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

		for r.schema.ContentJSight.TokenType == jschema.JSONTypeShortcut {
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
	if s.ContentJSight.TokenType != jschema.JSONTypeObject {
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
		if v.TokenType == jschema.JSONTypeObject || v.TokenType == jschema.JSONTypeArray {
			return fmt.Errorf("the multi-level property %q is not allowed in the Path directive", v.Key)
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
					return v.pathDirective.KeywordError(fmt.Sprintf("The parameter %q has already been defined earlier", p.parameter))
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

	err := core.catalog.HttpInteractions.Map(func(id catalog.HttpInteractionId, resourceMethod *catalog.HttpInteraction) (*catalog.HttpInteraction, error) {
		pp := pathParameters(resourceMethod.Path.String())
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
			resourceMethod.PathVariables = pv
		}
		return resourceMethod, nil
	})
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
		res[v.Key] = v
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
	return adoptError(core.catalog.HttpInteractions.Each(func(_ catalog.HttpInteractionId, v *catalog.HttpInteraction) error {
		q := v.Query
		if q != nil && q.Schema != nil && q.Schema.Notation == notation.SchemaNotationJSight {
			if err := core.processSchemaContentJSightAllOf(q.Schema.ContentJSight, q.Schema.UsedUserTypes); err != nil {
				return q.Directive.BodyError(err.Error())
			}
		}
		return nil
	}))
}

func (core *JApiCore) processRequestHeaderAllOf() *jerr.JApiError {
	return adoptError(core.catalog.HttpInteractions.Each(func(_ catalog.HttpInteractionId, v *catalog.HttpInteraction) error {
		r := v.Request
		if r != nil && r.HTTPRequestHeaders != nil && r.HTTPRequestHeaders.Schema != nil && r.HTTPRequestHeaders.Schema.Notation == notation.SchemaNotationJSight {
			h := r.HTTPRequestHeaders
			if err := core.processSchemaContentJSightAllOf(h.Schema.ContentJSight, h.Schema.UsedUserTypes); err != nil {
				return r.HTTPRequestHeaders.Directive.BodyError(err.Error())
			}
		}
		return nil
	}))
}

func (core *JApiCore) processRequestAllOf() *jerr.JApiError {
	return adoptError(core.catalog.HttpInteractions.Each(func(_ catalog.HttpInteractionId, v *catalog.HttpInteraction) error {
		r := v.Request
		if r != nil && r.HTTPRequestBody != nil && r.HTTPRequestBody.Schema != nil && r.HTTPRequestBody.Schema.Notation == notation.SchemaNotationJSight {
			b := r.HTTPRequestBody
			if err := core.processSchemaContentJSightAllOf(b.Schema.ContentJSight, b.Schema.UsedUserTypes); err != nil {
				return r.HTTPRequestBody.Directive.BodyError(err.Error())
			}
		}
		return nil
	}))
}

func (core *JApiCore) processResponseHeaderAllOf() *jerr.JApiError {
	return adoptError(core.catalog.HttpInteractions.Each(func(_ catalog.HttpInteractionId, v *catalog.HttpInteraction) error {
		for _, resp := range v.Responses {
			h := resp.Headers
			if h != nil && h.Schema != nil && h.Schema.Notation == notation.SchemaNotationJSight {
				if err := core.processSchemaContentJSightAllOf(h.Schema.ContentJSight, h.Schema.UsedUserTypes); err != nil {
					return resp.Headers.Directive.BodyError(err.Error())
				}
			}
		}
		return nil
	}))
}

func (core *JApiCore) processResponseAllOf() *jerr.JApiError {
	return adoptError(core.catalog.HttpInteractions.Each(func(_ catalog.HttpInteractionId, v *catalog.HttpInteraction) error {
		for _, resp := range v.Responses {
			b := resp.Body
			if b != nil && b.Schema != nil && b.Schema.Notation == notation.SchemaNotationJSight {
				if err := core.processSchemaContentJSightAllOf(b.Schema.ContentJSight, b.Schema.UsedUserTypes); err != nil {
					return resp.Body.Directive.BodyError(err.Error())
				}
			}
		}
		return nil
	}))
}

func (core *JApiCore) processSchemaContentJSightAllOf(sc *catalog.SchemaContentJSight, uut *catalog.StringSet) error {
	if sc.TokenType != jschema.JSONTypeObject {
		return nil
	}

	for _, v := range sc.Children {
		if err := core.processSchemaContentJSightAllOf(v, uut); err != nil {
			return err
		}
	}

	if rule, ok := sc.Rules.Get("allOf"); ok {
		switch rule.TokenType {
		case jschema.JSONTypeArray:
			for i := len(rule.Children) - 1; i >= 0; i-- {
				r := rule.Children[i]
				if err := core.inheritPropertiesFromUserType(sc, uut, r.ScalarValue); err != nil {
					return err
				}
			}
		case jschema.JSONTypeString:
			if err := core.inheritPropertiesFromUserType(sc, uut, rule.ScalarValue); err != nil {
				return err
			}
		}
	}
	return nil
}

func (core *JApiCore) inheritPropertiesFromUserType(sc *catalog.SchemaContentJSight, uut *catalog.StringSet, userTypeName string) error {
	ut, ok := core.catalog.UserTypes.Get(userTypeName)
	if !ok {
		return fmt.Errorf(`the user type %q not found`, userTypeName)
	}

	if ut.Schema.ContentJSight.TokenType != jschema.JSONTypeObject {
		return fmt.Errorf(`the user type %q is not an object`, userTypeName)
	}

	if sc.Children == nil {
		sc.Children = make([]*catalog.SchemaContentJSight, 0, 10)
	}

	for i := len(ut.Schema.ContentJSight.Children) - 1; i >= 0; i-- {
		v := ut.Schema.ContentJSight.Children[i]

		if sc.IsObjectHaveProperty(v.Key) {
			return fmt.Errorf(`it is not allowed to override the "%s" property from the user type "%s"`, v.Key, userTypeName)
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

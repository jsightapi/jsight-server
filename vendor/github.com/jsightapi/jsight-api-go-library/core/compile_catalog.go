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

func (core *JApiCore) compileCatalog() *jerr.JAPIError {
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

func (core *JApiCore) ExpandRawPathVariableShortcuts() *jerr.JAPIError {
	for i := 0; i < len(core.rawPathVariables); i++ {
		r := &core.rawPathVariables[i]

		for r.schema.ContentJSight.JsonType == jschema.JSONTypeShortcut {
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

func (core *JApiCore) CheckRawPathVariableSchemas() *jerr.JAPIError {
	for i := 0; i < len(core.rawPathVariables); i++ {
		if err := checkPathSchema(core.rawPathVariables[i].schema); err != nil {
			return core.rawPathVariables[i].pathDirective.KeywordError(err.Error())
		}
	}
	return nil
}

func checkPathSchema(s catalog.Schema) error {
	if s.ContentJSight.JsonType != jschema.JSONTypeObject {
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

	if s.ContentJSight.Properties == nil || s.ContentJSight.Properties.Len() == 0 {
		return errors.New("an empty object in the Path directive")
	}

	kv, ok := s.ContentJSight.Properties.Find(func(_ string, v *catalog.SchemaContentJSight) bool {
		return v.JsonType == jschema.JSONTypeObject || v.JsonType == jschema.JSONTypeArray
	})
	if ok {
		return fmt.Errorf("the multi-level property %q is not allowed in the Path directive", kv.Key)
	}

	return nil
}

func (core *JApiCore) BuildResourceMethodsPathVariables() *jerr.JAPIError {
	allProjectProperties := make(map[catalog.Path]prop)
	for _, v := range core.rawPathVariables {
		pp := core.propertiesToMap(v.schema.ContentJSight.Properties)

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

	err := core.catalog.ResourceMethods.Map(func(id catalog.ResourceMethodId, resourceMethod *catalog.ResourceMethod) (*catalog.ResourceMethod, error) {
		properties := make(map[string]prop)
		pp := pathParameters(resourceMethod.Path.String())

		for _, p := range pp {
			if pr, ok := allProjectProperties[p.path]; ok {
				properties[p.parameter] = pr
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
		return err.(*jerr.JAPIError) //nolint:errorlint
	}

	return nil
}

func (*JApiCore) propertiesToMap(pp *catalog.Properties) map[string]*catalog.SchemaContentJSight {
	if pp == nil || pp.Len() == 0 {
		return nil
	}

	res := make(map[string]*catalog.SchemaContentJSight, pp.Len())
	pp.EachSafe(func(k string, v *catalog.SchemaContentJSight) {
		res[k] = v
	})
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

func (core *JApiCore) ProcessAllOf() *jerr.JAPIError {
	var err *jerr.JAPIError

	err = core.processUserTypes()
	if err != nil {
		return err
	}

	err = core.processBaseUrlAllOf()
	if err != nil {
		return err
	}

	err = core.processRawPathVariablesAllOf()
	if err != nil {
		return err
	}

	err = core.processQueryAllOf()
	if err != nil {
		return err
	}

	err = core.processRequestHeaderAllOf()
	if err != nil {
		return err
	}

	err = core.processRequestAllOf()
	if err != nil {
		return err
	}

	err = core.processResponseHeaderAllOf()
	if err != nil {
		return err
	}

	err = core.processResponseAllOf()
	if err != nil {
		return err
	}

	return nil
}

func (core *JApiCore) processUserTypes() *jerr.JAPIError {
	return adoptError(core.catalog.UserTypes.Each(func(k string, v *catalog.UserType) error {
		if v.Schema.Notation == notation.SchemaNotationJSight {
			if err := core.processSchemaContentJSightAllOf(v.Schema.ContentJSight, v.Schema.UsedUserTypes); err != nil {
				return v.Directive.BodyError(err.Error())
			}
		}
		return nil
	}))
}

func (core *JApiCore) processBaseUrlAllOf() *jerr.JAPIError {
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

func (core *JApiCore) processRawPathVariablesAllOf() *jerr.JAPIError {
	for _, r := range core.rawPathVariables {
		if r.schema.Notation == notation.SchemaNotationJSight {
			if err := core.processSchemaContentJSightAllOf(r.schema.ContentJSight, r.schema.UsedUserTypes); err != nil {
				return r.pathDirective.BodyError(err.Error())
			}
		}
	}
	return nil
}

func (core *JApiCore) processQueryAllOf() *jerr.JAPIError {
	return adoptError(core.catalog.ResourceMethods.Each(func(_ catalog.ResourceMethodId, v *catalog.ResourceMethod) error {
		q := v.Query
		if q != nil && q.Schema != nil && q.Schema.Notation == notation.SchemaNotationJSight {
			if err := core.processSchemaContentJSightAllOf(q.Schema.ContentJSight, q.Schema.UsedUserTypes); err != nil {
				return q.Directive.BodyError(err.Error())
			}
		}
		return nil
	}))
}

func (core *JApiCore) processRequestHeaderAllOf() *jerr.JAPIError {
	return adoptError(core.catalog.ResourceMethods.Each(func(_ catalog.ResourceMethodId, v *catalog.ResourceMethod) error {
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

func (core *JApiCore) processRequestAllOf() *jerr.JAPIError {
	return adoptError(core.catalog.ResourceMethods.Each(func(_ catalog.ResourceMethodId, v *catalog.ResourceMethod) error {
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

func (core *JApiCore) processResponseHeaderAllOf() *jerr.JAPIError {
	return adoptError(core.catalog.ResourceMethods.Each(func(_ catalog.ResourceMethodId, v *catalog.ResourceMethod) error {
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

func (core *JApiCore) processResponseAllOf() *jerr.JAPIError {
	return adoptError(core.catalog.ResourceMethods.Each(func(_ catalog.ResourceMethodId, v *catalog.ResourceMethod) error {
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
	if sc.JsonType != jschema.JSONTypeObject {
		return nil
	}

	err := sc.Properties.Each(func(_ string, v *catalog.SchemaContentJSight) error {
		return core.processSchemaContentJSightAllOf(v, uut)
	})
	if err != nil {
		return err
	}

	if rule, ok := sc.Rules.Get("allOf"); ok {
		switch rule.JsonType {
		case jschema.JSONTypeArray:
			for i := len(rule.Items) - 1; i >= 0; i-- {
				r := rule.Items[i]
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

	if ut.Schema.ContentJSight.JsonType != jschema.JSONTypeObject {
		return fmt.Errorf(`the user type %q is not an object`, userTypeName)
	}

	if sc.Properties == nil {
		sc.Properties = &catalog.Properties{}
	}

	return ut.Schema.ContentJSight.Properties.EachReverse(func(k string, v *catalog.SchemaContentJSight) error {
		if sc.Properties.Has(k) {
			return fmt.Errorf(`it is not allowed to override the "%s" property from the user type "%s"`, k, userTypeName)
		}
		vv := *v
		if vv.InheritedFrom == "" {
			uut.Add(userTypeName)
		}
		vv.InheritedFrom = userTypeName
		sc.Properties.SetToTop(k, &vv)
		return nil
	})
}

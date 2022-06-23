package catalog

import (
	"errors"
	"fmt"

	jschemaLib "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/bytes"
	"github.com/jsightapi/jsight-schema-go-library/kit"

	"github.com/jsightapi/jsight-api-go-library/directive"
	"github.com/jsightapi/jsight-api-go-library/jerr"
	"github.com/jsightapi/jsight-api-go-library/notation"
)

// tag returns a Tag from the collection, or creates a new one and adds it to the collection
func (c *Catalog) tag(r InteractionId) *Tag {
	title := tagTitle(r.Path().String())
	name := tagName(title)

	t, ok := c.Tags.Get(name)
	if !ok {
		t = newEmptyTag(r)
		c.Tags.Set(name, t)
	}

	return t
}

func (c *Catalog) AddJSight(version string) error {
	if c.JSightVersion != "" {
		return errors.New("directive JSIGHT gotta be only one time")
	}

	c.JSightVersion = version

	return nil
}

func (c *Catalog) AddInfo(d directive.Directive) error {
	if c.Info != nil {
		return errors.New("directive INFO gotta be only one time")
	}
	c.Info = &Info{Directive: d}
	return nil
}

func (c *Catalog) AddTitle(name string) error {
	if c.Info.Title != "" {
		return errors.New(jerr.NotUniqueDirective)
	}
	c.Info.Title = name
	return nil
}

func (c *Catalog) AddVersion(version string) error {
	if c.Info.Version != "" {
		return errors.New(jerr.NotUniqueDirective)
	}
	c.Info.Version = version
	return nil
}

func (c *Catalog) AddDescriptionToInfo(text string) error {
	if c.Info.Description != nil {
		return errors.New(jerr.NotUniqueDirective)
	}

	c.Info.Description = &text

	return nil
}

func (c *Catalog) AddMethod(d directive.Directive) error {
	rk, err := newHttpInteractionId(d)
	if err != nil {
		return err
	}

	if c.HttpInteractions.Has(rk) {
		return fmt.Errorf("method is already defined in resource %s", rk.String())
	}

	t := c.tag(rk)
	t.appendInteractionId(rk)

	rm := initHttpInteraction(rk.path, rk.method, d.Annotation, t.Name)
	c.HttpInteractions.Set(rk, &rm)

	return nil
}

func (c *Catalog) AddDescriptionToHttpMethod(d directive.Directive, text string) error {
	rk, err := newHttpInteractionId(d)
	if err != nil {
		return err
	}

	if !c.HttpInteractions.Has(rk) {
		return fmt.Errorf("%s %q", jerr.ResourceNotFound, rk.String())
	}

	v := c.HttpInteractions.GetValue(rk)

	if v.Description != nil {
		return errors.New(jerr.NotUniqueDirective)
	}

	c.HttpInteractions.Update(rk, func(v *HttpInteraction) *HttpInteraction {
		v.Description = &text
		return v
	})

	return nil
}

func (c *Catalog) AddQueryToCurrentMethod(d directive.Directive, q Query) error {
	rk, err := newHttpInteractionId(d)
	if err != nil {
		return err
	}

	if !c.HttpInteractions.Has(rk) {
		return fmt.Errorf("%s %q", jerr.ResourceNotFound, rk.String())
	}

	v := c.HttpInteractions.GetValue(rk)

	if v.Query != nil {
		return errors.New(jerr.NotUniqueDirective)
	}

	c.HttpInteractions.Update(rk, func(v *HttpInteraction) *HttpInteraction {
		v.Query = &q
		return v
	})

	return nil
}

func (c *Catalog) AddResponse(code string, annotation string, d directive.Directive) error {
	rk, err := newHttpInteractionId(d)
	if err != nil {
		return err
	}

	r := HTTPResponse{Code: code, Annotation: annotation, Directive: d}

	c.HttpInteractions.Update(rk, func(v *HttpInteraction) *HttpInteraction {
		v.Responses = append(v.Responses, r)
		return v
	})

	return nil
}

func (c *Catalog) AddResponseBody(
	schemaName string,
	schemaBytes bytes.Bytes,
	bodyFormat SerializeFormat,
	sn notation.SchemaNotation,
	d directive.Directive,
	tt *UserSchemas,
	rr map[string]jschemaLib.Rule,
) *jerr.JApiError {
	rk, err := newHttpInteractionId(d)
	if err != nil {
		return d.KeywordError(err.Error())
	}

	if !c.HttpInteractions.Has(rk) {
		return d.KeywordError(fmt.Sprintf("%s %q", jerr.ResourceNotFound, rk.String()))
	}

	v := c.HttpInteractions.GetValue(rk)

	i := len(v.Responses) - 1
	if i == -1 {
		return d.KeywordError(fmt.Sprintf("%s for %q", jerr.ResponsesIsEmpty, rk.String()))
	}

	httpResponseBody, je := NewHTTPResponseBody(schemaName, schemaBytes, bodyFormat, sn, d, tt, rr)
	if je != nil {
		return je
	}

	c.HttpInteractions.Update(rk, func(v *HttpInteraction) *HttpInteraction {
		v.Responses[i].Body = &httpResponseBody
		return v
	})

	return nil
}

func (c *Catalog) AddResponseHeaders(s Schema, d directive.Directive) error {
	rk, err := newHttpInteractionId(d)
	if err != nil {
		return err
	}

	if !c.HttpInteractions.Has(rk) {
		return fmt.Errorf("%s %q", jerr.ResourceNotFound, rk.String())
	}

	v := c.HttpInteractions.GetValue(rk)

	i := len(v.Responses) - 1
	if i == -1 {
		return fmt.Errorf("%s for %q", jerr.ResponsesIsEmpty, rk.String())
	}

	if v.Responses[i].Headers != nil {
		return errors.New(jerr.NotUniqueDirective)
	}

	c.HttpInteractions.Update(rk, func(v *HttpInteraction) *HttpInteraction {
		v.Responses[i].Headers = &HTTPResponseHeaders{Schema: &s, Directive: d}
		return v
	})

	return nil
}

func (c *Catalog) AddServer(name string, annotation string) error {
	if c.Servers.Has(name) {
		return fmt.Errorf("duplicate server name %q", name)
	}

	server := new(Server)
	server.Annotation = annotation

	c.Servers.Set(name, server)

	return nil
}

func (c *Catalog) AddBaseUrl(serverName string, path string) error {
	if !c.Servers.Has(serverName) {
		return fmt.Errorf("server not found for %s", serverName)
	}

	v := c.Servers.GetValue(serverName)

	if v.BaseUrl != "" {
		return errors.New("BaseUrl already defined")
	}

	v.BaseUrl = path

	// if d.BodyCoords.IsSet() {
	// 	// baseurl has jschema body
	// 	s, err := UnmarshalSchema("", d.BodyCoords.Read())
	// 	if err != nil {
	// 		if e, ok := err.(kit.Error); ok {
	// 			return c.japiError(e.Message(), d.BodyCoords.B()+bytes.Index(e.Position()))
	// 		}
	// 		return c.japiError(err.Error(), d.BodyBegin())
	// 	}
	//
	// 	if s.ContentJSight.TokenType != objectStr && s.ContentJSight.TokenType != shortcutStr {
	// 		return c.japiError("the body of the BaseUrl directive can contain an object schema", d.BodyBegin())
	// 	}
	//
	// 	c.Servers.Update(serverName, func(v *Server) *Server {
	// 		v.BaseUrlVariables = &baseUrlVariables{
	// 			Schema:    &s,
	// 			directive: d,
	// 		}
	// 		return v
	// 	})
	// }

	return nil
}

func (c *Catalog) AddType(
	d directive.Directive,
	tt *UserSchemas,
	rr map[string]jschemaLib.Rule,
) *jerr.JApiError {
	name := d.Parameter("Name")

	if c.UserTypes.Has(name) {
		return d.KeywordError(fmt.Sprintf("duplicate type name %q", name))
	}

	userType := &UserType{
		Annotation: d.Annotation,
		Directive:  d,
	}
	typeNotation, err := notation.NewSchemaNotation(d.Parameter("SchemaNotation"))
	if err != nil {
		return d.KeywordError(err.Error())
	}

	switch typeNotation {
	case notation.SchemaNotationJSight:
		if !d.BodyCoords.IsSet() {
			return d.KeywordError(jerr.EmptyBody)
		}
		b := d.BodyCoords.Read()
		schema, err := UnmarshalSchema(name, b, tt, rr)
		if err != nil {
			var e kit.Error
			if errors.As(err, &e) {
				return d.BodyErrorIndex(e.Message(), e.Position())
			}
			return d.KeywordError(err.Error())
		}
		userType.Schema = schema
	case notation.SchemaNotationRegex:
		if !d.BodyCoords.IsSet() {
			return d.KeywordError(jerr.EmptyBody)
		}
		userType.Schema = NewRegexSchema(d.BodyCoords.Read())
	case notation.SchemaNotationAny, notation.SchemaNotationEmpty:
		userType.Schema = NewSchema(typeNotation)
	}

	c.UserTypes.Set(name, userType)

	return nil
}

func (c *Catalog) AddRequest(d directive.Directive) error {
	rk, err := newHttpInteractionId(d)
	if err != nil {
		return err
	}

	c.HttpInteractions.Update(rk, func(v *HttpInteraction) *HttpInteraction {
		if v.Request == nil {
			v.Request = &HTTPRequest{
				Directive: d,
			}
		}
		return v
	})

	return nil
}

func (c *Catalog) AddRequestBody(s Schema, f SerializeFormat, d directive.Directive) error {
	rk, err := newHttpInteractionId(d)
	if err != nil {
		return err
	}

	if !c.HttpInteractions.Has(rk) {
		return fmt.Errorf("%s %q", jerr.ResourceNotFound, rk.String())
	}

	v := c.HttpInteractions.GetValue(rk)

	if v.Request == nil {
		return fmt.Errorf("%s for %q", jerr.RequestIsEmpty, rk.String())
	}

	if v.Request.HTTPRequestBody != nil {
		return errors.New(jerr.NotUniqueDirective)
	}

	c.HttpInteractions.Update(rk, func(v *HttpInteraction) *HttpInteraction {
		v.Request.HTTPRequestBody = &HTTPRequestBody{Format: f, Schema: &s, Directive: d}
		return v
	})

	return nil
}

func (c *Catalog) AddRequestHeaders(s Schema, d directive.Directive) error {
	rk, err := newHttpInteractionId(d)
	if err != nil {
		return err
	}

	if !c.HttpInteractions.Has(rk) {
		return fmt.Errorf("%s %q", jerr.ResourceNotFound, rk.String())
	}

	v := c.HttpInteractions.GetValue(rk)

	if v.Request == nil {
		return fmt.Errorf("%s for %q", jerr.RequestIsEmpty, rk.String())
	}

	if v.Request.HTTPRequestHeaders != nil {
		return errors.New(jerr.NotUniqueDirective)
	}

	c.HttpInteractions.Update(rk, func(v *HttpInteraction) *HttpInteraction {
		v.Request.HTTPRequestHeaders = &HTTPRequestHeaders{Schema: &s, Directive: d}
		return v
	})

	return nil
}

func (_ *Catalog) AddProtocol(_ directive.Directive) error {
	// TODO
	return nil
}

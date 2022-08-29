package catalog

import (
	"errors"
	"fmt"

	jschemaLib "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/bytes"
	"github.com/jsightapi/jsight-schema-go-library/kit"
	"github.com/jsightapi/jsight-schema-go-library/rules/enum"

	"github.com/jsightapi/jsight-api-go-library/directive"
	"github.com/jsightapi/jsight-api-go-library/jerr"
	"github.com/jsightapi/jsight-api-go-library/notation"
)

func (c *Catalog) AddTag(name, title string) error {
	if c.Tags.Has(TagName(name)) {
		return fmt.Errorf("%s (%q)", jerr.DuplicateNames, name)
	}

	if title == "" {
		title = name
	}

	t := NewTag(name, title)

	c.Tags.Set(t.Name, t)

	return nil
}

func (c *Catalog) tagNames(d directive.Directive, id InteractionID) ([]TagName, *jerr.JApiError) {
	tt, je := c.tags(d, id)
	if je != nil {
		return nil, je
	}

	tns := make([]TagName, 0, len(tt))

	for _, t := range tt {
		t.appendInteractionID(id)
		tns = append(tns, t.Name)
	}

	return tns, nil
}

// tags return tag names for directives HTTP request method or JSON-RPC method.
func (c *Catalog) tags(d directive.Directive, id InteractionID) ([]*Tag, *jerr.JApiError) {
	if td := getChildrenTagsDirective(d); td != nil { // child directive Tags for HTTP or JSON-RPC methods
		return c.tagsFromTagsDirective(td)
	}

	if td := getParentTagsDirective(d); td != nil { // parent URL
		return c.tagsFromTagsDirective(td)
	}

	return []*Tag{c.pathTag(id)}, nil
}

// pathTag returns a Tag from the collection, or creates a new one and adds it to the collection
func (c *Catalog) pathTag(r InteractionID) *Tag {
	t := newPathTag(r)

	if tt, ok := c.Tags.Get(t.Name); ok {
		return tt
	}

	c.Tags.Set(t.Name, t)

	return t
}

func getChildrenTagsDirective(d directive.Directive) *directive.Directive {
	for _, dd := range d.Children {
		if dd.Type() == directive.Tags {
			return dd
		}
	}
	return nil
}

func getParentTagsDirective(d directive.Directive) *directive.Directive {
	if d.Parent != nil && d.Parent.Type() == directive.URL {
		return getChildrenTagsDirective(*d.Parent)
	}
	return nil
}

func (c *Catalog) tagsFromTagsDirective(d *directive.Directive) ([]*Tag, *jerr.JApiError) {
	if je := checkTagsDirective(d); je != nil {
		return nil, je
	}

	tt := make([]*Tag, 0, d.UnnamedParametersLen())

	for _, name := range d.UnnamedParameter() {
		tn := TagName(name)

		t, ok := c.Tags.Get(tn)
		if !ok {
			return nil, d.KeywordError(fmt.Sprintf("%s %q", jerr.TagNotFound, tn))
		}

		tt = append(tt, t)
	}

	return tt, nil
}

func checkTagsDirective(d *directive.Directive) *jerr.JApiError {
	if d.Annotation != "" {
		return d.KeywordError(jerr.AnnotationIsForbiddenForTheDirective)
	}

	if !d.HasUnnamedParameter() {
		return d.KeywordError(jerr.RequiredParameterNotSpecified)
	}

	return nil
}

func (c *Catalog) AddDescriptionToTag(name, description string) error {
	n := TagName(name)

	t, ok := c.Tags.Get(n)
	if !ok {
		return fmt.Errorf("%s (%q)", jerr.TagNotFound, n)
	}

	if t.Description != nil {
		return errors.New(jerr.NotUniqueDirective)
	}

	c.Tags.Update(n, func(v *Tag) *Tag {
		v.Description = &description
		return v
	})

	return nil
}

func (c *Catalog) AddTagDescription(name, title string) error {
	t := NewTag(name, title)

	if _, ok := c.Tags.Get(t.Name); ok {
		return errors.New("")
	}

	c.Tags.Set(t.Name, t)

	return nil
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

func (c *Catalog) AddHTTPMethod(d directive.Directive) *jerr.JApiError {
	httpID, err := newHTTPInteractionID(d)
	if err != nil {
		return d.KeywordError(err.Error())
	}

	if c.Interactions.Has(httpID) {
		return d.KeywordError(fmt.Sprintf("method is already defined in resource %s", httpID.String()))
	}

	in := newHTTPInteraction(httpID, d.Annotation)

	tns, je := c.tagNames(d, httpID)
	if je != nil {
		return je
	}
	for _, tn := range tns {
		in.appendTagName(tn)
	}

	c.Interactions.Set(httpID, in)

	return nil
}

func (c *Catalog) AddDescriptionToHTTPMethod(d directive.Directive, text string) error {
	httpID, err := newHTTPInteractionID(d)
	if err != nil {
		return err
	}

	if !c.Interactions.Has(httpID) {
		return fmt.Errorf("%s %q", jerr.HTTPResourceNotFound, httpID.String())
	}

	v := c.Interactions.GetValue(httpID).(*HTTPInteraction) //nolint:errcheck

	if v.Description != nil {
		return errors.New(jerr.NotUniqueDirective)
	}

	c.Interactions.Update(httpID, func(v Interaction) Interaction {
		v.(*HTTPInteraction).Description = &text
		return v
	})

	return nil
}

func (c *Catalog) AddDescriptionToJsonRpcMethod(d directive.Directive, text string) error {
	rpcId, err := newJsonRpcInteractionId(d)
	if err != nil {
		return err
	}

	if !c.Interactions.Has(rpcId) {
		return fmt.Errorf("%s %q", jerr.HTTPResourceNotFound, rpcId.String())
	}

	v := c.Interactions.GetValue(rpcId).(*JsonRpcInteraction) //nolint:errcheck

	if v.Description != nil {
		return errors.New(jerr.NotUniqueDirective)
	}

	c.Interactions.Update(rpcId, func(v Interaction) Interaction {
		v.(*JsonRpcInteraction).Description = &text
		return v
	})

	return nil
}

func (c *Catalog) AddQueryToCurrentMethod(d directive.Directive, q Query) error {
	httpID, err := newHTTPInteractionID(d)
	if err != nil {
		return err
	}

	if !c.Interactions.Has(httpID) {
		return fmt.Errorf("%s %q", jerr.HTTPResourceNotFound, httpID.String())
	}

	v := c.Interactions.GetValue(httpID).(*HTTPInteraction) //nolint:errcheck

	if v.Query != nil {
		return errors.New(jerr.NotUniqueDirective)
	}

	c.Interactions.Update(httpID, func(v Interaction) Interaction {
		v.(*HTTPInteraction).Query = &q
		return v
	})

	return nil
}

func (c *Catalog) AddResponse(code string, annotation string, d directive.Directive) error {
	httpID, err := newHTTPInteractionID(d)
	if err != nil {
		return err
	}

	r := HTTPResponse{Code: code, Annotation: annotation, Directive: d}

	c.Interactions.Update(httpID, func(v Interaction) Interaction {
		v.(*HTTPInteraction).Responses = append(v.(*HTTPInteraction).Responses, r)
		return v
	})

	return nil
}

func (c *Catalog) AddResponseBody(
	schemaBytes bytes.Bytes,
	bodyFormat SerializeFormat,
	sn notation.SchemaNotation,
	d directive.Directive,
	tt *UserSchemas,
	rr map[string]jschemaLib.Rule,
) *jerr.JApiError {
	httpID, err := newHTTPInteractionID(d)
	if err != nil {
		return d.KeywordError(err.Error())
	}

	if !c.Interactions.Has(httpID) {
		return d.KeywordError(fmt.Sprintf("%s %q", jerr.HTTPResourceNotFound, httpID.String()))
	}

	v := c.Interactions.GetValue(httpID).(*HTTPInteraction) //nolint:errcheck

	i := len(v.Responses) - 1
	if i == -1 {
		return d.KeywordError(fmt.Sprintf("%s for %q", jerr.ResponsesIsEmpty, httpID.String()))
	}

	httpResponseBody, je := NewHTTPResponseBody(schemaBytes, bodyFormat, sn, d, tt, rr)
	if je != nil {
		return je
	}

	c.Interactions.Update(httpID, func(v Interaction) Interaction {
		v.(*HTTPInteraction).Responses[i].Body = &httpResponseBody
		return v
	})

	return nil
}

func (c *Catalog) AddResponseHeaders(s Schema, d directive.Directive) error {
	httpID, err := newHTTPInteractionID(d)
	if err != nil {
		return err
	}

	if !c.Interactions.Has(httpID) {
		return fmt.Errorf("%s %q", jerr.HTTPResourceNotFound, httpID.String())
	}

	v := c.Interactions.GetValue(httpID).(*HTTPInteraction) //nolint:errcheck

	i := len(v.Responses) - 1
	if i == -1 {
		return fmt.Errorf("%s for %q", jerr.ResponsesIsEmpty, httpID.String())
	}

	if v.Responses[i].Headers != nil {
		return errors.New(jerr.NotUniqueDirective)
	}

	c.Interactions.Update(httpID, func(v Interaction) Interaction {
		v.(*HTTPInteraction).Responses[i].Headers = &HTTPResponseHeaders{Schema: &s, Directive: d}
		return v
	})

	return nil
}

func (c *Catalog) AddServer(name string, annotation string) error {
	if c.Servers.Has(name) {
		return fmt.Errorf("%s (%q)", jerr.DuplicateNames, name)
	}

	server := new(Server)
	server.Annotation = annotation

	c.Servers.Set(name, server)

	return nil
}

func (c *Catalog) AddBaseURL(serverName string, path string) error {
	if !c.Servers.Has(serverName) {
		return fmt.Errorf("server not found for %s", serverName)
	}

	v := c.Servers.GetValue(serverName)

	if v.BaseUrl != "" {
		return errors.New("BaseURL already defined")
	}

	v.BaseUrl = path

	// if d.BodyCoords.IsSet() {
	// 	// baseurl has jschema body
	// 	s, err := UnmarshalJSightSchema("", d.BodyCoords.Read())
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
	// 		v.BaseUrlVariables = &baseURLVariables{
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
	name := d.NamedParameter("Name")

	if c.UserTypes.Has(name) {
		return d.KeywordError(fmt.Sprintf("%s (%q)", jerr.DuplicateNames, name))
	}

	userType := &UserType{
		Annotation: d.Annotation,
		Directive:  d,
	}
	typeNotation, err := notation.NewSchemaNotation(d.NamedParameter("SchemaNotation"))
	if err != nil {
		return d.KeywordError(err.Error())
	}

	switch typeNotation {
	case notation.SchemaNotationJSight:
		if !d.BodyCoords.IsSet() {
			return d.KeywordError(jerr.EmptyBody)
		}
		schema, err := UnmarshalJSightSchema(name, d.BodyCoords.Read(), tt, rr)
		if err != nil {
			return adoptErrorForAddType(d, err)
		}
		userType.Schema = schema
	case notation.SchemaNotationRegex:
		if !d.BodyCoords.IsSet() {
			return d.KeywordError(jerr.EmptyBody)
		}
		schema, err := UnmarshalRegexSchema(name, d.BodyCoords.Read())
		if err != nil {
			return adoptErrorForAddType(d, err)
		}
		userType.Schema = schema
	case notation.SchemaNotationAny, notation.SchemaNotationEmpty:
		userType.Schema = NewSchema(typeNotation)
	}

	c.UserTypes.Set(name, userType)

	return nil
}

func (c *Catalog) AddRequest(d directive.Directive) error {
	httpID, err := newHTTPInteractionID(d)
	if err != nil {
		return err
	}

	c.Interactions.Update(httpID, func(v Interaction) Interaction {
		if v.(*HTTPInteraction).Request == nil {
			v.(*HTTPInteraction).Request = &HTTPRequest{
				Directive: d,
			}
		}
		return v
	})

	return nil
}

func (c *Catalog) AddRequestBody(s Schema, f SerializeFormat, d directive.Directive) error {
	httpID, err := newHTTPInteractionID(d)
	if err != nil {
		return err
	}

	if !c.Interactions.Has(httpID) {
		return fmt.Errorf("%s %q", jerr.HTTPResourceNotFound, httpID.String())
	}

	v := c.Interactions.GetValue(httpID).(*HTTPInteraction) //nolint:errcheck

	if v.Request == nil {
		return fmt.Errorf("%s for %q", jerr.RequestIsEmpty, httpID.String())
	}

	if v.Request.HTTPRequestBody != nil {
		return errors.New(jerr.NotUniqueDirective)
	}

	c.Interactions.Update(httpID, func(v Interaction) Interaction {
		v.(*HTTPInteraction).Request.HTTPRequestBody = &HTTPRequestBody{Format: f, Schema: &s, Directive: d}
		return v
	})

	return nil
}

func (c *Catalog) AddRequestHeaders(s Schema, d directive.Directive) error {
	httpID, err := newHTTPInteractionID(d)
	if err != nil {
		return err
	}

	if !c.Interactions.Has(httpID) {
		return fmt.Errorf("%s %q", jerr.HTTPResourceNotFound, httpID.String())
	}

	v := c.Interactions.GetValue(httpID).(*HTTPInteraction) //nolint:errcheck

	if v.Request == nil {
		return fmt.Errorf("%s for %q", jerr.RequestIsEmpty, httpID.String())
	}

	if v.Request.HTTPRequestHeaders != nil {
		return errors.New(jerr.NotUniqueDirective)
	}

	c.Interactions.Update(httpID, func(v Interaction) Interaction {
		v.(*HTTPInteraction).Request.HTTPRequestHeaders = &HTTPRequestHeaders{Schema: &s, Directive: d}
		return v
	})

	return nil
}

func (c *Catalog) AddJsonRpcMethod(d directive.Directive) *jerr.JApiError {
	rpcId, err := newJsonRpcInteractionId(d)
	if err != nil {
		return d.KeywordError(err.Error())
	}

	if c.Interactions.Has(rpcId) {
		return d.KeywordError(fmt.Sprintf("method is already defined in resource %s", rpcId.String()))
	}

	in := newJsonRpcInteraction(rpcId, d.NamedParameter("MethodName"), d.Annotation)

	tns, je := c.tagNames(d, rpcId)
	if je != nil {
		return je
	}
	for _, tn := range tns {
		in.appendTagName(tn)
	}

	c.Interactions.Set(rpcId, in)

	return nil
}

func (c *Catalog) AddJsonRpcParams(s Schema, d directive.Directive) error {
	rpcId, err := newJsonRpcInteractionId(d)
	if err != nil {
		return err
	}

	if !c.Interactions.Has(rpcId) {
		return fmt.Errorf("%s %q", jerr.JsonRpcResourceNotFound, rpcId.String())
	}

	v := c.Interactions.GetValue(rpcId).(*JsonRpcInteraction) //nolint:errcheck

	if v.Params != nil {
		return errors.New(jerr.NotUniqueDirective)
	}

	c.Interactions.Update(rpcId, func(v Interaction) Interaction {
		v.(*JsonRpcInteraction).Params = &jsonRpcParams{Schema: &s, Directive: d}
		return v
	})

	return nil
}

func (c *Catalog) AddJsonRpcResult(s Schema, d directive.Directive) error {
	rpcId, err := newJsonRpcInteractionId(d)
	if err != nil {
		return err
	}

	if !c.Interactions.Has(rpcId) {
		return fmt.Errorf("%s %q", jerr.JsonRpcResourceNotFound, rpcId.String())
	}

	v := c.Interactions.GetValue(rpcId).(*JsonRpcInteraction) //nolint:errcheck

	if v.Result != nil {
		return errors.New(jerr.NotUniqueDirective)
	}

	c.Interactions.Update(rpcId, func(v Interaction) Interaction {
		v.(*JsonRpcInteraction).Result = &jsonRpcResult{Schema: &s, Directive: d}
		return v
	})

	return nil
}

func (c *Catalog) AddEnum(d *directive.Directive, e *enum.Enum) *jerr.JApiError {
	name := d.NamedParameter("Name")
	if c.UserEnums.Has(name) {
		return d.KeywordError(fmt.Sprintf("%s (%q)", jerr.DuplicateNames, name))
	}

	r, err := c.enumDirectiveToUserRule(d, e)
	if err != nil {
		return d.KeywordError(err.Error())
	}

	c.UserEnums.Set(name, r)
	return nil
}

func (*Catalog) enumDirectiveToUserRule(d *directive.Directive, e *enum.Enum) (*UserRule, error) {
	vv, err := e.Values()
	if err != nil {
		return nil, err
	}

	r := Rule{
		TokenType: RuleTokenTypeArray,
	}

	for _, v := range vv {
		r.Children = append(r.Children, Rule{
			TokenType:   RuleTokenType(v.Type.ToTokenType()),
			ScalarValue: v.Value.Unquote().String(),
			Note:        v.Comment,
		})
	}

	return &UserRule{
		Annotation: d.Annotation,
		Value:      r,
		Directive:  d,
	}, nil
}

func adoptErrorForAddType(d directive.Directive, err error) *jerr.JApiError {
	var e kit.Error
	if errors.As(err, &e) {
		return d.BodyErrorIndex(e.Message(), e.Position())
	}
	return d.KeywordError(err.Error())
}

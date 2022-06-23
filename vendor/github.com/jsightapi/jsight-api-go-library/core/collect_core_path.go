package core

import (
	"errors"

	"github.com/jsightapi/jsight-schema-go-library/kit"

	"github.com/jsightapi/jsight-api-go-library/catalog"
	"github.com/jsightapi/jsight-api-go-library/directive"
	"github.com/jsightapi/jsight-api-go-library/jerr"
)

func (core *JApiCore) collectPaths(dd []*directive.Directive) *jerr.JApiError {
	for i := 0; i != len(dd); i++ {
		switch dd[i].Type() {
		case directive.Macro:
			continue
		case directive.Path:
			if je := core.collectPathVariables(dd[i]); je != nil {
				return je
			}
		default:
			// does nothing
		}

		if dd[i].Children != nil {
			if je := core.collectPaths(dd[i].Children); je != nil {
				return je
			}
		}
	}
	return nil
}

func (core *JApiCore) collectPathVariables(d *directive.Directive) *jerr.JApiError {
	if d.Annotation != "" {
		return d.KeywordError(jerr.AnnotationIsForbiddenForTheDirective)
	}
	if !d.BodyCoords.IsSet() {
		return d.KeywordError("there is no body for the Path directive")
	}

	s, err := catalog.UnmarshalSchema("", d.BodyCoords.Read(), core.userTypes, core.rules)
	if err != nil {
		var e kit.Error
		if errors.As(err, &e) {
			return d.BodyErrorIndex(e.Message(), e.Position())
		}
		return d.KeywordError(err.Error())
	}

	path, err := d.Path()
	if err != nil {
		return d.KeywordError(err.Error())
	}

	pp, err := PathParameters(path)
	if err != nil {
		return d.KeywordError(err.Error())
	}

	if d.Parent == nil {
		return d.KeywordError("parent directive not found")
	}

	parentDirective := *d.Parent

	if len(core.rawPathVariables) != 0 {
		prevParent := core.rawPathVariables[len(core.rawPathVariables)-1].parentDirective
		if prevParent.Equal(parentDirective) {
			return d.KeywordError(jerr.NotUniqueDirective)
		}
	}

	core.rawPathVariables = append(core.rawPathVariables, rawPathVariable{
		pathDirective:   *d,
		parentDirective: parentDirective,
		schema:          s,
		parameters:      pp,
	})

	return nil
}

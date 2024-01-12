package core

import (
	"github.com/jsightapi/jsight-api-core/directive"
	"github.com/jsightapi/jsight-api-core/jerr"
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
		return d.KeywordError(jerr.BodyIsEmpty)
	}

	s, err := newPathVariablesSchema(d.BodyCoords.Read(), core.UserTypesData())
	if err != nil {
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
		return d.KeywordError(jerr.ParentNotFound)
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
		schema:          s.JSchema,
		parameters:      pp,
	})

	return nil
}

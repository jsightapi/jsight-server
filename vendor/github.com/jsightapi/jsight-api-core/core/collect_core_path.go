package core

import (
	"github.com/jsightapi/jsight-api-core/directive"
	"github.com/jsightapi/jsight-api-core/jerr"
)

// runs through all directives which has path as a param (URL, Get, etc.)
// if this directive was not processed during processing of its child Path directive (did not have that child),
// then it may have unspecified path params, that must be added to core.rawPathVariables for future processing
// this is a temp workaround, until better architecture for path params is designed
func (core *JApiCore) addMissedUndefindedPathVariables(dd []*directive.Directive) *jerr.JApiError {
	for i := 0; i != len(dd); i++ {
		d := dd[i]
		switch d.Type() {
		case directive.URL, directive.Get, directive.Post, directive.Put, directive.Delete, directive.Patch:
			var used bool
			for i := 0; i < len(core.rawPathVariables); i++ {
				if core.rawPathVariables[i].pathDirective.Parent == d {
					used = true
				}
			}
			if !used { // register this params as rawPathVariables
				path, err := d.Path()
				if err != nil {
					return d.KeywordError(err.Error())
				}

				pp, err := PathParameters(path)
				if err != nil {
					return d.KeywordError(err.Error())
				}

				core.rawPathVariables = append(core.rawPathVariables, rawPathVariable{
					pathDirective: *d,
					parameters:    pp,
					imitated:      true,
				})
			}

		default:
			// does nothing
		}
	}
	return nil
}

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

// NOTE: works specifically with Path directive.
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

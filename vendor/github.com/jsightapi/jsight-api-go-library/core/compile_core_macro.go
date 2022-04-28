package core

import (
	"fmt"

	"github.com/jsightapi/jsight-api-go-library/directive"
	"github.com/jsightapi/jsight-api-go-library/jerr"
)

func (core *JApiCore) collectMacro() *jerr.JAPIError {
	for i := 0; i != len(core.directives); i++ {
		if core.directives[i].Type() == directive.Macro {
			if je := core.addMacro(core.directives[i]); je != nil {
				return je
			}

			core.directives = removeDirectiveFromSlice(core.directives, i)
			i--
		}
	}
	return nil
}

func (core *JApiCore) addMacro(d *directive.Directive) *jerr.JAPIError {
	if d.Annotation != "" {
		return d.KeywordError(jerr.AnnotationIsForbiddenForTheDirective)
	}

	name := d.Parameter("Name")

	if name == "" {
		return d.KeywordError(fmt.Sprintf("%s (%s)", jerr.RequiredParameterNotSpecified, "Name"))
	}
	if d.Children == nil {
		return d.KeywordError("empty macro")
	}

	if _, ok := core.macro[name]; ok {
		return d.KeywordError(fmt.Sprintf("duplicate macro name %q", name))
	}

	core.macro[name] = d

	return nil
}

package core

import (
	"fmt"

	"github.com/jsightapi/jsight-api-core/directive"
	"github.com/jsightapi/jsight-api-core/jerr"
)

func (core *JApiCore) collectMacro() *jerr.JApiError {
	for i := 0; i != len(core.directives); i++ {
		if core.directives[i].Type() == directive.Macro {
			if je := core.addMacro(core.directives[i]); je != nil {
				return je
			}

			core.directives = append(core.directives[:i], core.directives[i+1:]...)
			i--
		}
	}
	return nil
}

func (core *JApiCore) addMacro(d *directive.Directive) *jerr.JApiError {
	if d.Annotation != "" {
		return d.KeywordError(jerr.AnnotationIsForbiddenForTheDirective)
	}

	name := d.NamedParameter("Name")

	if name == "" {
		return d.KeywordError(fmt.Sprintf("%s (%s)", jerr.RequiredParameterNotSpecified, "Name"))
	}
	if d.Children == nil {
		return d.KeywordError(jerr.MacroIsEmpty)
	}

	if _, ok := core.macro[name]; ok {
		return d.KeywordError(fmt.Sprintf("%s %q", jerr.DuplicateNames, name))
	}

	core.macro[name] = d

	return nil
}

func (core *JApiCore) checkMacroForRecursion() *jerr.JApiError {
	for macroName, macro := range core.macro {
		if je := findPaste(macroName, macro); je != nil {
			return je
		}
	}
	return nil
}

func findPaste(macroName string, d *directive.Directive) *jerr.JApiError {
	if d.Type() == directive.Paste {
		switch d.NamedParameter("Name") {
		case "":
			return d.KeywordError(fmt.Sprintf("%s (%s)", jerr.RequiredParameterNotSpecified, "Name"))

		case macroName:
			return d.KeywordError(jerr.RecursionIsProhibited)
		}
	} else if d.Children != nil {
		for _, c := range d.Children {
			if je := findPaste(macroName, c); je != nil {
				return je
			}
		}
	}
	return nil
}

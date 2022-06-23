package core

import (
	"fmt"

	"github.com/jsightapi/jsight-api-go-library/directive"
	"github.com/jsightapi/jsight-api-go-library/jerr"
)

func (core *JApiCore) compileCore() *jerr.JApiError {
	if je := core.collectMacro(); je != nil {
		return je
	}
	if je := core.checkMacroForRecursion(); je != nil {
		return je
	}
	if je := core.processPaste(); je != nil {
		return je
	}

	if je := core.collectRules(); je != nil {
		return je
	}

	core.collectUserTypes()

	if je := core.compileUserTypes(); je != nil {
		return je
	}

	return core.collectPaths(core.directivesWithPastes)
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
		switch d.Parameter("Name") {
		case "":
			return d.KeywordError(fmt.Sprintf("%s (%s)", jerr.RequiredParameterNotSpecified, "Name"))

		case macroName:
			return d.KeywordError("recursion is prohibited")
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

func (core *JApiCore) collectUserTypes() {
	for _, d := range core.directivesWithPastes {
		if d.Type() == directive.Type {
			core.catalog.AddRawUserType(d)
		}
	}
}

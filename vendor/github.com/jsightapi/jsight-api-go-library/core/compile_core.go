package core

import (
	"fmt"

	"github.com/jsightapi/jsight-api-go-library/directive"
	"github.com/jsightapi/jsight-api-go-library/jerr"
)

func (core *JApiCore) compileCore() *jerr.JAPIError {
	if je := core.collectMacro(); je != nil {
		return je
	}
	if je := core.checkMacroForRecursion(); je != nil {
		return je
	}
	if je := core.processPaste(); je != nil {
		return je
	}

	core.collectUserTypes()

	if err := core.compileUserTypes(); err != nil {
		return err
	}

	return core.collectPaths(core.directivesWithPastes)
}

func (core *JApiCore) checkMacroForRecursion() *jerr.JAPIError {
	for macroName, macro := range core.macro {
		if je := findPaste(macroName, macro); je != nil {
			return je
		}
	}
	return nil
}

func findPaste(macroName string, d *directive.Directive) *jerr.JAPIError {
	if d.Type() == directive.Paste {
		switch d.Parameter("Name") {
		case "":
			return d.KeywordError(fmt.Sprintf("%s (%s)", jerr.RequiredParameterNotSpecified, "Name"))

		case macroName:
			return d.KeywordError("recursion is prohibited")
		}
	} else if d.Children != nil {
		for i := 0; i != len(d.Children); i++ {
			if je := findPaste(macroName, d.Children[i]); je != nil {
				return je
			}
		}
	}
	return nil
}

func (core *JApiCore) collectUserTypes() {
	for i := 0; i != len(core.directivesWithPastes); i++ {
		if core.directivesWithPastes[i].Type() == directive.Type {
			core.catalog.AddRawUserType(core.directivesWithPastes[i])
		}
	}
}

package core

import (
	"github.com/jsightapi/jsight-api-core/jerr"
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

	if je := core.collectTags(); je != nil {
		return je
	}

	if je := core.collectUserTypes(); je != nil {
		return je
	}

	return core.collectPaths(core.directivesWithPastes)
}

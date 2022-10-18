package core

import (
	"fmt"

	"github.com/jsightapi/jsight-api-go-library/directive"
	"github.com/jsightapi/jsight-api-go-library/jerr"
)

func (core *JApiCore) processPaste() *jerr.JApiError {
	core.directivesWithPastes = make([]*directive.Directive, 0, 200)
	core.currentContextDirective = nil
	return core.processPasteDirectiveList(core.directives)
}

func (core *JApiCore) processPasteDirectiveList(list []*directive.Directive) *jerr.JApiError {
	for _, d := range list {
		if je := core.processDirective(d); je != nil {
			return je
		}
	}
	return nil
}

func (core *JApiCore) processDirective(d *directive.Directive) *jerr.JApiError {
	if d.Type() == directive.Paste {
		if je := core.processPasteDirective(d); je != nil {
			return d.KeywordError(je.Error())
		}
		return nil
	}

	dd := d.CopyWoParentAndChildren()
	if je := core.processContext(&dd, &core.directivesWithPastes); je != nil {
		return je
	}

	if d.Children != nil {
		if je := core.processPasteDirectiveList(d.Children); je != nil {
			return je
		}
	}

	if d.HasExplicitContext {
		core.currentContextDirective = dd.Parent
	}

	return nil
}

func (core *JApiCore) processPasteDirective(paste *directive.Directive) *jerr.JApiError {
	if paste.Annotation != "" {
		return paste.KeywordError(jerr.AnnotationIsForbiddenForTheDirective)
	}

	name := paste.NamedParameter("Name")

	if name == "" {
		return paste.KeywordError(fmt.Sprintf("%s (%s)", jerr.RequiredParameterNotSpecified, "Name"))
	}

	macro, ok := core.macro[name]
	if !ok {
		return paste.KeywordError("macro not found")
	}

	if je := core.collectRulesFromDirectives(macro.Children); je != nil {
		return je
	}

	// macro.Children != nil - checked above
	return core.processPasteDirectiveList(macro.Children)
}

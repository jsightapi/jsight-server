package core

import (
	"fmt"

	"github.com/jsightapi/jsight-api-go-library/directive"
	"github.com/jsightapi/jsight-api-go-library/jerr"
)

func (core *JApiCore) collectTags() *jerr.JApiError {
	for _, d := range core.directivesWithPastes {
		if d.Type() == directive.TAG {
			if je := core.collectTag(d); je != nil {
				return je
			}
		}
	}
	return nil
}

func (core *JApiCore) collectTag(d *directive.Directive) *jerr.JApiError {
	if d.NamedParameter("TagName") == "" {
		return d.KeywordError(fmt.Sprintf("%s (%s)", jerr.RequiredParameterNotSpecified, "TagName"))
	}

	if err := core.catalog.AddTag(d.NamedParameter("TagName"), d.Annotation); err != nil {
		return d.KeywordError(err.Error())
	}

	return nil
}

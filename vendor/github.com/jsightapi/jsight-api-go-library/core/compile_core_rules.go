package core

import (
	"github.com/jsightapi/jsight-schema-go-library/rules/enum"

	"github.com/jsightapi/jsight-api-go-library/directive"
	"github.com/jsightapi/jsight-api-go-library/jerr"
)

func (core *JApiCore) collectRules() *jerr.JApiError {
	return core.collectRulesFromDirectives(core.directives)
}

func (core *JApiCore) collectRulesFromDirectives(dd []*directive.Directive) *jerr.JApiError {
	for _, d := range dd {
		if je := core.buildRule(d); je != nil {
			return je
		}
	}
	return nil
}

func (core *JApiCore) buildRule(d *directive.Directive) *jerr.JApiError {
	if d.Type() != directive.Enum {
		return nil
	}

	if !d.BodyCoords.IsSet() {
		return nil
	}

	name := d.NamedParameter("Name")

	r := enum.New(name, d.BodyCoords.Read())
	if err := r.Check(); err != nil {
		return jschemaToJAPIError(err, d)
	}

	core.rules[name] = r
	return core.catalog.AddEnum(d, r)
}

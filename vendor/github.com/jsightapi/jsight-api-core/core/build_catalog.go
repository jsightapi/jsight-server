package core

import (
	stdErrors "errors"
	"fmt"

	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/errs"
	"github.com/jsightapi/jsight-schema-core/kit"

	"github.com/jsightapi/jsight-api-core/directive"
	"github.com/jsightapi/jsight-api-core/jerr"
)

func (core *JApiCore) buildCatalog() *jerr.JApiError {
	if len(core.directivesWithPastes) != 0 && core.directivesWithPastes[0].Type() != directive.Jsight {
		return core.directivesWithPastes[0].KeywordError(jerr.DirectiveJSIGHTShouldBeTheFirst)
	}

	return core.addDirectives()
}

func adoptError(err error) (e *jerr.JApiError) {
	if err == nil {
		return nil
	}

	if stdErrors.As(err, &e) {
		return e
	}

	panic(fmt.Sprintf("Invalid error was given: %#v", err))
}

func safeAddType(curr schema.Schema, n string, ut schema.Schema) error {
	err := curr.AddType(n, ut)
	var e interface{ Code() errs.Code }
	if stdErrors.As(err, &e) && e.Code() == errs.ErrDuplicationOfNameOfTypes {
		err = nil
	}
	return err
}

func (core *JApiCore) checkUserType(name string) *jerr.JApiError {
	err := core.userTypes.GetValue(name).Check()
	if err == nil {
		return nil
	}

	d := core.rawUserTypes.GetValue(name)
	var e kit.Error
	if !stdErrors.As(err, &e) {
		return d.KeywordError(err.Error())
	}

	if e.IncorrectUserType() != "" && e.IncorrectUserType() != name {
		return core.checkUserType(e.IncorrectUserType())
	}

	return d.BodyErrorIndex(e.Message(), e.Index())
}

func jschemaToJAPIError(err error, d *directive.Directive) *jerr.JApiError {
	var e kit.Error
	if stdErrors.As(err, &e) {
		return d.BodyErrorIndex(e.Message(), e.Index())
	}
	return d.KeywordError(err.Error())
}

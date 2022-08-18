package core

import (
	"errors"
	"fmt"

	jschemaLib "github.com/jsightapi/jsight-schema-go-library"
	jerrors "github.com/jsightapi/jsight-schema-go-library/errors"
	"github.com/jsightapi/jsight-schema-go-library/kit"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema"
	"github.com/jsightapi/jsight-schema-go-library/notations/regex"

	"github.com/jsightapi/jsight-api-go-library/directive"
	"github.com/jsightapi/jsight-api-go-library/jerr"
	"github.com/jsightapi/jsight-api-go-library/notation"
)

func (core *JApiCore) buildCatalog() *jerr.JApiError {
	if len(core.directivesWithPastes) != 0 && core.directivesWithPastes[0].Type() != directive.Jsight {
		return core.directivesWithPastes[0].KeywordError("JSIGHT should be the first directive")
	}

	return core.addDirectives()
}

func (core *JApiCore) compileUserTypes() *jerr.JApiError {
	// Two-phase algorithm. On the first step we just create schema for each user
	// type. On the second step we will add all schema to all.
	// This is the simplest solution which allows us to skip building dependency
	// graph between types.

	if err := core.buildUserTypes(); err != nil {
		return err
	}

	err := core.userTypes.Each(func(k string, _ jschemaLib.Schema) error {
		return core.checkUserType(k)
	})
	return adoptError(err)
}

func (core *JApiCore) buildUserTypes() *jerr.JApiError {
	core.Catalog().GetRawUserTypes().EachSafe(func(k string, v *directive.Directive) {
		switch notation.SchemaNotation(v.NamedParameter("SchemaNotation")) {
		case "", notation.SchemaNotationJSight:
			if v.BodyCoords.IsSet() {
				core.userTypes.Set(k, jschema.New(k, v.BodyCoords.Read()))
			}
		case notation.SchemaNotationRegex:
			var oo []regex.Option
			if core.useFixedSeedForRegex {
				oo = append(oo, regex.WithGeneratorSeed(0))
			}
			core.userTypes.Set(k, regex.New(k, v.BodyCoords.Read(), oo...))
		default:
			// nothing
		}
	})

	err := core.userTypes.Each(func(n string, _ jschemaLib.Schema) error {
		return core.compileUserTypeWithAllDependencies(n)
	})
	return adoptError(err)
}

func adoptError(err error) (e *jerr.JApiError) {
	if err == nil {
		return nil
	}

	if errors.As(err, &e) {
		return e
	}

	panic(fmt.Sprintf("Invalid error was given: %#v", err))
}

func (core *JApiCore) compileUserTypeWithAllDependencies(name string) error {
	if _, ok := core.processedUserTypes[name]; ok {
		// This user type already processed, skip.
		return nil
	}
	core.processedUserTypes[name] = struct{}{}

	currUT := core.userTypes.GetValue(name)
	if currUT == nil {
		return nil
	}

	dd := core.catalog.GetRawUserTypes()

	// Add rules before we try to do something with the type.
	for n, r := range core.rules {
		if err := currUT.AddRule(n, r); err != nil {
			return jschemaToJAPIError(err, dd.GetValue(n))
		}
	}

	tt, err := fetchUsedUserTypes(currUT, core.userTypes)
	if err != nil {
		return jschemaToJAPIError(err, dd.GetValue(name))
	}

	for _, n := range tt {
		ut := core.userTypes.GetValue(n)
		if ut == nil {
			continue
		}

		if n != name {
			if err := core.compileUserTypeWithAllDependencies(n); err != nil {
				return err
			}

			if err := core.checkUserTypeDuringBuild(n, ut); err != nil {
				return jschemaToJAPIError(err, dd.GetValue(n))
			}
		}

		if err := safeAddType(currUT, n, ut); err != nil {
			return jschemaToJAPIError(err, dd.GetValue(n))
		}
	}

	// Check user type is correct.
	// We should do it here 'cause it will simplify further processing.
	if err := currUT.Check(); err != nil {
		return jschemaToJAPIError(err, dd.GetValue(name))
	}

	core.userTypes.Set(name, currUT)

	return nil
}

func (core *JApiCore) checkUserTypeDuringBuild(name string, ut jschemaLib.Schema) error {
	// In order to prevent errors in type recursion.
	if _, ok := core.processedUserTypes[name]; ok {
		return nil
	}

	return ut.Check()
}

func safeAddType(curr jschemaLib.Schema, n string, ut jschemaLib.Schema) error {
	err := curr.AddType(n, ut)
	var e interface{ Code() jerrors.ErrorCode }
	if errors.As(err, &e) && e.Code() == jerrors.ErrDuplicationOfNameOfTypes {
		err = nil
	}
	return err
}

func (core *JApiCore) checkUserType(name string) *jerr.JApiError {
	err := core.userTypes.GetValue(name).Check()
	if err == nil {
		return nil
	}

	d := core.catalog.GetRawUserTypes().GetValue(name)
	var e kit.Error
	if !errors.As(err, &e) {
		return d.KeywordError(err.Error())
	}

	if e.IncorrectUserType() != "" && e.IncorrectUserType() != name {
		return core.checkUserType(e.IncorrectUserType())
	}

	return d.BodyErrorIndex(e.Message(), e.Position())
}

func jschemaToJAPIError(err error, d *directive.Directive) *jerr.JApiError {
	var e kit.Error
	if errors.As(err, &e) {
		return d.BodyErrorIndex(e.Message(), e.Position())
	}
	return d.KeywordError(err.Error())
}

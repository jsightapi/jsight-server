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

func (core *JApiCore) buildCatalog() *jerr.JAPIError {
	if len(core.directivesWithPastes) != 0 && core.directivesWithPastes[0].Type() != directive.Jsight {
		return core.directivesWithPastes[0].KeywordError("JSIGHT should be the first directive")
	}

	return core.addDirectives()
}

func (core *JApiCore) compileUserTypes() *jerr.JAPIError {
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

func (core *JApiCore) buildUserTypes() *jerr.JAPIError {
	core.Catalog().GetRawUserTypes().EachSafe(func(k string, v *directive.Directive) {
		switch notation.SchemaNotation(v.Parameter("SchemaNotation")) {
		case "", notation.SchemaNotationJSight:
			if v.BodyCoords.IsSet() {
				core.userTypes.Set(k, jschema.New(k, v.BodyCoords.Read()))
			}
		case notation.SchemaNotationRegex:
			core.userTypes.Set(k, regex.New(k, v.BodyCoords.Read()))
		default:
			// nothing
		}
	})

	err := core.userTypes.Each(func(k string, v jschemaLib.Schema) error {
		if _, ok := core.processedUserTypes[k]; ok {
			// This user type already built, skip.
			return nil
		}

		dd := core.catalog.GetRawUserTypes()

		tt, err := core.getUsedUserTypes(v)
		if err != nil {
			return jschemaToJAPIError(err, dd.GetValue(k))
		}

		alreadyAddedTypes := map[string]struct{}{}

		for _, n := range tt {
			if n != k {
				if err := core.buildUserType(n); err != nil {
					return err
				}
			}

			ut := core.userTypes.GetValue(n)
			if ut == nil {
				continue
			}

			if _, ok := alreadyAddedTypes[n]; !ok {
				if err := safeAddType(v, n, ut); err != nil {
					return jschemaToJAPIError(err, dd.GetValue(n))
				}
				alreadyAddedTypes[n] = struct{}{}
			}
		}

		core.processedUserTypes[k] = struct{}{}
		core.userTypes.Set(k, v)

		return nil
	})
	return adoptError(err)
}

func adoptError(err error) (e *jerr.JAPIError) {
	if err == nil {
		return nil
	}

	if errors.As(err, &e) {
		return e
	}

	panic(fmt.Sprintf("Invalid error was given: %#v", err))
}

func (core *JApiCore) buildUserType(name string) *jerr.JAPIError {
	if _, ok := core.processedUserTypes[name]; ok {
		// This user type already processed, skip.
		return nil
	}

	currUT := core.userTypes.GetValue(name)
	if currUT == nil {
		return nil
	}

	dd := core.catalog.GetRawUserTypes()

	tt, err := core.getUsedUserTypes(currUT)
	if err != nil {
		return jschemaToJAPIError(err, dd.GetValue(name))
	}

	core.processedUserTypes[name] = struct{}{}
	alreadyAddedTypes := map[string]struct{}{}

	for _, n := range tt {
		if n != name {
			if err := core.buildUserType(n); err != nil {
				return err
			}
		}

		ut := core.userTypes.GetValue(n)
		if ut == nil {
			continue
		}

		if _, ok := alreadyAddedTypes[n]; !ok {
			if err := safeAddType(currUT, n, ut); err != nil {
				return jschemaToJAPIError(err, dd.GetValue(n))
			}
			alreadyAddedTypes[n] = struct{}{}
		}
	}

	core.userTypes.Set(name, currUT)

	return nil
}

func safeAddType(curr jschemaLib.Schema, n string, ut jschemaLib.Schema) error {
	err := curr.AddType(n, ut)
	var e interface{ Code() jerrors.ErrorCode }
	if errors.As(err, &e) && e.Code() == jerrors.ErrDuplicationOfNameOfTypes {
		err = nil
	}
	return err
}

func (core *JApiCore) getUsedUserTypes(ut jschemaLib.Schema) ([]string, error) {
	alreadyProcessed := map[string]struct{}{}
	if err := core.fetchUsedUserTypes(ut, alreadyProcessed); err != nil {
		return nil, err
	}

	ss := make([]string, 0, len(alreadyProcessed))
	for s := range alreadyProcessed {
		ss = append(ss, s)
	}

	return ss, nil
}

func (core *JApiCore) fetchUsedUserTypes(
	ut jschemaLib.Schema,
	alreadyProcessed map[string]struct{},
) error {
	if ut == nil {
		return nil
	}

	tt, err := ut.UsedUserTypes()
	if err != nil {
		return err
	}

	if len(tt) == 0 {
		return nil
	}

	for _, t := range tt {
		if _, ok := alreadyProcessed[t]; ok {
			continue
		}

		alreadyProcessed[t] = struct{}{}
		if err := core.fetchUsedUserTypes(core.userTypes.GetValue(t), alreadyProcessed); err != nil {
			return err
		}
	}
	return nil
}

func (core *JApiCore) checkUserType(name string) *jerr.JAPIError {
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

func jschemaToJAPIError(err error, d *directive.Directive) *jerr.JAPIError {
	var e kit.Error
	if errors.As(err, &e) {
		return d.BodyErrorIndex(e.Message(), e.Position())
	}
	return d.KeywordError(err.Error())
}

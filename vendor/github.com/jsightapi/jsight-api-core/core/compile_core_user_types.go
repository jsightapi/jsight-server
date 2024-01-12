package core

import (
	stdErrors "errors"

	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/notations/jschema"
	"github.com/jsightapi/jsight-schema-core/notations/regex"

	"github.com/jsightapi/jsight-api-core/directive"
	"github.com/jsightapi/jsight-api-core/jerr"
	"github.com/jsightapi/jsight-api-core/notation"
)

func (core *JApiCore) collectUserTypes() *jerr.JApiError {
	core.collectRawUserTypes()
	return core.compileUserTypes()
}

func (core *JApiCore) collectRawUserTypes() {
	for _, d := range core.directivesWithPastes {
		if d.Type() == directive.Type {
			core.AddRawUserType(d)
		}
	}
}

func (core *JApiCore) compileUserTypes() *jerr.JApiError {
	// Two-phase algorithm. On the first step we just create schema for each user
	// type. On the second step we will add all schema to all.
	// This is the simplest solution which allows us to skip building dependency
	// graph between types.

	if err := core.buildUserTypes(); err != nil {
		return err
	}

	err := core.userTypes.Each(func(k string, _ schema.Schema) error {
		return core.checkUserType(k)
	})
	return adoptError(err)
}

func (core *JApiCore) buildUserTypes() *jerr.JApiError {
	err := core.rawUserTypes.Each(func(k string, d *directive.Directive) error {
		switch notation.SchemaNotation(d.NamedParameter("SchemaNotation")) {
		case "", notation.SchemaNotationJSight:
			if !d.BodyCoords.IsSet() {
				return d.KeywordError(jerr.BodyIsEmpty)
			}
			core.userTypes.Set(k, jschema.New(k, d.BodyCoords.Read()))
		case notation.SchemaNotationRegex:
			if !d.BodyCoords.IsSet() {
				return d.KeywordError(jerr.BodyIsEmpty)
			}
			core.userTypes.Set(k, regex.New(k, d.BodyCoords.Read()))
		default:
			// nothing
		}
		return nil
	})
	if err != nil {
		return adoptError(err)
	}

	err = core.userTypes.Each(func(n string, _ schema.Schema) error {
		return core.compileUserTypeWithAllDependencies(n)
	})
	return adoptError(err)
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

	dd := core.rawUserTypes

	// Add rules before we try to do something with the type.
	for n, r := range core.rules {
		if err := currUT.AddRule(n, r); err != nil {
			return jschemaToJAPIError(err, dd.GetValue(n))
		}
	}

	tt, err := fetchUsedUserTypes(currUT, core.userTypes)
	if err != nil {
		var ute userTypeError
		if stdErrors.As(err, &ute) {
			name = ute.userTypeName
			err = ute.err
		}
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

func (core *JApiCore) checkUserTypeDuringBuild(name string, ut schema.Schema) error {
	// In order to prevent errors in type recursion.
	if _, ok := core.processedUserTypes[name]; ok {
		return nil
	}

	return ut.Check()
}

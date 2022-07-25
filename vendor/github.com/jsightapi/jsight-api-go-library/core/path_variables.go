package core

import (
	"fmt"
	"strings"

	jschema "github.com/jsightapi/jsight-schema-go-library"

	"github.com/jsightapi/jsight-api-go-library/catalog"
	"github.com/jsightapi/jsight-api-go-library/directive"
	"github.com/jsightapi/jsight-api-go-library/jerr"
	"github.com/jsightapi/jsight-api-go-library/notation"
)

type prop struct {
	parameter           string
	schemaContentJSight *catalog.SchemaContentJSight
	directive           directive.Directive
}

func (core *JApiCore) newPathVariables(properties []prop) (*catalog.PathVariables, *jerr.JApiError) {
	s := catalog.NewSchema(notation.SchemaNotationJSight)
	s.ContentJSight = &catalog.SchemaContentJSight{
		IsKeyUserTypeRef: false,
		TokenType:        jschema.JSONTypeObject,
		Type:             jschema.JSONTypeObject,
		Optional:         false,
		Children:         make([]*catalog.SchemaContentJSight, 0, len(properties)),
	}

	for _, p := range properties {
		if err := core.collectUsedUserTypes(p.schemaContentJSight, s.UsedUserTypes /* &s.UsedUserEnums */); err != nil {
			return nil, p.directive.KeywordError(err.Error())
		}

		p.schemaContentJSight.Key = catalog.SrtPtr(p.parameter)
		s.ContentJSight.Children = append(s.ContentJSight.Children, p.schemaContentJSight)
	}

	return &catalog.PathVariables{Schema: s}, nil
}

func (core *JApiCore) collectUsedUserTypes(sc *catalog.SchemaContentJSight, usedUserTypes *catalog.StringSet) error {
	if sc.TokenType == jschema.JSONTypeShortcut {
		// We have two different cases under "reference" type:
		// 1. Single type like "@foo"
		// 2. A list of types like "@foo | @bar"
		//
		// For the first case we have valid user type in the `v.Type` property.
		// But for the second case we got "mixed" there. So we should use `v.ScalarValue`
		// instead. This property should always be string.
		for _, t := range strings.Split(sc.ScalarValue, " | ") {
			if err := core.appendUsedUserType(usedUserTypes, t); err != nil {
				return err
			}
		}
	} else {
		err := sc.Rules.Each(func(k string, v catalog.Rule) error {
			switch k {
			case "type":
				if v.ScalarValue[0] == '@' {
					if err := core.appendUsedUserType(usedUserTypes, v.ScalarValue); err != nil {
						return err
					}
				}

			case "or":
				for _, i := range v.Children {
					var userType string
					if i.ScalarValue != "" {
						userType = i.ScalarValue
					} else {
						for _, c := range i.Children {
							if c.Key == "type" {
								userType = c.ScalarValue
								break
							}
						}
						// TODO if not found ???
					}

					// Schema types shouldn't be added.
					if jschema.IsValidType(userType) {
						continue
					}

					if err := core.appendUsedUserType(usedUserTypes, userType); err != nil {
						return err
					}
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (core *JApiCore) appendUsedUserType(usedUserTypes *catalog.StringSet, s string) error {
	if t, ok := core.catalog.UserTypes.Get(s); ok {
		switch t.Schema.Notation {
		case notation.SchemaNotationJSight:
			switch t.Schema.ContentJSight.TokenType {
			case "string", "number", "boolean", "null":
				usedUserTypes.Add(s)
				return nil
			default:
				return fmt.Errorf(`unavailable JSON type "%s" of the UserType "%s" in Path directive`, t.Schema.ContentJSight.TokenType, s)
			}
		case notation.SchemaNotationRegex:
			usedUserTypes.Add(s)
			return nil
		default:
			// case notation.SchemaNotationAny, notation.SchemaNotationEmpty:
			// return err (see below)
		}
	}
	return fmt.Errorf(`UserType not found "%s" for Path directive`, s)
}

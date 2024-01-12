package core

import (
	"errors"
	"fmt"
	"strings"

	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/notations/jschema"
	"github.com/jsightapi/jsight-schema-core/notations/jschema/ischema"

	"github.com/jsightapi/jsight-api-core/catalog"
	"github.com/jsightapi/jsight-api-core/jerr"
)

func (core *JApiCore) compileCatalog() *jerr.JApiError {
	je := core.collectPiecesOfPathVariables()
	if je != nil {
		return je
	}

	return core.setPathVariablesToCatalog()
}

func (core *JApiCore) collectPiecesOfPathVariables() *jerr.JApiError {
	ut := core.UserTypesData()

	for i := 0; i < len(core.rawPathVariables); i++ {
		if err := core.checkPathSchema(core.rawPathVariables[i].schema); err != nil {
			return core.rawPathVariables[i].pathDirective.KeywordError(err.Error())
		}

		obj := catalog.JSchemaObject{
			JSchema: core.rawPathVariables[i].schema,
		}
		schemaProps := obj.ObjectFirstLevelProperties(ut)
		types := core.rawPathVariables[i].schema.InnerTypesList()

		for _, pp := range core.rawPathVariables[i].parameters {
			if v, ok := schemaProps[pp.parameter]; ok {
				if _, ok := core.piecesOfPathVariables[pp]; ok {
					return core.rawPathVariables[i].pathDirective.KeywordError(
						fmt.Sprintf("The parameter %q has already been defined earlier",
							pp.path))
				}

				core.piecesOfPathVariables[pp] = PieceOfPathVariable{
					node:  v,
					types: types,
				}
			}

			delete(schemaProps, pp.parameter)
		}

		if len(schemaProps) != 0 {
			ss := core.getPropertiesNames(schemaProps)
			return core.rawPathVariables[i].pathDirective.KeywordError(
				fmt.Sprintf("Has unused parameters %q in schema", ss))
		}
	}

	return nil
}

func (core *JApiCore) setPathVariablesToCatalog() *jerr.JApiError {
	err := core.catalog.Interactions.Map(
		func(_ catalog.InteractionID, v catalog.Interaction) (catalog.Interaction, error) {
			if hi, ok := v.(*catalog.HTTPInteraction); ok {
				pp := pathParameters(v.Path().String())

				b := catalog.NewPathVariablesBuilder(core.catalog.UserTypes)
				for _, p := range pp {
					if piece, ok := core.piecesOfPathVariables[p]; ok {
						b.AddProperty(p.parameter, piece.node.Copy(), piece.types)
					}
				}
				if b.Len() != 0 {
					hi.SetPathVariables(b.Build())
				}
			}
			return v, nil
		},
	)
	if err != nil {
		return err.(*jerr.JApiError) //nolint:errorlint
	}

	return nil
}

func (core *JApiCore) checkPathSchema(s *jschema.JSchema) error {
	if err := core.checkPathSchemaRoot(s); err != nil {
		return err
	}

	for i := range s.ASTNode.Children {
		if err := core.checkPathSchemaProperty(s.ASTNode.Children[i]); err != nil {
			return err
		}
	}

	return nil
}

func (core *JApiCore) checkPathSchemaRoot(s *jschema.JSchema) error {
	if s.ASTNode.TokenType == schema.TokenTypeShortcut {
		typeName := s.ASTNode.SchemaType
		if typeName == "mixed" {
			return errors.New(jerr.PathOrErr)
		}

		ut, ok := core.catalog.UserTypes.Get(typeName)
		if !ok {
			return fmt.Errorf(`%s (%s)`, jerr.UserTypeNotFound, typeName)
		}

		return core.checkPathSchemaRoot(ut.Schema.(*catalog.ExchangeJSightSchema).JSchema)
	}

	if s.ASTNode.TokenType != schema.TokenTypeObject {
		return errors.New(jerr.PathObjectErr)
	}

	if s.ASTNode.Rules.Has("additionalProperties") {
		return errors.New(jerr.PathAdditionalPropertiesErr)
	}

	if s.ASTNode.Rules.Has("nullable") {
		return errors.New(jerr.PathNullableErr)
	}

	if s.ASTNode.Rules.Has("or") {
		return errors.New(jerr.PathOrErr)
	}

	if len(s.ASTNode.Children) == 0 && !s.ASTNode.Rules.Has("allOf") {
		return errors.New(jerr.PathEmptyErr)
	}

	if s.ASTNode.Rules.Has("allOf") {
		allOf, _ := s.ASTNode.Rules.Get("allOf")

		switch allOf.TokenType { //nolint:exhaustive // We expects only this types.
		case schema.TokenTypeArray:
			for i := len(allOf.Items) - 1; i >= 0; i-- {
				r := allOf.Items[i]
				if err := core.checkPathSchemaPropertyInAllOf(r.Value); err != nil {
					return err
				}
			}
		case schema.TokenTypeShortcut:
			if err := core.checkPathSchemaPropertyInAllOf(allOf.Value); err != nil {
				return err
			}
		}
	}

	return nil
}

func (core *JApiCore) checkPathSchemaPropertyInAllOf(typeName string) error {
	ut, ok := core.catalog.UserTypes.Get(typeName)
	if !ok {
		return fmt.Errorf(`%s (%s)`, jerr.UserTypeNotFound, typeName)
	}

	es, ok := ut.Schema.(*catalog.ExchangeJSightSchema)
	if !ok {
		return fmt.Errorf(`%s (%s)`, jerr.UserTypeNotFound, typeName)
	}

	if err := core.checkPathSchema(es.JSchema); err != nil {
		return err
	}

	return nil
}

func (core *JApiCore) checkPathSchemaProperty(an schema.ASTNode) error {
	if an.TokenType == schema.TokenTypeObject || an.TokenType == schema.TokenTypeArray {
		return fmt.Errorf("%s (%s)", jerr.PathMultiLevelPropertyErr, an.Key)
	}

	rule, ok := an.Rules.Get("or")
	if ok {
		for _, v := range rule.Items {
			switch v.TokenType {
			case schema.TokenTypeShortcut:
				if err := core.checkPathSchemaPropertyUserType(v.Value); err != nil {
					return err
				}
			case schema.TokenTypeObject:
				if t, ok := v.Properties.Get("type"); ok {
					if t.TokenType == "string" &&
						(t.Value == schema.TokenTypeObject || t.Value == schema.TokenTypeArray) {
						return fmt.Errorf("%s (%s)", jerr.PathMultiLevelPropertyErr, an.Key)
					}
				}
			}
		}
	} else if an.TokenType == schema.TokenTypeShortcut {
		if err := core.checkPathSchemaPropertyUserType(an.Value); err != nil {
			return err
		}
	}

	return nil
}

func (core *JApiCore) checkPathSchemaPropertyUserType(typeName string) error {
	ut, ok := core.catalog.UserTypes.Get(typeName)
	if !ok {
		return fmt.Errorf(`%s (%s)`, jerr.UserTypeNotFound, typeName)
	}

	rootNode, err := ut.Schema.GetAST()
	if err != nil {
		return errors.New(jerr.RuntimeFailure)
	}

	if err := core.checkPathSchemaProperty(rootNode); err != nil {
		return err
	}

	return nil
}

func (*JApiCore) getPropertiesNames(m map[string]ischema.Node) string {
	if len(m) == 0 {
		return ""
	}

	buf := strings.Builder{}
	for k := range m {
		buf.WriteString(k)
		buf.WriteString(", ")
	}
	return strings.TrimSuffix(buf.String(), ", ")
}

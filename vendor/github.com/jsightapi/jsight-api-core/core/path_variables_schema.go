package core

import (
	"github.com/jsightapi/jsight-schema-core/bytes"
	"github.com/jsightapi/jsight-schema-core/notations/jschema"
	"github.com/jsightapi/jsight-schema-core/notations/jschema/loader"
	"github.com/jsightapi/jsight-schema-core/notations/jschema/scanner"
	"github.com/jsightapi/jsight-schema-core/panics"
)

type pSchema struct {
	*jschema.JSchema
}

func newPathVariablesSchema(content bytes.Bytes, userTypes map[string]*jschema.JSchema) (*pSchema, error) {
	s := &pSchema{
		JSchema: jschema.New("", content),
	}

	err := s.loadPathVariables()
	if err != nil {
		return nil, err
	}

	for k, v := range userTypes {
		if err = s.AddType(k, v); err != nil { //nolint:gocritic
			return nil, err
		}
	}

	err = s.compilePathVariables()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *pSchema) loadPathVariables() error {
	return s.LoadOnce.Do(func() (err error) {
		defer func() {
			err = panics.Handle(recover(), err)
		}()
		sc := loader.LoadSchemaWithoutCompile(
			scanner.New(s.File),
			nil,
			s.Rules,
		)
		s.Inner = &sc
		s.ASTNode = s.BuildASTNode()
		return nil
	})
}

func (s *pSchema) compilePathVariables() error {
	return s.CompileOnce.Do(func() (err error) {
		defer func() {
			err = panics.Handle(recover(), err)
		}()
		loader.CompileAllOf(s.Inner)
		return nil
	})
}

package jschema

import (
	"fmt"

	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/bytes"
	"github.com/jsightapi/jsight-schema-core/errs"
	"github.com/jsightapi/jsight-schema-core/fs"
	"github.com/jsightapi/jsight-schema-core/internal/sync"
	"github.com/jsightapi/jsight-schema-core/kit"
	"github.com/jsightapi/jsight-schema-core/notations/jschema/checker"
	"github.com/jsightapi/jsight-schema-core/notations/jschema/ischema"
	"github.com/jsightapi/jsight-schema-core/notations/jschema/loader"
	"github.com/jsightapi/jsight-schema-core/notations/jschema/scanner"
	"github.com/jsightapi/jsight-schema-core/notations/regex"
	"github.com/jsightapi/jsight-schema-core/panics"
)

type JSchema struct {
	File  *fs.File
	Inner *ischema.ISchema

	Rules map[string]schema.Rule

	UsedUserTypes_ *StringSet

	LenOnce     sync.ErrOnceWithValue[uint]
	LoadOnce    sync.ErrOnce
	CompileOnce sync.ErrOnce

	ASTNode                  schema.ASTNode
	AreKeysOptionalByDefault bool
}

var _ schema.Schema = (*JSchema)(nil)

// New creates a Jsight schema with specified name and content.
func New[T bytes.ByteKeeper](name string, content T, oo ...Option) *JSchema {
	return FromFile(fs.NewFile(name, content), oo...)
}

// FromFile creates a Jsight schema from file.
func FromFile(f *fs.File, oo ...Option) *JSchema {
	s := &JSchema{
		File:           f,
		Rules:          map[string]schema.Rule{},
		UsedUserTypes_: &StringSet{},
	}

	for _, o := range oo {
		o(s)
	}

	return s
}

func FromRSchema(s *regex.RSchema) (*JSchema, error) {
	pattern, err := s.Pattern()
	if err != nil {
		return nil, err
	}

	example, err := s.Example()
	if err != nil {
		return nil, errs.ErrRegexExample.F(err)
	}

	ss := New(s.File.Name(), fmt.Sprintf("%q // {regex: %q}", example, pattern))
	if err = ss.load(); err != nil {
		return nil, errs.ErrLoadError.F(err)
	}

	return ss, nil
}

type Option func(s *JSchema)

func (s *JSchema) Len() (uint, error) {
	return s.LenOnce.Do(func() (uint, error) {
		return s.computeLen()
	})
}

func (s *JSchema) computeLen() (length uint, err error) {
	// Iterate through all lexemes until we reach the end
	// We should rewind here in case we call NextLexeme method.
	defer func() {
		err = panics.Handle(recover(), err)
	}()

	return scanner.New(s.File, scanner.ComputeLength).Length(), err
}

func (s *JSchema) Example() (b []byte, err error) {
	defer func() {
		err = panics.Handle(recover(), err)
	}()

	if err := s.Compile(); err != nil {
		return nil, err
	}

	if s.Inner.RootNode() == nil {
		return nil, kit.NewJSchemaError(s.File, errs.ErrEmptySchema.F())
	}

	return newExampleBuilder(s.Inner.TypesList()).Build(s.Inner.RootNode())
}

func (s *JSchema) AddType(name string, sc schema.Schema) (err error) {
	defer func() {
		err = panics.Handle(recover(), err)
	}()

	if err := s.load(); err != nil {
		return err
	}

	switch typ := sc.(type) {
	case *JSchema:
		if err := typ.load(); err != nil {
			return errs.ErrLoadError.F(err)
		}

		s.Inner.AddNamedType(name, typ.Inner, s.File, 0)
	case *regex.RSchema:
		typSc, err := FromRSchema(typ)
		if err != nil {
			return err
		}

		s.Inner.AddNamedType(name, typSc.Inner, s.File, 0)

	default:
		return errs.ErrRuntimeFailure.F()
	}

	return nil
}

func (s *JSchema) AddRule(n string, r schema.Rule) error {
	if s.Inner != nil {
		return errs.ErrRuleIsAlreadyCompiled.F()
	}

	if r == nil {
		return errs.ErrRuleIsNil.F()
	}

	if err := r.Check(); err != nil {
		return err
	}
	s.Rules[n] = r
	return nil
}

func (s *JSchema) Check() (err error) {
	defer func() {
		err = panics.Handle(recover(), err)
	}()
	return s.Compile()
}

func (s *JSchema) GetAST() (an schema.ASTNode, err error) {
	if err := s.Compile(); err != nil {
		return schema.ASTNode{}, err
	}

	return s.ASTNode, nil
}

func (s *JSchema) UsedUserTypes() ([]string, error) {
	if err := s.load(); err != nil {
		return nil, err
	}
	return s.UsedUserTypes_.Data(), nil
}

func (s *JSchema) AddUserTypeName(name string) {
	s.UsedUserTypes_.Add(name)
}

func (s *JSchema) load() error {
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
		s.CollectUserTypes()
		loader.CompileBasic(s.Inner, s.AreKeysOptionalByDefault)
		return nil
	})
}

func (s *JSchema) CollectUserTypes() {
	node := s.Inner.RootNode()
	// This is possible when schema isn't valid.
	if node == nil {
		return
	}

	for _, str := range collectUserTypes(node) {
		s.UsedUserTypes_.Add(str)
	}
}

func (s *JSchema) BuildASTNode() schema.ASTNode {
	root := s.Inner.RootNode()
	if root == nil {
		// This case will be handled in loader.CompileBasic.
		return schema.ASTNode{
			Rules: &schema.RuleASTNodes{},
		}
	}

	an, err := root.ASTNode()
	if err != nil {
		panic(err)
	}
	return an
}

func (s *JSchema) Compile() error {
	return s.CompileOnce.Do(func() (err error) {
		defer func() {
			err = panics.Handle(recover(), err)
		}()
		if err := s.load(); err != nil {
			return err
		}
		loader.CompileAllOf(s.Inner)
		loader.AddUnnamedTypes(s.Inner)
		checker.CheckRootSchema(s.Inner)
		return checker.CheckRecursion(s.File.Name(), s.Inner)
	})
}

func (s *JSchema) InnerTypesList() map[string]ischema.Type {
	return s.Inner.TypesList()
}

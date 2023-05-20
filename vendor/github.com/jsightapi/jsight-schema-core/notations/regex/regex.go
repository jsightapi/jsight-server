package regex

import (
	"regexp"

	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/bytes"
	"github.com/jsightapi/jsight-schema-core/errs"
	"github.com/jsightapi/jsight-schema-core/fs"
	"github.com/jsightapi/jsight-schema-core/internal/sync"
	"github.com/jsightapi/jsight-schema-core/kit"

	"github.com/lucasjones/reggen"
)

type RSchema struct {
	File *fs.File
	RE   *regexp.Regexp

	pattern       string
	compileOnce   sync.ErrOnce
	generatorOnce sync.ErrOnceWithValue[*reggen.Generator]
	generatorSeed int64
}

var _ schema.Schema = &RSchema{}

type Option func(*RSchema)

// WithGeneratorSeed pass specific seed to regex example generator.
// Necessary for test.
func WithGeneratorSeed(seed int64) Option {
	return func(s *RSchema) {
		s.generatorSeed = seed
	}
}

// New creates a Regex schema with specified name and content.
func New[T bytes.ByteKeeper](name string, content T, oo ...Option) *RSchema {
	return FromFile(fs.NewFile(name, content), oo...)
}

// FromFile creates a Regex schema from file.
func FromFile(f *fs.File, oo ...Option) *RSchema {
	s := &RSchema{
		File: f,
	}

	for _, o := range oo {
		o(s)
	}

	return s
}

func (s *RSchema) Pattern() (string, error) {
	if err := s.Compile(); err != nil {
		return "", err
	}
	return s.pattern, nil
}

func (s *RSchema) Len() (uint, error) {
	if err := s.Compile(); err != nil {
		return 0, err
	}
	// Add 2 for beginning and ending '/' character.
	return uint(len(s.pattern)) + 2, nil
}

func (s *RSchema) Example() ([]byte, error) {
	if err := s.Compile(); err != nil {
		return nil, err
	}

	return s.generateExample()
}

func (s *RSchema) generateExample() ([]byte, error) {
	g, err := s.generatorOnce.Do(func() (*reggen.Generator, error) {
		g, err := reggen.NewGenerator(s.pattern)
		if err != nil {
			return nil, err
		}
		g.SetSeed(s.generatorSeed)
		return g, nil
	})
	if err != nil {
		return nil, err
	}

	return []byte(g.Generate(1)), nil
}

func (*RSchema) AddType(string, schema.Schema) error {
	// Regex doesn't use any user types at all.
	return nil
}

func (*RSchema) AddRule(string, schema.Rule) error {
	// Regex doesn't use any rules at all.
	return nil
}

func (s *RSchema) Check() error {
	return s.Compile()
}

func (s *RSchema) GetAST() (schema.ASTNode, error) {
	if err := s.Compile(); err != nil {
		return schema.ASTNode{}, err
	}
	return schema.ASTNode{
		IsKeyShortcut: false,
		TokenType:     schema.TokenTypeString,
		SchemaType:    string(schema.SchemaTypeString),
		Rules:         &schema.RuleASTNodes{},
		Value:         "/" + s.pattern + "/",
	}, nil
}

func (*RSchema) UsedUserTypes() ([]string, error) {
	// Regex doesn't use any user types at all.
	return nil, nil
}

func (s *RSchema) Compile() error {
	return s.compileOnce.Do(func() error {
		return s.doCompile()
	})
}

func (s *RSchema) doCompile() error {
	content := s.File.Content()

	if content.Byte(0) != '/' {
		return s.newJSchemaError(errs.ErrRegexUnexpectedStart, 0, content.Byte(0))
	}

	var escaped bool

loop:
	for i, c := range content.SubLow(1).Data() {
		switch c {
		case '\\':
			escaped = !escaped

		case '/':
			if !escaped {
				s.pattern = content.Sub(1, i+1).String()
				break loop
			}
			escaped = false

		default:
			escaped = false
		}
	}

	if s.pattern == "" {
		idx := uint(content.Len() - 1)
		return s.newJSchemaError(errs.ErrRegexUnexpectedEnd, idx, content.Byte(idx))
	}

	var err error

	if s.RE, err = regexp.Compile(s.pattern); err != nil {
		e := errs.ErrRegexInvalid.F(content)
		err := kit.NewJSchemaError(s.File, e)
		err.SetIndex(bytes.Index(0))
		return err
	}
	return nil
}

func (s *RSchema) newJSchemaError(code errs.Code, idx uint, c byte) kit.JSchemaError {
	e := code.F(bytes.QuoteChar(c))
	err := kit.NewJSchemaError(s.File, e)
	err.SetIndex(bytes.Index(idx))
	return err
}

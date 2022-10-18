package regex

import (
	stdErrors "errors"
	"regexp"

	"github.com/lucasjones/reggen"

	jschema "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/bytes"
	"github.com/jsightapi/jsight-schema-go-library/errors"
	"github.com/jsightapi/jsight-schema-go-library/fs"
	"github.com/jsightapi/jsight-schema-go-library/internal/sync"
)

type Schema struct {
	file *fs.File

	pattern       string
	compileOnce   sync.ErrOnce
	generatorOnce sync.ErrOnceWithValue[*reggen.Generator]
	generatorSeed int64
}

var _ jschema.Schema = &Schema{}

type Option func(*Schema)

// WithGeneratorSeed pass specific seed to regex example generator.
// Necessary for test.
func WithGeneratorSeed(seed int64) Option {
	return func(s *Schema) {
		s.generatorSeed = seed
	}
}

// New creates a Regex schema with specified name and content.
func New[T fs.FileContent](name string, content T, oo ...Option) *Schema {
	return FromFile(fs.NewFile(name, content), oo...)
}

// FromFile creates a Regex schema from file.
func FromFile(f *fs.File, oo ...Option) *Schema {
	s := &Schema{
		file: f,
	}

	for _, o := range oo {
		o(s)
	}

	return s
}

func (s *Schema) Pattern() (string, error) {
	if err := s.compile(); err != nil {
		return "", err
	}
	return s.pattern, nil
}

func (s *Schema) Len() (uint, error) {
	if err := s.compile(); err != nil {
		return 0, err
	}
	// Add 2 for beginning and ending '/' character.
	return uint(len(s.pattern)) + 2, nil
}

func (s *Schema) Example() ([]byte, error) {
	if err := s.compile(); err != nil {
		return nil, err
	}

	return s.generateExample()
}

func (s *Schema) generateExample() ([]byte, error) {
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

func (*Schema) AddType(string, jschema.Schema) error {
	// Regex doesn't use any user types at all.
	return nil
}

func (*Schema) AddRule(string, jschema.Rule) error {
	// Regex doesn't use any rules at all.
	return nil
}

func (s *Schema) Check() error {
	return s.compile()
}

func (*Schema) Validate(jschema.Document) error {
	return stdErrors.New("unimplemented")
}

func (s *Schema) GetAST() (jschema.ASTNode, error) {
	if err := s.compile(); err != nil {
		return jschema.ASTNode{}, err
	}
	return jschema.ASTNode{
		IsKeyShortcut: false,
		TokenType:     jschema.TokenTypeString,
		SchemaType:    string(jschema.SchemaTypeString),
		Rules:         nil,
		Value:         "/" + s.pattern + "/",
	}, nil
}

func (*Schema) UsedUserTypes() ([]string, error) {
	// Regex doesn't use any user types at all.
	return nil, nil
}

func (s *Schema) compile() error {
	return s.compileOnce.Do(func() error {
		return s.doCompile()
	})
}

func (s *Schema) doCompile() error {
	content := s.file.Content()

	if content[0] != '/' {
		return s.newDocumentError(errors.ErrRegexUnexpectedStart, 0, content[0])
	}

	var escaped bool

loop:
	for i, c := range content[1:] {
		switch c {
		case '\\':
			escaped = !escaped

		case '/':
			if !escaped {
				s.pattern = string(content[1 : i+1])
				break loop
			}
			escaped = false

		default:
			escaped = false
		}
	}

	if s.pattern == "" {
		idx := uint(len(content) - 1)
		return s.newDocumentError(errors.ErrRegexUnexpectedEnd, idx, content[idx])
	}

	if _, err := regexp.Compile(s.pattern); err != nil {
		e := errors.Format(errors.ErrRegexInvalid, content)
		err := errors.NewDocumentError(s.file, e)
		err.SetIndex(bytes.Index(0))
		return err
	}
	return nil
}

func (s *Schema) newDocumentError(code errors.ErrorCode, idx uint, c byte) errors.DocumentError {
	e := errors.Format(code, bytes.QuoteChar(c))
	err := errors.NewDocumentError(s.file, e)
	err.SetIndex(bytes.Index(idx))
	return err
}

package enum

import (
	stdErrors "errors"
	"sync"

	jschema "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/bytes"
	"github.com/jsightapi/jsight-schema-go-library/fs"
	"github.com/jsightapi/jsight-schema-go-library/internal/lexeme"
)

// The Enum rule.
type Enum struct {
	compileErr       error
	computeLengthErr error
	buildAStNodeErr  error

	file   *fs.File
	values []Value

	astNode jschema.ASTNode
	length  uint

	compileOnce       sync.Once
	computeLengthOnce sync.Once
	buildASTNodeOnce  sync.Once
}

// Value represents single enum's value.
type Value struct {
	// Comment value's comment.
	Comment string

	// Type value type.
	Type jschema.SchemaType

	// Value enum value.
	Value bytes.Bytes
}

var _ jschema.Rule = (*Enum)(nil)

// New creates new Enum rule with specified name and content.
func New(name string, content []byte) *Enum {
	return FromFile(fs.NewFile(name, content))
}

// FromFile creates Enum rule from specified file.
func FromFile(f *fs.File) *Enum {
	return &Enum{file: f}
}

func (e *Enum) Len() (uint, error) {
	e.computeLengthOnce.Do(func() {
		e.length, e.computeLengthErr = newScanner(e.file, scannerComputeLength).Length()
	})
	return e.length, e.computeLengthErr
}

// Check checks that enum is valid.
func (e *Enum) Check() error {
	return e.compile()
}

func (e *Enum) GetAST() (jschema.ASTNode, error) {
	if err := e.compile(); err != nil {
		return jschema.ASTNode{}, err
	}

	return e.buildASTNode()
}

func (e *Enum) buildASTNode() (jschema.ASTNode, error) {
	e.buildASTNodeOnce.Do(func() {
		e.astNode = jschema.ASTNode{
			JSONType:   jschema.JSONTypeArray,
			SchemaType: string(jschema.SchemaTypeEnum),
		}

		if len(e.values) == 0 {
			return
		}

		e.astNode.Children = make([]jschema.ASTNode, 0, len(e.values))

		for _, v := range e.values {
			n := jschema.ASTNode{
				Value:   v.Value.String(),
				Comment: v.Comment,
			}

			if v.Value == nil {
				n.JSONType = jschema.JSONTypeNull
				n.SchemaType = string(jschema.SchemaTypeComment)
			} else {
				n.JSONType = v.Type.ToTokenType()
				n.SchemaType = string(v.Type)
			}

			e.astNode.Children = append(e.astNode.Children, n)
		}
	})
	return e.astNode, e.buildAStNodeErr
}

// Values returns a list of values defined in this enum.
func (e *Enum) Values() ([]Value, error) {
	if err := e.compile(); err != nil {
		return nil, err
	}
	return e.values, nil
}

func (e *Enum) compile() error {
	e.compileOnce.Do(func() {
		e.compileErr = e.doCompile()
	})
	return e.compileErr
}

func (e *Enum) doCompile() (err error) {
	scan := newScanner(e.file)

	collectLiteral := false
	for {
		lex, err := scan.Next()
		if stdErrors.Is(err, errEOS) {
			break
		}
		if err != nil {
			return err
		}

		// Collect enum values.
		switch lex.Type() {
		case lexeme.LiteralEnd:
			collectLiteral = true
			v := lex.Value()
			t, err := jschema.GuessSchemaType(v)
			if err != nil {
				return err
			}

			e.values = append(e.values, Value{
				Value: v,
				Type:  t,
			})

		case lexeme.NewLine:
			collectLiteral = false

		case lexeme.InlineAnnotationTextEnd:
			comment := lex.Value().TrimSpaces().String()
			if collectLiteral {
				e.values[len(e.values)-1].Comment = comment
			} else {
				e.values = append(e.values, Value{
					Comment: comment,
				})
			}
		}
	}
	return nil
}

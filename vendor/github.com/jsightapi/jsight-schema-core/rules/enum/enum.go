package enum

import (
	stdErrors "errors"

	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/bytes"
	"github.com/jsightapi/jsight-schema-core/fs"
	"github.com/jsightapi/jsight-schema-core/internal/sync"
	"github.com/jsightapi/jsight-schema-core/lexeme"
)

// The Enum rule.
type Enum struct {
	file   *fs.File
	values []Value

	compileOnce       sync.ErrOnce
	computeLengthOnce sync.ErrOnceWithValue[uint]
	buildASTNodeOnce  sync.ErrOnceWithValue[schema.ASTNode]
}

// Value represents single enum's value.
type Value struct {
	// Comment value's comment.
	Comment string

	// Type value type.
	Type schema.SchemaType

	// Value enum value.
	Value bytes.Bytes
}

var _ schema.Rule = (*Enum)(nil)

// New creates new Enum rule with specified name and content.
func New[T bytes.ByteKeeper](name string, content T) *Enum {
	return FromFile(fs.NewFile(name, content))
}

// FromFile creates Enum rule from specified file.
func FromFile(f *fs.File) *Enum {
	return &Enum{file: f}
}

func (e *Enum) Len() (uint, error) {
	return e.computeLengthOnce.Do(func() (uint, error) {
		return newScanner(e.file, scannerComputeLength).Length()
	})
}

// Check checks that enum is valid.
func (e *Enum) Check() error {
	return e.compile()
}

func (e *Enum) GetAST() (schema.ASTNode, error) {
	if err := e.compile(); err != nil {
		return schema.ASTNode{}, err
	}

	return e.buildASTNode()
}

func (e *Enum) buildASTNode() (schema.ASTNode, error) {
	return e.buildASTNodeOnce.Do(func() (schema.ASTNode, error) {
		an := schema.ASTNode{
			TokenType:  schema.TokenTypeArray,
			SchemaType: string(schema.SchemaTypeEnum),
		}

		if len(e.values) == 0 {
			return an, nil
		}

		an.Children = make([]schema.ASTNode, 0, len(e.values))

		for _, v := range e.values {
			n := schema.ASTNode{
				Value:   v.Value.String(),
				Comment: v.Comment,
			}

			if v.Value.IsNil() {
				n.TokenType = schema.TokenTypeNull
				n.SchemaType = string(schema.SchemaTypeComment)
			} else {
				n.TokenType = v.Type.ToTokenType()
				n.SchemaType = string(v.Type)
			}

			an.Children = append(an.Children, n)
		}

		return an, nil
	})
}

// Values returns a list of values defined in this enum.
func (e *Enum) Values() ([]Value, error) {
	if err := e.compile(); err != nil {
		return nil, err
	}
	return e.values, nil
}

func (e *Enum) compile() error {
	return e.compileOnce.Do(func() error {
		return e.doCompile()
	})
}

func (e *Enum) doCompile() (err error) {
	scan := newScanner(e.file)

	collectLiteral := false
	inAnnotation := false
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
			if err := e.handleLiteralEnd(lex); err != nil {
				return err
			}

		case lexeme.NewLine:
			if !inAnnotation {
				collectLiteral = false
			}

		case lexeme.MultiLineAnnotationTextBegin:
			inAnnotation = true

		case lexeme.InlineAnnotationTextEnd, lexeme.MultiLineAnnotationTextEnd:
			e.handleEndOfComment(lex, collectLiteral)
			inAnnotation = false
		}
	}
	return nil
}

func (e *Enum) handleLiteralEnd(lex lexeme.LexEvent) error {
	v := lex.Value()
	t, err := schema.GuessSchemaType(v.Data())
	if err != nil {
		return err
	}

	e.values = append(e.values, Value{
		Value: v,
		Type:  t,
	})
	return nil
}

func (e *Enum) handleEndOfComment(lex lexeme.LexEvent, collectLiteral bool) {
	comment := lex.Value().TrimSpaces().String()
	if collectLiteral {
		e.values[len(e.values)-1].Comment = comment
	} else {
		e.values = append(e.values, Value{
			Comment: comment,
			Type:    schema.SchemaTypeComment,
		})
	}
}

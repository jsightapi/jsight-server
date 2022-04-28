package jschema

import (
	"github.com/jsightapi/jsight-schema-go-library/errors"
	"github.com/jsightapi/jsight-schema-go-library/fs"
	"github.com/jsightapi/jsight-schema-go-library/internal/lexeme"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema/internal/scanner"
)

type Enum struct {
	file *fs.File
}

func NewEnum(name string, content []byte) Enum {
	return EnumFromFile(fs.NewFile(name, content))
}

func EnumFromFile(f *fs.File) Enum {
	return Enum{f}
}

func (e Enum) Check() (err error) {
	defer func() {
		err = handlePanic(recover(), nil)
	}()

	scan := scanner.NewSchemaScanner(e.file, false)
	checker := newEnumChecker()

	for {
		lex, ok := scan.Next() // can panic
		if !ok {
			break
		}
		checker.Check(lex) // can panic
	}
	return nil
}

type enumChecker struct {
	// stateFunc a function for running a state machine (the current state of the
	// state machine).
	stateFunc func(lexeme.LexEvent)
}

func newEnumChecker() *enumChecker {
	l := new(enumChecker)
	l.stateFunc = l.begin
	return l
}

// Check checks the lexeme sequence to make sure it is an enum. When any error is
// detected, sends DocumentError into a panic.
func (l *enumChecker) Check(lex lexeme.LexEvent) {
	defer lexeme.CatchLexEventError(lex)
	l.stateFunc(lex)
}

// begin of array "["
func (l *enumChecker) begin(lex lexeme.LexEvent) {
	switch lex.Type() {
	case lexeme.NewLine:
		return
	case lexeme.ArrayBegin:
		l.stateFunc = l.arrayItemBeginOrArrayEnd
	default:
		panic(errors.ErrArrayWasExpectedInEnumRule)
	}
}

// begin of array item begin or array end
// ex: [1 <--
// ex: [" <--
// ex: ] <--
func (l *enumChecker) arrayItemBeginOrArrayEnd(lex lexeme.LexEvent) {
	switch lex.Type() {
	case lexeme.NewLine:
		return
	case lexeme.ArrayItemBegin:
		l.stateFunc = l.literal
	case lexeme.ArrayEnd:
		l.stateFunc = l.afterEndOfEnum
	default:
		panic(errors.ErrLoader)
	}
}

// array item value (literal)
func (l *enumChecker) literal(lex lexeme.LexEvent) {
	switch lex.Type() {
	case lexeme.LiteralBegin:
		return
	case lexeme.LiteralEnd:
		l.stateFunc = l.arrayItemEnd
	default:
		panic(errors.ErrIncorrectArrayItemTypeInEnumRule)
	}
}

// array item end
func (l *enumChecker) arrayItemEnd(lex lexeme.LexEvent) {
	switch lex.Type() {
	case lexeme.NewLine:
		return
	case lexeme.ArrayItemEnd:
		l.stateFunc = l.arrayItemBeginOrArrayEnd
	default:
		panic(errors.ErrLoader)
	}
}

// afterEndOfEnum the method should not be called during normal operation. Ensures
// that the checker will not continue to work after the load is complete.
func (*enumChecker) afterEndOfEnum(lex lexeme.LexEvent) {
	if lex.Type() != lexeme.NewLine {
		panic(errors.ErrUnnecessaryLexemeAfterTheEndOfEnum)
	}
}

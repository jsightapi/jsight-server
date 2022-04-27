package lexeme

import (
	"fmt"

	"github.com/jsightapi/jsight-schema-go-library/errors"
)

func CatchLexEventError(lex LexEvent) {
	r := recover() //nolint:revive // It's okay.
	if r == nil {
		return
	}

	switch val := r.(type) {
	case errors.DocumentError:
		panic(r)
	case errors.Err:
		panic(NewLexEventError(lex, val))
	default:
		panic(NewLexEventError(lex, errors.Format(errors.ErrGeneric, fmt.Sprintf("%s", r))))
	}
}

func CatchLexEventErrorWithIncorrectUserType(lex LexEvent, name string) {
	if name == "" {
		CatchLexEventError(lex)
		return
	}
	r := recover() //nolint:revive // It's okay.
	if r == nil {
		return
	}

	switch val := r.(type) {
	case errors.DocumentError:
		panic(r)
	case errors.Err:
		e := NewLexEventError(lex, val)
		e.SetIncorrectUserType(name)
		panic(e)
	default:
		e := NewLexEventError(lex, errors.Format(errors.ErrGeneric, fmt.Sprintf("%s", r)))
		e.SetIncorrectUserType(name)
		panic(e)
	}
}

func NewLexEventError(lex LexEvent, err errors.Err) errors.DocumentError {
	e := errors.NewDocumentError(lex.File(), err)
	e.SetIndex(lex.Begin())
	return e
}

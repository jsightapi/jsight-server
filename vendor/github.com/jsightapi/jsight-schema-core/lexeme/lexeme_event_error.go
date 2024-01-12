package lexeme

import (
	"fmt"

	"github.com/jsightapi/jsight-schema-core/errs"
	"github.com/jsightapi/jsight-schema-core/kit"
)

func CatchLexEventError(lex LexEvent) {
	r := recover() //nolint:revive // It's okay.
	if r == nil {
		return
	}
	panic(ConvertError(lex, r))
}

func ConvertError(lex LexEvent, err any) kit.JSchemaError {
	switch e := err.(type) {
	case kit.JSchemaError:
		return e
	case errs.Code:
		return NewError(lex, e.F())
	case *errs.Err:
		return NewError(lex, e)
	default:
		return NewError(lex, errs.ErrGeneric.F(fmt.Sprintf("%s", err)))
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
	case kit.JSchemaError:
		panic(r)
	case errs.Code:
		e := NewError(lex, val.F())
		e.SetIncorrectUserType(name)
		panic(e)
	case *errs.Err:
		e := NewError(lex, val)
		e.SetIncorrectUserType(name)
		panic(e)
	default:
		e := NewError(lex, errs.ErrGeneric.F(fmt.Sprintf("%s", r)))
		e.SetIncorrectUserType(name)
		panic(e)
	}
}

func NewError(lex LexEvent, e *errs.Err) kit.JSchemaError {
	ee := kit.NewJSchemaError(lex.File(), e)
	ee.SetIndex(lex.Begin())
	return ee
}

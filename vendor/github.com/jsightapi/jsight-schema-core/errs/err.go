package errs

import (
	"fmt"
	"strings"
)

type Err struct {
	Code_   Code
	message string
}

var (
	_ error = Err{}
	_ error = &Err{}
)

func f(code Code, args ...any) *Err {
	message, ok := errorFormat[code]
	if !ok {
		panic(ErrRuntimeFailure.F())
	}

	cnt := strings.Count(message, "%")
	if cnt != len(args) {
		panic(ErrRuntimeFailure.F())
	}
	if cnt != 0 {
		message = fmt.Sprintf(message, args...)
	}

	return &Err{
		Code_:   code,
		message: message,
	}
}

func (e Err) Code() Code {
	return e.Code_
}

func (e Err) Error() string {
	return e.message
}

func (e Err) Equal(b Err) bool {
	return e.Code_ == b.Code_ && e.message == b.message
}

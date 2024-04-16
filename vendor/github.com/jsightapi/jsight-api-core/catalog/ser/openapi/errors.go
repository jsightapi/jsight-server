package openapi

import (
	"fmt"
	"strings"
)

type Error interface {
	error
	wrapWith(string) Error
	wrapWithf(string, ...any) Error
	text() string
	wrapped() Error
}

type errImpl struct {
	text_    string
	wrapped_ Error
}

func (e errImpl) text() string {
	return e.text_
}

func (e errImpl) wrapped() Error {
	return e.wrapped_
}

func (e errImpl) Error() string {
	var sb strings.Builder
	unwrap(&sb, e)
	return sb.String()
}

func (e errImpl) wrapWith(text string) Error {
	return errImpl{
		text_:    text,
		wrapped_: &e,
	}
}

func (e errImpl) wrapWithf(format string, a ...any) Error {
	text := fmt.Sprintf(format, a...)
	return errImpl{
		text_:    text,
		wrapped_: &e,
	}
}

func unwrap(sb *strings.Builder, e Error) {
	sb.WriteString(e.text())
	if e.wrapped() != nil {
		sb.WriteString(": ")
		unwrap(sb, e.wrapped())
	}
}

func castErr(e error) Error {
	if e == nil {
		return nil
	}
	return e.(errImpl)
}

func newErr(text string) Error {
	return errImpl{
		text_: text,
	}
}

package internal

import (
	jschema "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/errors"
)

type ValidationError struct {
	message string
	code    errors.ErrorCode
}

var _ jschema.ValidationError = ValidationError{}

func NewValidatorError(c errors.ErrorCode, msg string) ValidationError {
	return ValidationError{
		message: msg,
		code:    c,
	}
}

func (v ValidationError) Error() string {
	return errors.Format(v.code).Error()
}

func (v ValidationError) Message() string {
	return v.message
}

func (v ValidationError) ErrCode() int {
	return int(v.code)
}

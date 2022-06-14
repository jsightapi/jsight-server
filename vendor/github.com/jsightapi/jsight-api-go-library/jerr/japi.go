package jerr

import (
	"github.com/jsightapi/jsight-schema-go-library/bytes"
	"github.com/jsightapi/jsight-schema-go-library/fs"
)

type JAPIError struct {
	Msg string
	Location
}

func NewJAPIError(msg string, f *fs.File, i bytes.Index) *JAPIError {
	loc := NewLocation(f, i)
	return &JAPIError{Location: loc, Msg: msg}
}

func (e JAPIError) Error() string {
	return e.Msg
}

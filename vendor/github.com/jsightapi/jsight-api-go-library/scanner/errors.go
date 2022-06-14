package scanner

import (
	"fmt"
	"unicode/utf8"

	"github.com/jsightapi/jsight-schema-go-library/bytes"

	"github.com/jsightapi/jsight-api-go-library/jerr"
)

func (s Scanner) japiError(msg string, i bytes.Index) *jerr.JAPIError {
	return jerr.NewJAPIError(msg, s.file, i)
}

func (s Scanner) japiErrorBasic(msg string) *jerr.JAPIError {
	return jerr.NewJAPIError(msg, s.file, s.curIndex)
}

func (s Scanner) japiErrorUnexpectedChar(where string, expected string) *jerr.JAPIError {
	var msg string
	if s.curIndex < s.dataSize {
		r, _ := utf8.DecodeRune(s.data[s.curIndex:])
		if expected == "" {
			msg = fmt.Sprintf("invalid character %q %s", r, where)
		} else {
			msg = fmt.Sprintf("invalid character %q %s, expecting %s", r, where, expected)
		}
	} else {
		if expected == "" {
			msg = fmt.Sprintf("invalid end of file %s", where)
		} else {
			msg = fmt.Sprintf("invalid end of file %s, expecting %s", where, expected)
		}
	}
	return s.japiError(msg, s.curIndex)
}

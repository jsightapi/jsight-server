package errors

import (
	"fmt"
	"strings"
)

type Errorf struct { //nolint:errname // This is okay.
	args []interface{}
	code ErrorCode
}

func Format(code ErrorCode, args ...interface{}) Errorf {
	return Errorf{
		code: code,
		args: args,
	}
}

func (e Errorf) Code() ErrorCode {
	return e.code
}

func (e Errorf) Error() string {
	if format, ok := errorFormat[e.code]; ok {
		cnt := strings.Count(format, "%s")
		cnt += strings.Count(format, "%q")
		if cnt != len(e.args) {
			panic("Invalid error message: " + format)
		}
		if cnt == 0 {
			return format
		} else {
			return fmt.Sprintf(format, e.args...)
		}
	}
	panic("Unknown error code")
}

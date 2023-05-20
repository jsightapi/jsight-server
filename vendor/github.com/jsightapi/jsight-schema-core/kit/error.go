package kit

import (
	"fmt"

	"github.com/jsightapi/jsight-schema-core/errs"
	"github.com/jsightapi/jsight-schema-core/fs"
)

type Error interface {
	Filename() string
	Index() uint
	Line() uint
	Column() uint
	Message() string
	ErrCode() int
	IncorrectUserType() string
}

// ConvertError converts error to Error interface.
// Used in JSight API Core.
func ConvertError(f *fs.File, err any) Error {
	switch e := err.(type) {
	case JSchemaError:
		return e
	case errs.Code:
		return NewJSchemaError(f, e.F())
	case *errs.Err:
		return NewJSchemaError(f, e)
	case error:
		return NewJSchemaError(f, errs.ErrGeneric.F(e.Error()))
	}
	return NewJSchemaError(f, errs.ErrGeneric.F(fmt.Sprintf("%s", err)))
}

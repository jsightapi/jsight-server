package kit

import (
	"fmt"

	lib "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/errors"
	"github.com/jsightapi/jsight-schema-go-library/fs"
)

type Error interface {
	Filename() string
	Position() uint
	Message() string
	ErrCode() int
	IncorrectUserType() string
}

// ConvertError converts error to Error interface.
// Added for BC
func ConvertError(f *fs.File, err error) Error {
	switch e := err.(type) { //nolint:errorlint // This is okay.
	case errors.ErrorCode:
		return sdkError{
			filename: f.Name(),
			position: 0,
			message:  e.Error(),
			errCode:  int(e.Code()),
		}

	case errors.DocumentError:
		return e

	case lib.ParsingError:
		return sdkError{
			filename: f.Name(),
			position: e.Position(),
			message:  e.Message(),
			errCode:  e.ErrCode(),
		}

	case lib.ValidationError:
		return sdkError{
			filename: f.Name(),
			position: 0,
			message:  e.Message(),
			errCode:  e.ErrCode(),
		}
	}
	return errors.NewDocumentError(f, errors.Format(errors.ErrGeneric, fmt.Sprintf("%s", err)))
}

type sdkError struct {
	filename          string
	message           string
	incorrectUserType string
	position          uint
	errCode           int
}

func (s sdkError) Filename() string          { return s.filename }
func (s sdkError) Position() uint            { return s.position }
func (s sdkError) Message() string           { return s.message }
func (s sdkError) ErrCode() int              { return s.errCode }
func (s sdkError) IncorrectUserType() string { return s.incorrectUserType }

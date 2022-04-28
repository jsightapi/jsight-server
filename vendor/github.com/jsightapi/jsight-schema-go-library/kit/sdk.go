package kit

import (
	"fmt"

	lib "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/errors"
	"github.com/jsightapi/jsight-schema-go-library/formats/json"
	"github.com/jsightapi/jsight-schema-go-library/fs"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema"
)

type Error interface {
	Filename() string
	Position() uint
	Message() string
	ErrCode() int
	IncorrectUserType() string
}

// LengthOfSchema computes length of schema specified in th file.
// Deprecated: Use Len method of jschema.Schema instead.
func LengthOfSchema(f *fs.File) (uint, Error) {
	l, err := jschema.FromFile(
		f,
		jschema.AllowTrailingNonSpaceCharacters(),
	).
		Len()
	if err != nil {
		return 0, convertError(f, err)
	}
	return l, nil
}

// LengthOfJson computes length of JSON document specified in this file.
// Deprecated: Use Len method of jschema.Schema instead.
func LengthOfJson(f *fs.File) (uint, Error) {
	l, err := json.FromFile(f, json.AllowTrailingNonSpaceCharacters()).Len()
	if err != nil {
		return 0, convertError(f, err)
	}
	return l, nil
}

// ValidateJson the key of extraTypes parameter is the name of the type.
// The file name is used only for display in case of an error.
// They may not be the same.
// Deprecated: Use Validate method of jschema.Schema instead.
func ValidateJson(
	schemaFile *fs.File,
	extraTypes map[string]*fs.File,
	jsonFile *fs.File,
	areKeysOptionalByDefault bool,
) Error {
	var oo []jschema.Option
	if areKeysOptionalByDefault {
		oo = append(oo, jschema.KeysAreOptionalByDefault())
	}

	sc := jschema.FromFile(schemaFile, oo...)

	for name, f := range extraTypes {
		if len(f.Content()) == 0 {
			return errors.NewDocumentError(schemaFile, errors.Format(errors.ErrEmptyType, name))
		}
		if err := sc.AddType(name, jschema.FromFile(f, oo...)); err != nil {
			return convertError(f, err)
		}
	}

	err := sc.Validate(json.FromFile(jsonFile))
	if err != nil {
		return convertError(schemaFile, err)
	}
	return nil
}

// convertError converts error to Error interface.
// Added for BC
func convertError(f *fs.File, err error) Error {
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

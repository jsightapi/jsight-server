package kit

import (
	"fmt"

	"github.com/jsightapi/jsight-api-go-library/core"
	"github.com/jsightapi/jsight-api-go-library/jerr"

	"github.com/jsightapi/jsight-schema-go-library/bytes"
	"github.com/jsightapi/jsight-schema-go-library/fs"
	"github.com/jsightapi/jsight-schema-go-library/reader"
)

// JApi is an interface-level wrapper for JApiCore
type JApi struct {
	core *core.JApiCore
}

// NewJapi returns interface-level wrapper for JApiCore
// Does not include .jst file validation. File validation should be called explicitly.
func NewJapi(filepath string, oo ...core.Option) (JApi, error) {
	f, err := readPanicFree(filepath)
	if err != nil {
		return JApi{}, err
	}
	return NewJApiFromFile(f, oo...), nil
}

func readPanicFree(filename string) (f *fs.File, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s", r)
		}
	}()
	f = reader.Read(filename)
	return f, err
}

// Deprecated: use NewJApiFromFile
func NewJapiFromBytes(b bytes.Bytes, oo ...core.Option) JApi {
	return NewJApiFromFile(fs.NewFile("root", b), oo...)
}

func NewJApiFromFile(file *fs.File, oo ...core.Option) JApi {
	return JApi{
		core.NewJApiCore(file, oo...),
	}
}

// ValidateJAPI validates .jst file
func (j *JApi) ValidateJAPI() *jerr.JApiError {
	return j.core.ValidateJAPI()
}

func (j JApi) ToJson() ([]byte, error) {
	return j.core.Catalog().ToJson()
}

func (j JApi) Title() string {
	if j.core != nil && j.core.Catalog() != nil && j.core.Catalog().Info != nil {
		return j.core.Catalog().Info.Title
	}
	return ""
}

func (j JApi) ToJsonIndent() ([]byte, error) {
	return j.core.Catalog().ToJsonIndent()
}

package kit

import (
	"fmt"

	"github.com/jsightapi/jsight-schema-core/fs"
	"github.com/jsightapi/jsight-schema-core/reader"

	"github.com/jsightapi/jsight-api-core/catalog"
	"github.com/jsightapi/jsight-api-core/core"
	"github.com/jsightapi/jsight-api-core/jerr"
)

// JApi is an interface-level wrapper for JApiCore
type JApi struct {
	core *core.JApiCore
}

func NewJapi(filepath string, oo ...core.Option) (JApi, *jerr.JApiError) {
	f, err := readPanicFree(filepath)
	if err != nil {
		return JApi{}, jerr.NewJApiError(err.Error(), f, 0)
	}
	return NewJApiFromFile(f, oo...)
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

func NewJApiFromFile(file *fs.File, oo ...core.Option) (JApi, *jerr.JApiError) {
	j := JApi{
		core.NewJApiCore(file, oo...),
	}
	je := j.core.BuildCatalog()
	if je != nil {
		return j, je
	}
	return j, nil
}

func (j *JApi) Catalog() *catalog.Catalog {
	return j.core.Catalog()
}

func (j *JApi) Title() string {
	if j.core != nil && j.Catalog() != nil && j.Catalog().Info != nil {
		return j.Catalog().Info.Title
	}
	return ""
}

func (j *JApi) ToJson() ([]byte, error) {
	return j.Catalog().ToJson()
}

func (j *JApi) ToJsonIndent() ([]byte, error) {
	return j.Catalog().ToJsonIndent()
}

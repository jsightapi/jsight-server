package reader

import (
	"os"

	"github.com/jsightapi/jsight-schema-go-library/errors"
	"github.com/jsightapi/jsight-schema-go-library/fs"
)

// Read reads the contents of the file, returns a slice of bytes.
func Read(filename string) *fs.File {
	return ReadWithName(filename, filename)
}

func ReadWithName(filename, name string) *fs.File {
	data, err := os.ReadFile(filename)
	if err != nil {
		docErr := errors.DocumentError{}
		docErr.SetMessage(err.Error())
		panic(docErr)
	}
	return fs.NewFile(name, data)
}

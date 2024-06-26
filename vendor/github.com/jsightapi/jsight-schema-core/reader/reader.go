package reader

import (
	"fmt"
	"os"

	"github.com/jsightapi/jsight-schema-core/fs"
	"github.com/jsightapi/jsight-schema-core/kit"
)

// Read reads the contents of the file, returns a slice of bytes.
func Read(filename string) *fs.File {
	return ReadWithName(filename, filename)
}

func ReadWithName(filename, name string) *fs.File {
	data, err := os.ReadFile(filename)
	if err != nil {
		docErr := kit.JSchemaError{}
		docErr.SetMessage(fmt.Sprintf(`The system cannot open the file %q.`, name))
		panic(docErr)
	}
	return fs.NewFile(name, data)
}

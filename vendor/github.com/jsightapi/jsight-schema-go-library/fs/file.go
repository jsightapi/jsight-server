package fs

import (
	"fmt"

	"github.com/jsightapi/jsight-schema-go-library/bytes"
)

// File represent a file.
type File struct {
	name    string
	content bytes.Bytes
}

// NewFile creates new File instance.
func NewFile[T FileContent](name string, content T) *File {
	return &File{
		name:    name,
		content: normalizeFileContent(content),
	}
}

// FileContent all allowed types for specifying File's content.
type FileContent interface {
	string | bytes.Bytes | []byte
}

// normalizeFileContent convert generic FileContent to bytes.Bytes 'cause we operate
// this type in the file.
func normalizeFileContent[T FileContent](content T) bytes.Bytes {
	switch c := any(content).(type) {
	case string:
		return bytes.Bytes(c)
	case []byte:
		return c
	case bytes.Bytes:
		return c
	}

	// This might happen only when we extend `FileContent` interface and forget
	// to add new case to the type switch above this point.
	panic(fmt.Sprintf("Unhandled content type %T", content))
}

// Name returns file name.
func (f File) Name() string {
	return f.name
}

// Content returns file content.
func (f File) Content() bytes.Bytes {
	return f.content
}

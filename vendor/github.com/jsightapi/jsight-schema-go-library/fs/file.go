package fs

import "github.com/jsightapi/jsight-schema-go-library/bytes"

// File represent a file.
type File struct {
	name    string
	content bytes.Bytes
}

// NewFile creates new File instance.
func NewFile(name string, content bytes.Bytes) *File {
	return &File{
		name:    name,
		content: content,
	}
}

// Name returns file name.
func (f File) Name() string {
	return f.name
}

// Content returns file content.
func (f File) Content() bytes.Bytes {
	return f.content
}

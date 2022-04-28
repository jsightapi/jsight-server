package fs

import "github.com/jsightapi/jsight-schema-go-library/bytes"

type File struct {
	name    string
	content bytes.Bytes
}

func NewFile(name string, content bytes.Bytes) *File {
	return &File{
		name:    name,
		content: content,
	}
}

func (f File) Name() string {
	return f.name
}

func (f *File) SetName(filename string) {
	f.name = filename
}

func (f File) Content() bytes.Bytes {
	return f.content
}

func (f *File) SetContent(content bytes.Bytes) {
	f.content = content
}

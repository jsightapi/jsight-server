package jerr

import (
	"github.com/jsightapi/jsight-schema-core/bytes"
	"github.com/jsightapi/jsight-schema-core/fs"
)

type Location struct {
	File   *fs.File
	Quote  string
	Index  bytes.Index
	Line   bytes.Index
	Column bytes.Index
}

// NewLocation a bit optimized version of getting all info
func NewLocation(f *fs.File, i bytes.Index) Location {
	loc := Location{
		File:  f,
		Index: i,
		Quote: quote(f.Content(), i),
	}
	loc.Line, loc.Column = f.Content().LineAndColumn(i)
	return loc
}

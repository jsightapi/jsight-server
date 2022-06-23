package jerr

import (
	"github.com/jsightapi/jsight-schema-go-library/bytes"
	"github.com/jsightapi/jsight-schema-go-library/fs"
)

type Location struct {
	file  *fs.File
	quote string
	index bytes.Index
	line  bytes.Index
}

// NewLocation a bit optimized version of getting all info
func NewLocation(f *fs.File, i bytes.Index) Location {
	nl := DetectNewLineSymbol(f.Content())
	lb := LineBeginning(f.Content(), i, nl)

	return Location{
		file:  f,
		index: i,
		quote: quote(f.Content(), i, lb, nl),
		line:  LineNumber(f.Content(), i, nl),
		//positionInLine: p - lb
	}
}

func (l Location) Quote() string {
	return l.quote
}

func (l Location) Line() bytes.Index {
	return l.line
}

func (l Location) Index() bytes.Index {
	return l.index
}

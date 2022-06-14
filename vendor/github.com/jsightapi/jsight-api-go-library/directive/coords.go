package directive

import (
	"fmt"

	"github.com/jsightapi/jsight-schema-go-library/bytes"
	"github.com/jsightapi/jsight-schema-go-library/fs"
)

type Coords struct {
	f *fs.File
	b bytes.Index
	e bytes.Index
}

func NewCoords(f *fs.File, b bytes.Index, e bytes.Index) Coords {
	return Coords{f, b, e}
}

func (c Coords) Read() bytes.Bytes {
	return c.f.Content().Slice(c.b, c.e)
}

func (c Coords) IsSet() bool {
	return c.f != nil && c.e != 0
}

func (c Coords) B() bytes.Index {
	return c.b
}

func (c Coords) String() string {
	return fmt.Sprintf("[%d:%d]", c.b, c.e)
}

func (c Coords) File() *fs.File {
	return c.f
}

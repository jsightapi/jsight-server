package directive

import (
	"fmt"

	"github.com/jsightapi/jsight-schema-core/bytes"
	"github.com/jsightapi/jsight-schema-core/fs"
)

type Coords struct {
	file  *fs.File
	begin bytes.Index
	end   bytes.Index
}

func NewCoords(f *fs.File, begin, end bytes.Index) Coords {
	return Coords{f, begin, end}
}

func (c Coords) Read() bytes.Bytes {
	return c.file.Content().Sub(c.begin, c.end+1)
}

func (c Coords) IsSet() bool {
	return c.file != nil && c.end != 0
}

func (c Coords) Begin() bytes.Index {
	return c.begin
}

func (c Coords) String() string {
	return fmt.Sprintf("[%d:%d]", c.begin, c.end)
}

func (c Coords) File() *fs.File {
	return c.file
}

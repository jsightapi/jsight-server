package ischema

import (
	"github.com/jsightapi/jsight-schema-core/bytes"
	"github.com/jsightapi/jsight-schema-core/fs"
)

type Type struct {
	Schema   *ISchema
	RootFile *fs.File
	Begin    bytes.Index
}

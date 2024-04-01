package openapi

import (
	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/errs"
	"github.com/jsightapi/jsight-schema-core/notations/jschema"
	"github.com/jsightapi/jsight-schema-core/notations/regex"
	"github.com/jsightapi/jsight-schema-core/openapi/internal/jsoac"
	"github.com/jsightapi/jsight-schema-core/openapi/internal/rsoac"
)

type SchemaObject interface {
	SetDescription(s string)
	MarshalJSON() (b []byte, err error)
}

var _ SchemaObject = (*jsoac.JSOAC)(nil)
var _ SchemaObject = (*rsoac.RSOAC)(nil)

func NewSchemaObject(s schema.Schema) SchemaObject {
	switch st := any(s).(type) {
	case *jschema.JSchema:
		return jsoac.New(st)
	case *regex.RSchema:
		return rsoac.New(st)
	}

	panic(errs.ErrRuntimeFailure.F())
}

package jsoac

import (
	"encoding/json"

	"github.com/jsightapi/jsight-schema-core/openapi/internal"

	schema "github.com/jsightapi/jsight-schema-core"
)

type Nullable struct {
	value []byte
}

var _ json.Marshaler = Nullable{}
var _ json.Marshaler = &Nullable{}

func newNullable(astNode schema.ASTNode) *Nullable {
	if astNode.Rules.Has("nullable") && astNode.Rules.GetValue("nullable").Value == internal.StringTrue {
		return newNullableFromBool(true)
	}
	return nil
}

func newNullableFromBool(b bool) *Nullable {
	if b {
		return &Nullable{[]byte(`true`)}
	}
	return nil
}

func (n Nullable) MarshalJSON() (b []byte, err error) {
	return n.value, nil
}

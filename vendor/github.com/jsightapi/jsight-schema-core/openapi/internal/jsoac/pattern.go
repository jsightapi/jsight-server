package jsoac

import (
	"encoding/json"

	"github.com/jsightapi/jsight-schema-core/openapi/internal"

	schema "github.com/jsightapi/jsight-schema-core"
)

type Pattern struct {
	value []byte
}

var _ json.Marshaler = Pattern{}
var _ json.Marshaler = &Pattern{}

func newPattern(astNode schema.ASTNode) *Pattern {
	if astNode.Rules.Has("regex") {
		return &Pattern{value: internal.ToJSONString(astNode.Rules.GetValue("regex").Value)}
	}
	return nil
}

func (ex Pattern) MarshalJSON() (b []byte, err error) {
	return ex.value, nil
}

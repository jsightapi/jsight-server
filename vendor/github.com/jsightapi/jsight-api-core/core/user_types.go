package core

import (
	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/notations/jschema"
)

func (core *JApiCore) UserTypesData() map[string]*jschema.JSchema {
	ut := make(map[string]*jschema.JSchema, core.userTypes.Len())
	_ = core.userTypes.Each(func(k string, s schema.Schema) error {
		if ss, ok := s.(*jschema.JSchema); ok {
			ut[k] = ss
		}
		return nil
	})
	return ut
}

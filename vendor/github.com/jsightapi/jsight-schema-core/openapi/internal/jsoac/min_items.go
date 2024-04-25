package jsoac

import (
	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/openapi/internal"
)

func newMinItems(astNode schema.ASTNode) *int64 {
	if astNode.Rules.Has("minItems") {
		return internal.Int64RefByString(astNode.Rules.GetValue("minItems").Value)
	}
	return nil
}

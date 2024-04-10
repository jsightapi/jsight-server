package jsoac

import (
	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/openapi/internal"
)

func newMaxItems(astNode schema.ASTNode) *int64 {
	if astNode.Rules.Has("maxItems") {
		return internal.Int64RefByString(astNode.Rules.GetValue("maxItems").Value)
	}
	return nil
}

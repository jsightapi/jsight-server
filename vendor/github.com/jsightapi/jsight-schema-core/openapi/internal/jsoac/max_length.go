package jsoac

import (
	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/openapi/internal"
)

func newMaxLength(astNode schema.ASTNode) *int64 {
	if astNode.Rules.Has("maxLength") {
		return internal.Int64RefByString(astNode.Rules.GetValue("maxLength").Value)
	}
	return nil
}

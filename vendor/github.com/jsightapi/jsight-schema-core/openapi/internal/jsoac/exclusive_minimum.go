package jsoac

import (
	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/openapi/internal"
)

func newExclusiveMinimum(astNode schema.ASTNode) *bool {
	if astNode.Rules.Has("exclusiveMinimum") {
		b := astNode.Rules.GetValue("exclusiveMinimum").Value == internal.StringTrue
		if b {
			return &b
		}
		return nil
	}
	return nil
}

package jsoac

import (
	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/openapi/internal"
)

func newExclusiveMaximum(astNode schema.ASTNode) *bool {
	if astNode.Rules.Has("exclusiveMaximum") {
		b := astNode.Rules.GetValue("exclusiveMaximum").Value == internal.StringTrue
		if b {
			return &b
		}
		return nil
	}
	return nil
}

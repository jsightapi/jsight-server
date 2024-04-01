package jsoac

import (
	schema "github.com/jsightapi/jsight-schema-core"
)

func newExclusiveMinimum(astNode schema.ASTNode) *bool {
	if astNode.Rules.Has("exclusiveMinimum") {
		b := astNode.Rules.GetValue("exclusiveMinimum").Value == stringTrue
		if b {
			return &b
		}
		return nil
	}
	return nil
}

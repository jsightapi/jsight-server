package jsoac

import (
	schema "github.com/jsightapi/jsight-schema-core"
)

func newExclusiveMaximum(astNode schema.ASTNode) *bool {
	if astNode.Rules.Has("exclusiveMaximum") {
		b := astNode.Rules.GetValue("exclusiveMaximum").Value == stringTrue
		if b {
			return &b
		}
		return nil
	}
	return nil
}

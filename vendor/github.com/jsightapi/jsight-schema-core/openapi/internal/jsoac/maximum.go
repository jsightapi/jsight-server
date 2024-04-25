package jsoac

import (
	schema "github.com/jsightapi/jsight-schema-core"
)

func newMaximum(astNode schema.ASTNode) *Number {
	if astNode.Rules.Has("max") {
		return newNumber(astNode.Rules.GetValue("max").Value)
	}
	return nil
}

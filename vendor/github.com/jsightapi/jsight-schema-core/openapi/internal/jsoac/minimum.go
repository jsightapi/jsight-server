package jsoac

import (
	schema "github.com/jsightapi/jsight-schema-core"
)

func newMinimum(astNode schema.ASTNode) *Number {
	if astNode.Rules.Has("min") {
		return newNumber(astNode.Rules.GetValue("min").Value)
	}
	return nil
}

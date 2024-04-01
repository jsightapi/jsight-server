package jsoac

import (
	schema "github.com/jsightapi/jsight-schema-core"
)

func newMinItems(astNode schema.ASTNode) *int64 {
	if astNode.Rules.Has("minItems") {
		return int64RefByString(astNode.Rules.GetValue("minItems").Value)
	}
	return nil
}

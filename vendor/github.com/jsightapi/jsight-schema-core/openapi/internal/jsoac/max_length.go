package jsoac

import (
	schema "github.com/jsightapi/jsight-schema-core"
)

func newMaxLength(astNode schema.ASTNode) *int64 {
	if astNode.Rules.Has("maxLength") {
		return int64RefByString(astNode.Rules.GetValue("maxLength").Value)
	}
	return nil
}

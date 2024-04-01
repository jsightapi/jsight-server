package jsoac

import (
	schema "github.com/jsightapi/jsight-schema-core"
)

func newMaxItems(astNode schema.ASTNode) *int64 {
	if astNode.Rules.Has("maxItems") {
		return int64RefByString(astNode.Rules.GetValue("maxItems").Value)
	}
	return nil
}

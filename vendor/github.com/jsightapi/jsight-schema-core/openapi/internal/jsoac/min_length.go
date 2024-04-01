package jsoac

import (
	schema "github.com/jsightapi/jsight-schema-core"
)

func newMinLength(astNode schema.ASTNode) *int64 {
	if astNode.Rules.Has("minLength") {
		return int64RefByString(astNode.Rules.GetValue("minLength").Value)
	}
	return nil
}

package jsoac

import (
	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/openapi/internal"
)

func newMinLength(astNode schema.ASTNode) *int64 {
	if astNode.Rules.Has("minLength") {
		return internal.Int64RefByString(astNode.Rules.GetValue("minLength").Value)
	}
	return nil
}

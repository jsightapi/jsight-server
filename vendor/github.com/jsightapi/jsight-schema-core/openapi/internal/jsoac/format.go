package jsoac

import (
	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/openapi/internal"
)

func newFormat(astNode schema.ASTNode) *string {
	if astNode.Rules.Has("type") {
		return internal.FormatFromSchemaType(astNode.Rules.GetValue("type").Value)
	}
	return nil
}

package jsoac

import (
	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/openapi/internal"
)

func newConst(astNode schema.ASTNode) *Enum {
	if astNode.Rules.Has("const") && astNode.Rules.GetValue("const").Value == internal.StringTrue {
		ex := newExample(astNode.Value, internal.IsString(astNode))
		bb := ex.jsonValue()

		enum := makeEmptyEnum()
		enum.append(bb)
		return enum
	}
	return nil
}

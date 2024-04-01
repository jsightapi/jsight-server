package jsoac

import (
	schema "github.com/jsightapi/jsight-schema-core"
)

func newConst(astNode schema.ASTNode) *Enum {
	if astNode.Rules.Has("const") && astNode.Rules.GetValue("const").Value == stringTrue {
		ex := newExample(astNode.Value, isString(astNode))
		bb := ex.jsonValue()

		enum := makeEmptyEnum()
		enum.append(bb)
		return enum
	}
	return nil
}

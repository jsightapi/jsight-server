package jsoac

import (
	"encoding/json"

	schema "github.com/jsightapi/jsight-schema-core"
)

type Enum struct {
	list [][]byte
}

var _ json.Marshaler = Enum{}
var _ json.Marshaler = &Enum{}

func newEnum(astNode schema.ASTNode) *Enum {
	if enum := newConst(astNode); enum != nil {
		return enum
	}

	if astNode.Rules.Has(stringEnum) &&
		astNode.Rules.GetValue(stringEnum).TokenType == stringArray &&
		0 < len(astNode.Rules.GetValue(stringEnum).Items) {
		enum := makeEmptyEnum()
		for _, s := range astNode.Rules.GetValue(stringEnum).Items {
			ex := newExample(s.Value, s.TokenType == "string")
			bb := ex.jsonValue()

			enum.append(bb)
		}
		return enum
	}

	if astNode.SchemaType == stringNull {
		enum := makeEmptyEnum()
		ex := newExample(stringNull, false)
		enum.append(ex.jsonValue())
		return enum
	}

	return nil
}

func makeEmptyEnum() *Enum {
	return &Enum{
		list: make([][]byte, 0, 3),
	}
}

func (e *Enum) append(b []byte) {
	e.list = append(e.list, b)
}

func (e Enum) MarshalJSON() ([]byte, error) {
	b := bufferPool.Get()
	defer bufferPool.Put(b)

	b.WriteByte('[')
	for i, item := range e.list {
		b.Write(item)

		if i+1 != len(e.list) {
			b.WriteByte(',')
		}
	}
	b.WriteByte(']')

	return b.Bytes(), nil
}

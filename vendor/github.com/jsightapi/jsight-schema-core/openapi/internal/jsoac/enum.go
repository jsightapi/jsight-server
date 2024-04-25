package jsoac

import (
	"encoding/json"

	"github.com/jsightapi/jsight-schema-core/openapi/internal"

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

	if astNode.Rules.Has(internal.StringEnum) &&
		astNode.Rules.GetValue(internal.StringEnum).TokenType == internal.StringArray &&
		0 < len(astNode.Rules.GetValue(internal.StringEnum).Items) {
		enum := makeEmptyEnum()
		for _, s := range astNode.Rules.GetValue(internal.StringEnum).Items {
			ex := newExample(s.Value, s.TokenType == "string")
			bb := ex.jsonValue()

			enum.append(bb)
		}
		return enum
	}

	if astNode.SchemaType == internal.StringNull {
		enum := makeEmptyEnum()
		ex := newExample(internal.StringNull, false)
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
	b := internal.BufferPool.Get()
	defer internal.BufferPool.Put(b)

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

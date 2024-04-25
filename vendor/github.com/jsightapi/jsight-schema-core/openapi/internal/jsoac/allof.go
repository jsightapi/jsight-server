package jsoac

import (
	"encoding/json"

	"github.com/jsightapi/jsight-schema-core/openapi/internal"

	schema "github.com/jsightapi/jsight-schema-core"

	"github.com/jsightapi/jsight-schema-core/errs"
)

type AllOf struct {
	userTypeNames []string
}

var _ json.Marshaler = AllOf{}
var _ json.Marshaler = &AllOf{}

func newAllOf(astNode schema.ASTNode) *AllOf {
	if astNode.Rules.Has("allOf") {
		a := &AllOf{
			userTypeNames: make([]string, 0, 5),
		}
		a.append(astNode.Rules.GetValue("allOf"))
		return a
	}
	return nil
}

func (a *AllOf) append(rule schema.RuleASTNode) {
	switch rule.TokenType {
	case schema.TokenTypeShortcut:
		a.userTypeNames = append(a.userTypeNames, rule.Value)
	case schema.TokenTypeArray:
		for _, child := range rule.Items {
			a.append(child)
		}
	default:
		panic(errs.ErrRuntimeFailure.F())
	}
}

func (a AllOf) MarshalJSON() ([]byte, error) {
	b := internal.BufferPool.Get()
	defer internal.BufferPool.Put(b)

	b.WriteByte('[')

	for i, name := range a.userTypeNames {
		ref := newRefFromUserTypeName(name, false)
		rb, err := ref.MarshalJSON()
		if err != nil {
			return nil, err
		}

		b.Write(rb)

		if i+1 != len(a.userTypeNames) {
			b.WriteByte(',')
		}
	}

	b.WriteByte(']')

	return b.Bytes(), nil
}

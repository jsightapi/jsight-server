package jsoac

// JSight schema to OpenAPi converter

import (
	"encoding/json"

	schema "github.com/jsightapi/jsight-schema-core"

	"github.com/jsightapi/jsight-schema-core/notations/jschema"
)

type JSOAC struct {
	root        Node
	description *string
}

func New(j *jschema.JSchema) *JSOAC {
	return NewFromASTNode(j.ASTNode)
}

func NewFromASTNode(astNode schema.ASTNode) *JSOAC {
	return &JSOAC{
		root: newNode(astNode),
	}
}

func (o *JSOAC) SetDescription(s string) {
	o.description = &s
}

func (o JSOAC) MarshalJSON() (b []byte, err error) {
	if o.description != nil {
		o.root.SetNodeDescription(*o.description)
	}
	return json.Marshal(o.root)
}

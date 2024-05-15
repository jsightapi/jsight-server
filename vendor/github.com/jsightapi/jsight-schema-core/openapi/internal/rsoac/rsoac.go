package rsoac

import (
	"encoding/json"

	schema "github.com/jsightapi/jsight-schema-core"

	"github.com/jsightapi/jsight-schema-core/notations/regex"
)

// RSOAC Regex schema to OpenAPi converter
type RSOAC struct {
	root        *RegexString
	description *string
}

func New(rs *regex.RSchema) *RSOAC {
	astNode := getASTNode(rs)
	return NewFromASTNode(astNode)
}

func getASTNode(rs *regex.RSchema) schema.ASTNode {
	an, err := rs.GetAST()
	if err != nil {
		panic(err)
	}
	return an
}

func NewFromASTNode(astNode schema.ASTNode) *RSOAC {
	return &RSOAC{
		root: newRegexString(astNode),
	}
}

// SetDescription has no effect on the resulting OpenAPI
func (o *RSOAC) SetDescription(s string) {
	o.description = &s
}

func (o RSOAC) MarshalJSON() (b []byte, err error) {
	return json.Marshal(o.root)
}

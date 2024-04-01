package openapi

import (
	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/errs"
	"github.com/jsightapi/jsight-schema-core/notations/jschema"
	"github.com/jsightapi/jsight-schema-core/notations/regex"
	"github.com/jsightapi/jsight-schema-core/openapi/internal/jsoac"
	"github.com/jsightapi/jsight-schema-core/openapi/internal/rsoac"
)

type SchemaInfo struct {
	regex *regex.RSchema
	node  *schema.ASTNode
}

var _ SchemaInformer = SchemaInfo{}
var _ SchemaInformer = (*SchemaInfo)(nil)

func NewRSchemaInfo(rs *regex.RSchema) SchemaInfo {
	return SchemaInfo{
		regex: rs,
	}
}

func NewJSchemaInfo(js *jschema.JSchema) SchemaInfo {
	return SchemaInfo{
		node: &(js.ASTNode),
	}
}

func newJSchemaInfoFromASTNode(astNode schema.ASTNode) SchemaInfo {
	return SchemaInfo{
		node: &astNode,
	}
}

func (e SchemaInfo) Type() SchemaInfoType {
	if e.regex != nil {
		return SchemaInfoTypeRegex
	}

	if e.node.SchemaType == "any" {
		return SchemaInfoTypeAny
	}

	switch e.node.TokenType {
	case schema.TokenTypeNumber, schema.TokenTypeString, schema.TokenTypeBoolean, schema.TokenTypeNull:
		return SchemaInfoTypeScalar
	case schema.TokenTypeObject:
		return SchemaInfoTypeObject
	case schema.TokenTypeArray:
		return SchemaInfoTypeArray
	case schema.TokenTypeShortcut:
		return SchemaInfoTypeReference
	default:
		panic(errs.ErrRuntimeFailure.F())
	}
}

func (e SchemaInfo) SchemaObject() SchemaObject {
	if e.regex != nil {
		return rsoac.New(e.regex)
	}

	return jsoac.NewFromASTNode(*e.node)
}

func (e SchemaInfo) Annotation() string {
	if e.node != nil {
		return e.node.Comment
	}
	return ""
}

func (e SchemaInfo) Children() []schema.ASTNode {
	return e.node.Children
}

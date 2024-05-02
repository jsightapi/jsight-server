package openapi

import (
	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/errs"
	"github.com/jsightapi/jsight-schema-core/notations/jschema"
	"github.com/jsightapi/jsight-schema-core/notations/regex"
	"github.com/jsightapi/jsight-schema-core/openapi/internal"
)

var empty struct{}

type dereference struct {
	userTypes    map[string]schema.Schema
	visitedTypes map[string]struct{}
	result       *schemaInfoList
}

func Dereference(s schema.Schema) []SchemaInformer {
	var d dereference

	if st, ok := s.(*jschema.JSchema); ok {
		d = newDereference(st.UserTypeCollection)
	} else {
		d = newDereference(nil)
	}

	d.schema(s)

	return d.result.list()
}

func newDereference(userTypes map[string]schema.Schema) dereference {
	return dereference{
		userTypes:    userTypes,
		visitedTypes: make(map[string]struct{}, len(userTypes)),
		result:       newSchemaInfoList(),
	}
}

func (d dereference) schema(s schema.Schema) {
	switch st := s.(type) {
	case *jschema.JSchema:
		d.jSchema(st.ASTNode)
	case *regex.RSchema:
		d.rSchema(st)
	default:
		panic(errs.ErrRuntimeFailure.F())
	}
}

func (d dereference) rSchema(rs *regex.RSchema) {
	info := NewRSchemaInfo(rs)
	d.result.append(info)
}

func (d *dereference) jSchema(astNode schema.ASTNode) {
	if rule, ok := astNode.Rules.Get("or"); ok {
		for _, item := range rule.Items {
			d.orItem(item)
		}
		return
	}

	switch astNode.TokenType {
	case schema.TokenTypeNumber, schema.TokenTypeString, schema.TokenTypeBoolean, schema.TokenTypeNull,
		schema.TokenTypeArray:
		info := newJSchemaInfoFromASTNode(astNode)
		d.result.append(info)
	case schema.TokenTypeObject:
		info := newObjectInfo(astNode, d.userTypes)
		d.result.append(info)
	case schema.TokenTypeShortcut:
		name := astNode.Value
		if _, ok := d.visitedTypes[name]; !ok {
			d.visitedTypes[name] = empty
			d.userType(name)
		}
	default:
		panic(errs.ErrRuntimeFailure.F())
	}
}

func (d dereference) userType(name string) {
	ut, ok := d.userTypes[name]
	if !ok {
		panic(errs.ErrUserTypeNotFound.F(name))
	}

	d.schema(ut)
}

func (d dereference) orItem(r schema.RuleASTNode) {
	mockAstNode := internal.RuleToASTNode(r)
	d.jSchema(mockAstNode)
}

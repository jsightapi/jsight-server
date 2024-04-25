package openapi

import (
	schema "github.com/jsightapi/jsight-schema-core"
)

type PropertyInfo struct {
	SchemaInfo
}

var _ PropertyInformer = PropertyInfo{}
var _ PropertyInformer = (*PropertyInfo)(nil)

func newPropertyInfo(astNode schema.ASTNode) PropertyInfo {
	return PropertyInfo{
		SchemaInfo: newJSchemaInfoFromASTNode(astNode),
	}
}

func (i PropertyInfo) Key() string {
	return i.node.Key
}

func (i PropertyInfo) Optional() bool {
	v, ok := i.node.Rules.Get("optional")
	if !ok {
		return false
	}
	return v.Value == "true"
}

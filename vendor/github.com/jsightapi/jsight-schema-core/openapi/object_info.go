package openapi

import schema "github.com/jsightapi/jsight-schema-core"

type ObjectInfo struct {
	SchemaInfo
}

var _ ObjectInformer = ObjectInfo{}
var _ ObjectInformer = (*ObjectInfo)(nil)

func newObjectInfo(astNode schema.ASTNode) ObjectInfo {
	return ObjectInfo{
		SchemaInfo: newJSchemaInfoFromASTNode(astNode),
	}
}

func (o ObjectInfo) PropertiesInfos() []PropertyInformer {
	props := o.SchemaInfo.Children()
	result := make([]PropertyInformer, len(props))

	for i, child := range props {
		result[i] = newPropertyInfo(child)
	}

	return result
}

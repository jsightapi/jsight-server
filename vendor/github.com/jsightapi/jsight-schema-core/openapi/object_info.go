package openapi

import (
	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/errs"
)

type ObjectInfo struct {
	SchemaInfo
	userTypes map[string]schema.Schema
}

var _ ObjectInformer = ObjectInfo{}
var _ ObjectInformer = (*ObjectInfo)(nil)

func newObjectInfo(astNode schema.ASTNode, userTypes map[string]schema.Schema) ObjectInfo {
	return ObjectInfo{
		SchemaInfo: newJSchemaInfoFromASTNode(astNode),
		userTypes:  userTypes,
	}
}

func (o ObjectInfo) PropertiesInfos() []PropertyInformer {
	props := o.SchemaInfo.Children()
	result := make([]PropertyInformer, len(props))

	for i, child := range props {
		result[i] = newPropertyInfo(child)
	}

	if rule, ok := o.node.Rules.Get("allOf"); ok {
		result = append(result, o.allOf(rule)...)
	}

	return result
}

func (o ObjectInfo) allOf(ruleAllOf schema.RuleASTNode) []PropertyInformer {
	result := make([]PropertyInformer, 0, 5)

	switch ruleAllOf.TokenType {
	case schema.TokenTypeShortcut:
		result = append(result, o.dereferenceUserTypeProperties(ruleAllOf.Value)...)
	case schema.TokenTypeArray:
		for _, item := range ruleAllOf.Items {
			result = append(result, o.dereferenceUserTypeProperties(item.Value)...)
		}
	default:
		panic(errs.ErrRuntimeFailure.F())
	}

	return result
}

func (o ObjectInfo) dereferenceUserTypeProperties(userTypeName string) []PropertyInformer {
	result := make([]PropertyInformer, 0, 5)

	d := newDereference(o.userTypes)
	d.userType(userTypeName)
	si := d.result.list()

	for _, child := range si {
		if oi, ok := child.(ObjectInformer); ok {
			result = append(result, oi.PropertiesInfos()...)
		} else {
			panic(errs.ErrRuntimeFailure.F())
		}
	}

	return result
}

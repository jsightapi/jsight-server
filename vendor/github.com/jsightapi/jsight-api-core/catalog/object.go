package catalog

import (
	"github.com/jsightapi/jsight-schema-core/json"
	"github.com/jsightapi/jsight-schema-core/notations/jschema"
	"github.com/jsightapi/jsight-schema-core/notations/jschema/ischema"
)

type JSchemaObject struct {
	*jschema.JSchema
}

func (s *JSchemaObject) RootObjectNode() (*ischema.ObjectNode, bool) {
	root := s.Inner.RootNode()
	if root.Type() != json.TypeObject {
		return nil, false
	}

	obj, ok := root.(*ischema.ObjectNode)
	if !ok {
		return nil, false
	}

	return obj, true
}

func (s *JSchemaObject) ObjectProperty(key string) (ischema.Node, bool) {
	obj, ok := s.RootObjectNode()
	if !ok {
		return nil, false
	}
	return obj.Child(key, false)
}

func (s *JSchemaObject) ObjectFirstLevelProperties(ut map[string]*jschema.JSchema) map[string]ischema.Node {
	m := make(map[string]ischema.Node, 5)
	s.objectFirstLevelProperties(m, ut)
	return m
}

func (s *JSchemaObject) objectFirstLevelProperties(m map[string]ischema.Node, ut map[string]*jschema.JSchema) {
	s.appendPropertiesFromShortcut(m, ut)
	s.appendPropertiesFromObject(m)
}

func (s *JSchemaObject) appendPropertiesFromShortcut(m map[string]ischema.Node, ut map[string]*jschema.JSchema) {
	names := jschema.UserTypeNamesFromEachTypeConstraint(s.Inner.RootNode())

	for _, name := range names {
		if ss, ok := ut[name]; ok {
			obj := JSchemaObject{JSchema: ss}
			obj.objectFirstLevelProperties(m, ut)
		}
	}
}

func (s *JSchemaObject) appendPropertiesFromObject(m map[string]ischema.Node) {
	obj, ok := s.RootObjectNode()
	if !ok {
		return
	}

	for _, v := range obj.Keys().Data {
		if !v.IsShortcut {
			if n, ok := obj.Child(v.Key, false); ok {
				n.SetParent(nil)
				m[v.Key] = n
			}
		}
	}
}

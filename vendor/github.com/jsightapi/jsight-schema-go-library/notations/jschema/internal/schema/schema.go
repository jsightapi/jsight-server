package schema

import (
	"fmt"
	"strings"

	"github.com/jsightapi/jsight-schema-go-library/bytes"
	"github.com/jsightapi/jsight-schema-go-library/errors"
	"github.com/jsightapi/jsight-schema-go-library/fs"
)

type Schema struct {
	// types the map where key is the name of the type (or included Schema).
	types    map[string]Type
	rootNode Node
}

func New() Schema {
	return Schema{
		types: make(map[string]Type, 5),
	}
}

func (s Schema) TypesList() map[string]Type {
	return s.types
}

// Type returns *Schema or panic if not found.
func (s Schema) Type(name string) *Schema { // todo confuses that here is returning Schema, instead Type. So historically. Perhaps it is worth remaking?
	t, ok := s.types[name]
	if ok {
		return t.schema
	}
	panic(errors.Format(errors.ErrTypeNotFound, name))
}

func (s Schema) RootNode() Node {
	return s.rootNode
}

func (s *Schema) AddNamedType(name string, typ *Schema, rootFile *fs.File, begin bytes.Index) {
	if !bytes.Bytes(name).IsUserTypeName() {
		panic(errors.Format(errors.ErrInvalidSchemaName, name))
	}
	s.addType(name, typ, rootFile, begin)
}

// AddUnnamedType Adds an unnamed TYPE to the SCHEMA. Returns a unique name for the added TYPE.
func (s *Schema) AddUnnamedType(typ *Schema, rootFile *fs.File, begin bytes.Index) string {
	name := fmt.Sprintf("#%p", typ)
	s.addType(name, typ, rootFile, begin)
	return name
}

func (s *Schema) addType(name string, schema *Schema, rootFile *fs.File, begin bytes.Index) {
	if _, ok := s.types[name]; ok {
		panic(errors.Format(errors.ErrDuplicationOfNameOfTypes, name))
	}
	s.types[name] = Type{schema, rootFile, begin}
}

func (s *Schema) AddType(n string, t Type) {
	s.types[n] = t
}

func (s *Schema) SetRootNode(node Node) {
	s.rootNode = node
}

func (s Schema) String() string {
	var str strings.Builder

	if len(s.types) != 0 {
		str.WriteString("Types:\n")
		for name, typ := range s.types {
			str.WriteString("\t" + name + "\n")
			str.WriteString(typ.schema.rootNode.IndentedTreeString(2) + "\n")
		}
	}

	if s.rootNode != nil {
		str.WriteString("Schema root node:\n" + s.rootNode.IndentedTreeString(1) + "\n")
	}

	return str.String()
}

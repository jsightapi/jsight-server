package ischema

import (
	"fmt"

	"github.com/jsightapi/jsight-schema-core/bytes"
	"github.com/jsightapi/jsight-schema-core/errs"
	"github.com/jsightapi/jsight-schema-core/fs"
)

type ISchema struct {
	// types the map where key is the name of the type (or included Schema).
	types    map[string]Type
	rootNode Node
}

func New() ISchema {
	return ISchema{
		types: make(map[string]Type, 5),
	}
}

func (s ISchema) TypesList() map[string]Type {
	return s.types
}

// MustType returns *ISchema or panic if not found.
func (s ISchema) MustType(name string) *ISchema {
	t, ok := s.types[name]
	if ok {
		return t.Schema
	}
	panic(errs.ErrUserTypeNotFound.F(name))
}

// Type returns specified type's schema.
func (s ISchema) Type(name string) (*ISchema, *errs.Err) {
	t, ok := s.types[name]
	if ok {
		return t.Schema, nil
	}

	return nil, errs.ErrUserTypeNotFound.F(name)
}

func (s ISchema) RootNode() Node {
	return s.rootNode
}

func (s *ISchema) AddNamedType(name string, typ *ISchema, rootFile *fs.File, begin bytes.Index) {
	if !bytes.NewBytes(name).IsUserTypeName() {
		panic(errs.ErrInvalidSchemaName.F(name))
	}
	s.addType(name, typ, rootFile, begin)
}

// AddUnnamedType Adds an unnamed TYPE to the SCHEMA. Returns a unique name for the added TYPE.
func (s *ISchema) AddUnnamedType(typ *ISchema, rootFile *fs.File, begin bytes.Index) string {
	name := fmt.Sprintf("#%p", typ)
	s.addType(name, typ, rootFile, begin)
	return name
}

func (s *ISchema) addType(name string, schema *ISchema, rootFile *fs.File, begin bytes.Index) {
	if _, ok := s.types[name]; ok {
		panic(errs.ErrDuplicationOfNameOfTypes.F(name))
	}
	s.types[name] = Type{schema, rootFile, begin}
}

func (s *ISchema) AddType(n string, t Type) {
	s.types[n] = t
}

func (s *ISchema) SetRootNode(node Node) {
	s.rootNode = node
}

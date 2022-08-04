package jschema

import (
	stdBytes "bytes"
	"fmt"

	"github.com/jsightapi/jsight-schema-go-library/bytes"
	"github.com/jsightapi/jsight-schema-go-library/errors"
	"github.com/jsightapi/jsight-schema-go-library/internal/sync"
	internalSchema "github.com/jsightapi/jsight-schema-go-library/notations/jschema/internal/schema"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema/internal/schema/constraint"
)

type exampleBuilder struct {
	// types all user types used in this schema.
	types map[string]internalSchema.Type

	// processedTypes an unordered set of processed types required for handling
	// recursion.
	// Infinity recursion can't happen here 'cause we check it before building
	// example, but optional recursion can be there.
	processedTypes map[string]int
}

func newExampleBuilder(types map[string]internalSchema.Type) *exampleBuilder {
	return &exampleBuilder{
		types:          types,
		processedTypes: map[string]int{},
	}
}

func (b *exampleBuilder) Build(node internalSchema.Node) ([]byte, error) {
	switch typedNode := node.(type) {
	case *internalSchema.ObjectNode:
		return b.buildExampleForObjectNode(typedNode)

	case *internalSchema.ArrayNode:
		return b.buildExampleForArrayNode(typedNode)

	case *internalSchema.LiteralNode:
		return typedNode.BasisLexEventOfSchemaForNode().Value(), nil

	case *internalSchema.MixedValueNode:
		return b.buildExampleForMixedValueNode(typedNode)

	default:
		return nil, fmt.Errorf("unhandled node type %T", node)
	}
}

func (b *exampleBuilder) buildExampleForObjectNode(node *internalSchema.ObjectNode) ([]byte, error) {
	if node.Constraint(constraint.TypesListConstraintType) != nil {
		return nil, errors.ErrUserTypeFound
	}

	buf := exampleBufferPool.Get()
	defer exampleBufferPool.Put(buf)

	buf.WriteRune('{')
	children := node.Children()
	length := len(children)
	for i, childNode := range children {
		ex, err := b.Build(childNode)
		if err != nil {
			return nil, err
		}

		if ex == nil {
			continue
		}

		k, err := b.buildObjectKey(node.Key(i))
		if err != nil {
			return nil, err
		}

		buf.WriteRune('"')
		buf.Write(k)
		buf.WriteString(`":`)
		buf.Write(ex)
		if i+1 != length {
			buf.WriteRune(',')
		}
	}
	buf.WriteRune('}')
	return buf.Bytes(), nil
}

func (b *exampleBuilder) buildObjectKey(k internalSchema.ObjectNodeKey) ([]byte, error) {
	if !bytes.Bytes(k.Key).IsUserTypeName() {
		return []byte(k.Key), nil
	}

	typ, ok := b.types[k.Key]
	if !ok {
		return nil, errors.Format(errors.ErrUnknownType, k.Key)
	}

	ex, err := b.Build(typ.Schema().RootNode())
	if err != nil {
		return nil, err
	}
	return stdBytes.Trim(ex, `"`), nil
}

func (b *exampleBuilder) buildExampleForArrayNode(node *internalSchema.ArrayNode) ([]byte, error) {
	if node.Constraint(constraint.TypesListConstraintType) != nil {
		return nil, errors.ErrUserTypeFound
	}

	buf := exampleBufferPool.Get()
	defer exampleBufferPool.Put(buf)

	buf.WriteRune('[')
	children := node.Children()
	length := len(children)
	for i, childNode := range children {
		ex, err := b.Build(childNode)
		if err != nil {
			return nil, err
		}

		if ex == nil {
			continue
		}

		buf.Write(ex)
		if i+1 != length {
			buf.WriteRune(',')
		}
	}
	buf.WriteRune(']')
	return buf.Bytes(), nil
}

func (b *exampleBuilder) buildExampleForMixedValueNode(node *internalSchema.MixedValueNode) ([]byte, error) {
	tt := node.GetTypes()
	if len(tt) == 0 {
		// Normally this shouldn't happen, but we still have to handle this case.
		return nil, errors.ErrLoader
	}

	typeName := tt[0]
	if !bytes.Bytes(typeName).IsUserTypeName() {
		return node.Value(), nil
	}

	if cnt := b.processedTypes[typeName]; cnt > 1 {
		// Do not process already processed type more than twice.
		return nil, nil
	}

	b.processedTypes[typeName]++
	defer func() {
		delete(b.processedTypes, typeName)
	}()

	t, ok := b.types[typeName]
	if !ok {
		return nil, errors.Format(errors.ErrTypeNotFound, typeName)
	}
	return b.Build(t.Schema().RootNode())
}

func buildExample(node internalSchema.Node, types map[string]internalSchema.Type) ([]byte, error) {
	switch typedNode := node.(type) {
	case *internalSchema.ObjectNode:
		return buildExampleForObjectNode(typedNode, types)

	case *internalSchema.ArrayNode:
		return buildExampleForArrayNode(typedNode, types)

	case *internalSchema.LiteralNode:
		return typedNode.BasisLexEventOfSchemaForNode().Value(), nil

	case *internalSchema.MixedValueNode:
		return buildExampleForMixedValueNode(typedNode, types)

	default:
		return nil, fmt.Errorf("unhandled node type %T", node)
	}
}

func buildExampleForObjectNode(
	node *internalSchema.ObjectNode,
	types map[string]internalSchema.Type,
) ([]byte, error) {
	if node.Constraint(constraint.TypesListConstraintType) != nil {
		return nil, errors.ErrUserTypeFound
	}

	b := exampleBufferPool.Get()
	defer exampleBufferPool.Put(b)

	b.WriteRune('{')
	children := node.Children()
	length := len(children)
	for i, childNode := range children {
		key := node.Key(i)
		b.WriteRune('"')
		b.WriteString(key.Key)
		b.WriteString(`":`)

		ex, err := buildExample(childNode, types)
		if err != nil {
			return nil, err
		}
		b.Write(ex)
		if i+1 != length {
			b.WriteRune(',')
		}
	}
	b.WriteRune('}')
	return b.Bytes(), nil
}

func buildExampleForArrayNode(
	node *internalSchema.ArrayNode,
	types map[string]internalSchema.Type,
) ([]byte, error) {
	if node.Constraint(constraint.TypesListConstraintType) != nil {
		return nil, errors.ErrUserTypeFound
	}

	b := exampleBufferPool.Get()
	defer exampleBufferPool.Put(b)

	b.WriteRune('[')
	children := node.Children()
	length := len(children)
	for i, childNode := range children {
		ex, err := buildExample(childNode, types)
		if err != nil {
			return nil, err
		}
		b.Write(ex)
		if i+1 != length {
			b.WriteRune(',')
		}
	}
	b.WriteRune(']')
	return b.Bytes(), nil
}

var exampleBufferPool = sync.NewBufferPool(512)

func buildExampleForMixedValueNode(
	node *internalSchema.MixedValueNode,
	types map[string]internalSchema.Type,
) ([]byte, error) {
	tt := node.GetTypes()
	if len(tt) == 0 {
		// Normally this shouldn't happen, but we still have to handle this case.
		return nil, errors.ErrLoader
	}

	typeName := tt[0]
	if !bytes.Bytes(typeName).IsUserTypeName() {
		return node.Value(), nil
	}

	t, ok := types[typeName]
	if !ok {
		return nil, errors.Format(errors.ErrTypeNotFound, typeName)
	}
	return buildExample(t.Schema().RootNode(), types)
}

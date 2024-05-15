package jschema

import (
	stdBytes "bytes"

	"github.com/jsightapi/jsight-schema-core/bytes"
	"github.com/jsightapi/jsight-schema-core/errs"
	"github.com/jsightapi/jsight-schema-core/internal/sync"
	"github.com/jsightapi/jsight-schema-core/notations/jschema/ischema"
	"github.com/jsightapi/jsight-schema-core/notations/jschema/ischema/constraint"
)

type exampleBuilder struct {
	// types all user types used in this schema.
	types map[string]ischema.Type

	// processedTypes an unordered set of processed types required for handling
	// recursion.
	// Infinite recursion can't happen here 'cause we check it before building
	// example, but optional recursion can be there.
	processedTypes map[string]int
}

func newExampleBuilder(types map[string]ischema.Type) *exampleBuilder {
	return &exampleBuilder{
		types:          types,
		processedTypes: map[string]int{},
	}
}

func (b *exampleBuilder) Build(node ischema.Node) ([]byte, error) {
	switch typedNode := node.(type) {
	case *ischema.ObjectNode:
		return b.buildExampleForObjectNode(typedNode)

	case *ischema.ArrayNode:
		return b.buildExampleForArrayNode(typedNode)

	case *ischema.LiteralNode:
		return typedNode.BasisLexEventOfSchemaForNode().Value().Data(), nil

	case *ischema.MixedValueNode:
		return b.buildExampleForMixedValueNode(typedNode)

	default:
		return nil, errs.ErrRuntimeFailure.F()
	}
}

func (b *exampleBuilder) buildExampleForObjectNode(node *ischema.ObjectNode) ([]byte, error) {
	if node.Constraint(constraint.TypesListConstraintType) != nil {
		return nil, errs.ErrUserTypeFound.F()
	}

	buf := exampleBufferPool.Get()
	defer exampleBufferPool.Put(buf)

	buf.WriteByte('{')
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

		buf.WriteByte('"')
		buf.Write(k)
		buf.WriteString(`":`)
		buf.Write(ex)
		if i+1 != length {
			buf.WriteByte(',')
		}
	}
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

func (b *exampleBuilder) buildObjectKey(k ischema.ObjectNodeKey) ([]byte, error) {
	if !k.IsShortcut {
		return []byte(k.Key), nil
	}

	typ, ok := b.types[k.Key]
	if !ok {
		return nil, errs.ErrUserTypeNotFound.F(k.Key)
	}

	ex, err := b.Build(typ.Schema.RootNode())
	if err != nil {
		return nil, err
	}
	return stdBytes.Trim(ex, `"`), nil
}

func (b *exampleBuilder) buildExampleForArrayNode(node *ischema.ArrayNode) ([]byte, error) {
	if node.Constraint(constraint.TypesListConstraintType) != nil {
		return nil, errs.ErrUserTypeFound.F()
	}

	buf := exampleBufferPool.Get()
	defer exampleBufferPool.Put(buf)

	buf.WriteByte('[')
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
			buf.WriteByte(',')
		}
	}
	buf.WriteByte(']')
	return buf.Bytes(), nil
}

func (b *exampleBuilder) buildExampleForMixedValueNode(node *ischema.MixedValueNode) ([]byte, error) {
	tt := node.GetTypes()
	if len(tt) == 0 {
		// Normally this shouldn't happen, but we still have to handle this case.
		return nil, errs.ErrLoader.F()
	}

	typeName := tt[0]
	if !bytes.NewBytes(typeName).IsUserTypeName() {
		return node.Value().Data(), nil
	}

	if cnt := b.processedTypes[typeName]; cnt > 1 {
		// Do not process already processed type more than twice.
		return nil, nil
	}

	b.processedTypes[typeName]++
	defer func() {
		b.processedTypes[typeName]--
	}()

	t, ok := b.types[typeName]
	if !ok {
		return nil, errs.ErrUserTypeNotFound.F(typeName)
	}
	return b.Build(t.Schema.RootNode())
}

func buildExample(node ischema.Node, types map[string]ischema.Type) ([]byte, error) {
	switch typedNode := node.(type) {
	case *ischema.ObjectNode:
		return buildExampleForObjectNode(typedNode, types)

	case *ischema.ArrayNode:
		return buildExampleForArrayNode(typedNode, types)

	case *ischema.LiteralNode:
		return typedNode.BasisLexEventOfSchemaForNode().Value().Data(), nil

	case *ischema.MixedValueNode:
		return buildExampleForMixedValueNode(typedNode, types)

	default:
		return nil, errs.ErrRuntimeFailure.F()
	}
}

func buildExampleForObjectNode(
	node *ischema.ObjectNode,
	types map[string]ischema.Type,
) ([]byte, error) {
	if node.Constraint(constraint.TypesListConstraintType) != nil {
		return nil, errs.ErrUserTypeFound.F()
	}

	b := exampleBufferPool.Get()
	defer exampleBufferPool.Put(b)

	b.WriteByte('{')
	children := node.Children()
	length := len(children)
	for i, childNode := range children {
		key := node.Key(i)
		b.WriteByte('"')
		b.WriteString(key.Key)
		b.WriteString(`":`)

		ex, err := buildExample(childNode, types)
		if err != nil {
			return nil, err
		}
		b.Write(ex)
		if i+1 != length {
			b.WriteByte(',')
		}
	}
	b.WriteByte('}')
	return b.Bytes(), nil
}

func buildExampleForArrayNode(
	node *ischema.ArrayNode,
	types map[string]ischema.Type,
) ([]byte, error) {
	if node.Constraint(constraint.TypesListConstraintType) != nil {
		return nil, errs.ErrUserTypeFound.F()
	}

	b := exampleBufferPool.Get()
	defer exampleBufferPool.Put(b)

	b.WriteByte('[')
	children := node.Children()
	length := len(children)
	for i, childNode := range children {
		ex, err := buildExample(childNode, types)
		if err != nil {
			return nil, err
		}
		b.Write(ex)
		if i+1 != length {
			b.WriteByte(',')
		}
	}
	b.WriteByte(']')
	return b.Bytes(), nil
}

var exampleBufferPool = sync.NewBufferPool(512)

func buildExampleForMixedValueNode(
	node *ischema.MixedValueNode,
	types map[string]ischema.Type,
) ([]byte, error) {
	tt := node.GetTypes()
	if len(tt) == 0 {
		// Normally this shouldn't happen, but we still have to handle this case.
		return nil, errs.ErrLoader.F()
	}

	typeName := tt[0]
	if !bytes.NewBytes(typeName).IsUserTypeName() {
		return node.Value().Data(), nil
	}

	t, ok := types[typeName]
	if !ok {
		return nil, errs.ErrUserTypeNotFound.F(typeName)
	}
	return buildExample(t.Schema.RootNode(), types)
}

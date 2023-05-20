package ischema

import (
	"strings"
	"sync"

	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/bytes"
	"github.com/jsightapi/jsight-schema-core/errs"
	"github.com/jsightapi/jsight-schema-core/json"
	"github.com/jsightapi/jsight-schema-core/lexeme"
	"github.com/jsightapi/jsight-schema-core/notations/jschema/ischema/constraint"
)

type ObjectNode struct {
	mu sync.Mutex

	// children node list.
	children []Node

	// keys stores the index of the node on the map for quick search.
	keys *ObjectNodeKeys

	baseNode

	// waitingForChild indicates that the Grow method will create a child node by
	// getting the next lexeme.
	waitingForChild bool
}

var _ Node = &ObjectNode{}

func NewObjectNode(lex lexeme.LexEvent) *ObjectNode {
	n := ObjectNode{
		baseNode: newBaseNode(lex),
		children: make([]Node, 0, 10),
		keys:     newObjectNodeKeys(),
	}
	n.setJsonType(json.TypeObject)
	return &n
}

func (*ObjectNode) Type() json.Type {
	return json.TypeObject
}

func (n *ObjectNode) Grow(lex lexeme.LexEvent) (Node, bool) {
	if n.waitingForChild {
		n.waitingForChild = false
		child := NewNode(lex)
		n.addChild(child)
		return child, true
	}

	switch lex.Type() {
	case lexeme.ObjectBegin, lexeme.ObjectKeyBegin, lexeme.ObjectValueEnd:

	case lexeme.KeyShortcutEnd:
		key := lex.Value().Unquote().String()
		n.AddKey(key, lex.Value().IsUserTypeName(), lex) // can panic

	case lexeme.ObjectKeyEnd:
		key := lex.Value().Unquote().String()
		n.AddKey(key, lex.Value().IsUserTypeName(), lex) // can panic

	case lexeme.ObjectValueBegin:
		n.waitingForChild = true

	case lexeme.ObjectEnd:
		return n.parent, false

	default:
		panic(errs.ErrUnexpectedLexicalEvent.F(lex.Type().String(), "in object node"))
	}

	return n, false
}

func (n *ObjectNode) Children() []Node {
	return n.children
}

func (n *ObjectNode) Len() int {
	return len(n.children)
}

// ChildByRawKey returns child by raw key as is present in schema.
// For instance: "foo" or @foo (shortcut).
func (n *ObjectNode) ChildByRawKey(rawKey bytes.Bytes) (Node, bool) {
	key := rawKey
	isShortcut := rawKey.IsUserTypeName()
	if !isShortcut {
		key = key.Unquote()
	}
	return n.Child(key.String(), isShortcut)
}

// Child returns child bye specified key.
func (n *ObjectNode) Child(key string, isShortcut bool) (Node, bool) {
	i, ok := n.keys.Get(key, isShortcut)
	if ok {
		return n.children[i.Index], true
	}
	return nil, false
}

func (n *ObjectNode) AddKey(key string, isShortcut bool, lex lexeme.LexEvent) {
	// Save child node index into map for faster search.
	n.keys.Set(ObjectNodeKey{
		Key:        key,
		Index:      len(n.children),
		IsShortcut: isShortcut,
		Lex:        lex,
	})
}

func (n *ObjectNode) addChild(child Node) {
	child.SetParent(n)
	n.children = append(n.children, child)
}

func (n *ObjectNode) AddChild(key ObjectNodeKey, child Node) {
	n.AddKey(key.Key, key.IsShortcut, key.Lex) // can panic
	n.addChild(child)
}

func (n *ObjectNode) Key(index int) ObjectNodeKey {
	if kv, ok := n.keys.Find(index); ok {
		return kv
	}
	panic(errs.ErrRuntimeFailure.F())
}

func (n *ObjectNode) Keys() *ObjectNodeKeys {
	return n.keys
}

func (n *ObjectNode) ASTNode() (schema.ASTNode, error) {
	an := astNodeFromNode(n)

	var err error
	an.Children, err = n.collectASTProperties()
	if err != nil {
		return schema.ASTNode{}, err
	}

	return an, nil
}

func (n *ObjectNode) collectASTProperties() ([]schema.ASTNode, error) {
	if len(n.keys.Data) == 0 {
		return nil, nil
	}

	pp := make([]schema.ASTNode, 0, len(n.keys.Data))

	for _, v := range n.keys.Data {
		c := n.children[v.Index]
		cn, err := c.ASTNode()
		if err != nil {
			return pp, err
		}

		cn.IsKeyShortcut = v.IsShortcut
		cn.Key = v.Key

		pp = append(pp, cn)
	}

	return pp, nil
}

func (n *ObjectNode) Copy() Node {
	nn := n.copyBase()
	nn.copyChildrenFrom(n)
	nn.copyKeysFrom(n)
	return &nn
}

func (n *ObjectNode) CopyAndLowercaseKeys() *ObjectNode {
	nn := n.copyBase()
	nn.copyChildrenFrom(n)
	nn.copyLowercaseKeysFrom(n)
	return &nn
}

func (n *ObjectNode) copyBase() ObjectNode {
	return ObjectNode{
		baseNode: n.baseNode.Copy(),
	}
}

func (n *ObjectNode) copyChildrenFrom(from *ObjectNode) {
	n.children = make([]Node, 0, len(from.children))
	for _, v := range from.children {
		n.children = append(n.children, v.Copy())
	}
}

func (n *ObjectNode) copyKeysFrom(from *ObjectNode) {
	n.keys = newObjectNodeKeys()
	for _, v := range from.keys.Data {
		n.keys.Set(v)
	}
}

func (n *ObjectNode) copyLowercaseKeysFrom(from *ObjectNode) {
	n.keys = newObjectNodeKeys()
	for _, v := range from.keys.Data {
		v.Key = strings.ToLower(v.Key)
		n.keys.Set(v)
	}
}

func (n *ObjectNode) EnsureAdditionalProperties() {
	n.mu.Lock()
	if c := n.Constraint(constraint.AdditionalPropertiesConstraintType); c == nil {
		c = constraint.NewAdditionalProperties(bytes.NewBytes("true"))
		n.AddConstraint(c)
	}
	n.mu.Unlock()
}

package schema

import (
	"fmt"

	jschema "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/bytes"
	"github.com/jsightapi/jsight-schema-go-library/internal/json"
	"github.com/jsightapi/jsight-schema-go-library/internal/lexeme"
)

type ObjectNode struct {
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

func newObjectNode(lex lexeme.LexEvent) *ObjectNode {
	n := ObjectNode{
		baseNode: newBaseNode(lex),
		children: make([]Node, 0, 10),
		keys:     newObjectNodeKeys(),
	}
	n.setJsonType(json.TypeObject)
	return &n
}

func (ObjectNode) Type() json.Type {
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
		n.addKey(key, lex.Value().IsUserTypeName(), lex) // can panic

	case lexeme.ObjectKeyEnd:
		key := lex.Value().Unquote().String()
		n.addKey(key, lex.Value().IsUserTypeName(), lex) // can panic

	case lexeme.ObjectValueBegin:
		n.waitingForChild = true

	case lexeme.ObjectEnd:
		return n.parent, false

	default:
		panic(`Unexpected lexical event "` + lex.Type().String() + `" in object node`)
	}

	return n, false
}

func (n ObjectNode) Children() []Node {
	return n.children
}

func (n ObjectNode) Len() int {
	return len(n.children)
}

// ChildByRawKey returns child by raw key as is present in schema.
// For instance: "foo" or @foo (shortcut).
func (n ObjectNode) ChildByRawKey(rawKey bytes.Bytes) (Node, bool) {
	key := rawKey
	isShortcut := rawKey.IsUserTypeName()
	if !isShortcut {
		key = key.Unquote()
	}
	return n.Child(key.String(), isShortcut)
}

// Child returns child bye specified key.
func (n ObjectNode) Child(key string, isShortcut bool) (Node, bool) {
	i, ok := n.keys.Get(key, isShortcut)
	if ok {
		return n.children[i.Index], true
	}
	return nil, false
}

func (n *ObjectNode) addKey(key string, isShortcut bool, lex lexeme.LexEvent) {
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
	n.addKey(key.Key, key.IsShortcut, key.Lex) // can panic
	n.addChild(child)
}

func (n ObjectNode) Key(index int) ObjectNodeKey {
	if kv, ok := n.keys.Find(index); ok {
		return kv
	}
	panic(fmt.Sprintf(`Schema key not found in index %d`, index))
}

func (n ObjectNode) Keys() *ObjectNodeKeys {
	return n.keys
}

func (n *ObjectNode) ASTNode() (jschema.ASTNode, error) {
	an := astNodeFromNode(n)

	var err error
	an.Children, err = n.collectASTProperties()
	if err != nil {
		return jschema.ASTNode{}, err
	}

	return an, nil
}

func (n *ObjectNode) collectASTProperties() ([]jschema.ASTNode, error) {
	if len(n.keys.Data) == 0 {
		return nil, nil
	}

	pp := make([]jschema.ASTNode, 0, len(n.keys.Data))

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

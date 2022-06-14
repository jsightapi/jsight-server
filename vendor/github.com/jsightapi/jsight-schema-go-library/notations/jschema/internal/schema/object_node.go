package schema

import (
	"fmt"
	"strings"
	"sync"

	jschema "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/errors"
	"github.com/jsightapi/jsight-schema-go-library/internal/json"
	"github.com/jsightapi/jsight-schema-go-library/internal/lexeme"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema/internal/schema/constraint"
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

// gen:OrderedMap
type ObjectNodeKeys struct {
	data  map[string]InnerObjectNodeKey
	order []string
	mx    sync.RWMutex
}

type InnerObjectNodeKey struct {
	Lex        lexeme.LexEvent
	Index      int
	IsShortcut bool
}

var _ Node = &ObjectNode{}

func newObjectNode(lex lexeme.LexEvent) *ObjectNode {
	n := ObjectNode{
		baseNode: newBaseNode(lex),
		children: make([]Node, 0, 10),
		keys:     &ObjectNodeKeys{},
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

func (n ObjectNode) Child(key string) (Node, bool) {
	i, ok := n.keys.Get(key)
	if ok {
		return n.children[i.Index], true
	}
	return nil, false
}

func (n *ObjectNode) addKey(key string, isShortcut bool, lex lexeme.LexEvent) {
	if n.keys.Has(key) {
		panic(errors.Format(errors.ErrDuplicateKeysInSchema, key))
	}

	// Save child node index into map for faster search.
	n.keys.Set(key, InnerObjectNodeKey{
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
	n.addKey(key.Name, key.IsShortcut, key.Lex) // can panic
	n.addChild(child)
}

type ObjectNodeKey struct {
	Name       string
	Lex        lexeme.LexEvent
	IsShortcut bool
}

func (n ObjectNode) Key(index int) ObjectNodeKey {
	kv, ok := n.keys.Find(func(k string, v InnerObjectNodeKey) bool {
		return v.Index == index
	})

	if !ok {
		panic(fmt.Sprintf(`Schema key not found in index %d`, index))
	}

	return ObjectNodeKey{
		Name:       kv.Key,
		IsShortcut: kv.Value.IsShortcut,
		Lex:        kv.Value.Lex,
	}
}

func (n ObjectNode) Keys() *ObjectNodeKeys {
	return n.keys
}

func (n ObjectNode) IndentedTreeString(depth int) string {
	indent := strings.Repeat("\t", depth)

	var str strings.Builder
	str.WriteString(n.IndentedNodeString(depth))

	for index, childNode := range n.children {
		key := n.Key(index) // can panic: Index not found in array
		str.WriteString(indent + "\t\"" + key.Name + "\":\n")
		str.WriteString(childNode.IndentedTreeString(depth + 2))
	}

	return str.String()
}

func (n ObjectNode) IndentedNodeString(depth int) string {
	indent := strings.Repeat("\t", depth)

	var str strings.Builder
	str.WriteString(indent + "* " + n.Type().String() + "\n")

	n.constraints.EachSafe(func(k constraint.Type, v constraint.Constraint) {
		str.WriteString(indent + "* " + v.String() + "\n")
	})

	return str.String()
}

func (n *ObjectNode) ASTNode() (jschema.ASTNode, error) {
	an := astNodeFromNode(n)

	var err error
	an.Properties, err = n.collectASTProperties()
	if err != nil {
		return jschema.ASTNode{}, err
	}

	return an, nil
}

func (n *ObjectNode) collectASTProperties() (*jschema.ASTNodes, error) {
	pp := &jschema.ASTNodes{}

	err := n.keys.Each(func(k string, v InnerObjectNodeKey) error {
		c := n.children[v.Index]
		cn, err := c.ASTNode()
		if err != nil {
			return err
		}

		cn.IsKeyShortcut = v.IsShortcut

		pp.Set(k, cn)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return pp, nil
}

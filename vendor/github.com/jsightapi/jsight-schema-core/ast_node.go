package schema

import "sync"

// ASTNode an AST node.
type ASTNode struct {
	// TokenType corresponding JSON type for this AST node's value.
	TokenType TokenType

	// SchemaType corresponding schema type for this AST node's value.
	SchemaType string

	// Key a node key (if this is the property of the object).
	Key string

	// Value a node value.
	// Make sense only for scalars and shortcuts.
	Value string

	// Comment a ast node comment.
	Comment string

	// Rules a map of attached rules.
	Rules *RuleASTNodes

	// Children contains all array items and object properties.
	// Make sense only for arrays and object.
	Children []ASTNode

	// IsKeyShortcut will be true if this property key is shortcut.
	// Make sense only for AST nodes which are represents object property.
	IsKeyShortcut bool

	// InheritedFrom a user type from which this property is inherited.
	InheritedFrom string
}

func (c *ASTNode) ObjectProperty(k string) *ASTNode {
	for i := range c.Children {
		if c.Children[i].Key == k {
			return &c.Children[i]
		}
	}
	return nil
}

func (c *ASTNode) Unshift(n ASTNode) {
	c.Children = append([]ASTNode{n}, c.Children...)
}

// ASTNodes an ordered map of AST nodes.
// gen:OrderedMap
type ASTNodes struct {
	data  map[string]ASTNode
	order []string
	mx    sync.RWMutex
}

type RuleASTNode struct {
	// TokenType corresponding JSON type for this AST node's value.
	TokenType TokenType

	// Value a node value.
	// Make sense only for scalars and shortcuts.
	Value string

	// Comment a ast node comment.
	Comment string

	// Properties contains all object properties.
	// Make sense only for objects.
	Properties *RuleASTNodes

	// Items contains all array items.
	// Make sense only for arrays.
	Items []RuleASTNode

	// Source a source of this rule.
	Source RuleASTNodeSource
}

func NewRuleASTNodes(data map[string]RuleASTNode, order []string) *RuleASTNodes {
	return &RuleASTNodes{
		data:  data,
		order: order,
	}
}

func MakeRuleASTNodes(capacity int) *RuleASTNodes {
	return &RuleASTNodes{
		data:  make(map[string]RuleASTNode, capacity),
		order: make([]string, 0, capacity),
	}
}

type RuleASTNodeSource int

const (
	RuleASTNodeSourceUnknown RuleASTNodeSource = iota

	// RuleASTNodeSourceManual indicates rule added manually by the user.
	RuleASTNodeSourceManual

	// RuleASTNodeSourceGenerated indicates rule generated inside the code.
	RuleASTNodeSourceGenerated
)

// RuleASTNodes an ordered map of rule AST nodes.
// gen:OrderedMap
type RuleASTNodes struct {
	data  map[string]RuleASTNode
	order []string
	mx    sync.RWMutex
}

// Rule represents a custom user-defined rule.
type Rule interface {
	// Len returns length of this rule in bytes.
	// Might return ParsingError if rule isn't valid.
	Len() (uint, error)

	// Check checks this rule is valid.
	// Can return ParsingError if rule isn't valid.
	Check() error

	// GetAST returns a root AST node for this schema.
	GetAST() (ASTNode, error)
}

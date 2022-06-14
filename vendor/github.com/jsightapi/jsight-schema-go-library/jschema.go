package jschema

import (
	"sync"

	"github.com/jsightapi/jsight-schema-go-library/internal/lexeme"
)

// Document represents a document.
// It's a concrete data. Data maybe a scalar type or complex type.
//
// Not a thead safe!
//
// Example of the valid documents:
// - "foo"
// - [1, 2, 3]
// - {"foo": "bar"}
type Document interface {
	// NextLexeme returns next lexeme from this document.
	// Might return ParsingError if document isn't valid.
	// Will return io.EOF when no more lexemes are available.
	NextLexeme() (lexeme.LexEvent, error)

	// Len returns length of document in bytes.
	// Might return ParsingError if document isn't valid.
	Len() (uint, error)

	// Check checks that this document is valid.
	// Can return ParsingError if document isn't valid.
	Check() error
}

// Schema represents a schema.
// Schema is a some description of expected structure of payload.
type Schema interface {
	// Len returns length of this schema in bytes.
	// Might return ParsingError if schema isn't valid.
	Len() (uint, error)

	// Example returns an example for this schema.
	// Might return ParsingError if schema isn't valid.
	Example() ([]byte, error)

	// AddType adds a new type to this schema.
	// Might return a ParsingError if add type isn't valid.
	AddType(name string, schema Schema) error

	// AddRule adds a new type to this schema.
	// Might return a ParsingError if add type isn't valid.
	AddRule(name string, schema Rule) error

	// Check checks that this schema is valid.
	// Can return ParsingError if schema isn't valid.
	Check() error

	// Validate validates specified document.
	// Might return a ParsingError if schema isn't valid.
	// Can return a ValidationError if specified document isn't valid for this
	// schema.
	Validate(Document) error

	// GetAST returns a root AST node for this schema.
	GetAST() (ASTNode, error)

	// UsedUserTypes return all used user types.
	UsedUserTypes() ([]string, error)
}

// ASTNode an AST node.
type ASTNode struct {
	// JSONType corresponding JSON type for this AST node's value.
	JSONType JSONType

	// SchemaType corresponding schema type for this AST node's value.
	SchemaType string

	// Value a node value.
	// Make sense only for scalars and shortcuts.
	Value string

	// Comment a ast node comment.
	Comment string

	// Rules a map of attached rules.
	Rules *RuleASTNodes

	// Properties contains all object properties.
	// Make sense only for objects.
	Properties *ASTNodes

	// Items contains all array items.
	// Make sense only for arrays.
	Items []ASTNode

	// IsKeyShortcut will be true if this property key is shortcut.
	// Make sense only for AST nodes which are represents object property.
	IsKeyShortcut bool
}

// ASTNodes an ordered map of AST nodes.
// gen:OrderedMap
type ASTNodes struct {
	data  map[string]ASTNode
	order []string
	mx    sync.RWMutex
}

func NewASTNodes(data map[string]ASTNode, order []string) *ASTNodes {
	return &ASTNodes{
		data:  data,
		order: order,
	}
}

type RuleASTNode struct {
	// JSONType corresponding JSON type for this AST node's value.
	JSONType JSONType

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
type Rule interface{}

// ParsingError indicates something bad was happened during parsing.
type ParsingError interface {
	error

	// Position returns position of buggy character.
	Position() uint

	// Message returns an error message.
	Message() string

	// ErrCode returns an error code.
	ErrCode() int
}

// ValidationError indicates that validation was failed.
type ValidationError interface {
	error

	// Message returns an error message.
	Message() string

	// ErrCode returns an error code.
	ErrCode() int
}

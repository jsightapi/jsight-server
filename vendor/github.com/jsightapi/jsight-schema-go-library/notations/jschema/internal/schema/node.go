package schema

//go:generate mockery --name Node --output ../mocks

import (
	"sync"

	jschema "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/bytes"
	"github.com/jsightapi/jsight-schema-go-library/internal/json"
	"github.com/jsightapi/jsight-schema-go-library/internal/lexeme"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema/internal/schema/constraint"
)

// The node of the internal representation of the scheme.
// Roughly corresponds to the JSON element in the EXAMPLE of schema.
// Contains information about the constraints imposed on the node.

type Node interface {
	// Type returns type of this node.
	Type() json.Type

	SetRealType(string)
	RealType() string

	// Parent returns a parent of this node.
	Parent() Node

	// SetParent sets a parent for this node.
	SetParent(Node)

	// BasisLexEventOfSchemaForNode returns a LexEvent from the scheme on the
	// basis of which the node is created. It is used to check on the schemes for
	// compliance with the example and the list of constraints. Also used to display
	// an error.
	BasisLexEventOfSchemaForNode() lexeme.LexEvent

	// Grow this method receives the input lexical event from the scanner, fill
	// yourself with data from them. If necessary, creates children. Returns the
	// node to which you want to pass the next lexeme (yourself, child, or parent).
	Grow(lexeme.LexEvent) (Node, bool)

	// Constraint returns a constraint by its type.
	Constraint(constraint.Type) constraint.Constraint

	// AddConstraint adds a constraint to this node.
	AddConstraint(constraint.Constraint)

	// DeleteConstraint removes a constraint from this node.
	DeleteConstraint(constraint.Type)

	// ConstraintMap returns a list of constraints or nil (if empty).
	ConstraintMap() *Constraints

	// NumberOfConstraints returns the number of constraints.
	NumberOfConstraints() int

	// IndentedTreeString returns the text representation of the node and recursively
	// all children. You can specify the indent (number of tabs) that is placed
	// before each line in the text view.
	IndentedTreeString(int) string

	// IndentedNodeString returns the text representation of the node without children.
	// You can specify the indent (number of tabs) that is placed before each line
	// in the text view.
	IndentedNodeString(int) string

	// Value returns this node's value.
	Value() bytes.Bytes

	// ASTNode returns proper ASTNode for this node.
	ASTNode() (jschema.ASTNode, error)

	// SetComment sets a comment for this node.
	SetComment(string)

	// Comment returns this node comment.
	Comment() string
}

// Constraints an ordered map of node constraints.
// gen:OrderedMap
type Constraints struct {
	data  map[constraint.Type]constraint.Constraint
	order []constraint.Type
	mx    sync.RWMutex
}

// BranchNode that can contain child elements (an array or an object).
type BranchNode interface {
	Children() []Node
	Len() int
}

func NewNode(lex lexeme.LexEvent) Node {
	switch lex.Type() { //nolint:exhaustive // We will throw a panic in over cases.
	case lexeme.LiteralBegin:
		return newLiteralNode(lex)
	case lexeme.ObjectBegin:
		return newObjectNode(lex)
	case lexeme.ArrayBegin:
		return newArrayNode(lex)
	case lexeme.MixedValueBegin:
		return NewMixedValueNode(lex)
	}
	panic(`Can not create node from the lexical event "` + lex.Type().String() + `"`)
}

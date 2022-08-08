package schema

import (
	"github.com/jsightapi/jsight-schema-go-library/bytes"
	"github.com/jsightapi/jsight-schema-go-library/errors"
	"github.com/jsightapi/jsight-schema-go-library/internal/json"
	"github.com/jsightapi/jsight-schema-go-library/internal/lexeme"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema/internal/schema/constraint"
)

type baseNode struct {
	// parent a parent node.
	parent Node

	realType string

	// comment a node comment.
	comment string

	// constraints a list of this node constraints.
	constraints *Constraints

	// schemaLexEvent used to check and display an error if the node value does
	// not match the constraints.
	schemaLexEvent lexeme.LexEvent

	// jsonType a JSON type for this node.
	jsonType json.Type
}

func newBaseNode(lex lexeme.LexEvent) baseNode {
	return baseNode{
		parent:         nil,
		jsonType:       json.TypeUndefined,
		schemaLexEvent: lex,
		constraints:    &Constraints{},
	}
}

func (n baseNode) Type() json.Type {
	return n.jsonType
}

func (n *baseNode) SetRealType(s string) bool {
	// Make sure current real type is compatible with node type.
	avail, ok := compatibleTypes[s]
	if !ok {
		return false
	}

	if _, ok := avail[n.jsonType]; !ok {
		return false
	}

	n.realType = s
	return true
}

var compatibleTypes = map[string]map[json.Type]struct{}{
	"mixed": availableJSONTypes(
		json.TypeObject,
		json.TypeArray,
		json.TypeString,
		json.TypeInteger,
		json.TypeFloat,
		json.TypeBoolean,
		json.TypeNull,
		json.TypeMixed,
	),
	"enum": availableJSONTypes(
		json.TypeString,
		json.TypeInteger,
		json.TypeFloat,
		json.TypeBoolean,
		json.TypeNull,
	),
	"any": availableJSONTypes(
		json.TypeObject,
		json.TypeArray,
		json.TypeString,
		json.TypeInteger,
		json.TypeFloat,
		json.TypeBoolean,
		json.TypeNull,
		json.TypeMixed,
	),
	"decimal":  availableJSONTypes(json.TypeFloat),
	"email":    availableJSONTypes(json.TypeString),
	"uri":      availableJSONTypes(json.TypeString),
	"uuid":     availableJSONTypes(json.TypeString),
	"date":     availableJSONTypes(json.TypeString),
	"datetime": availableJSONTypes(json.TypeString),
	"object":   availableJSONTypes(json.TypeObject),
	"array":    availableJSONTypes(json.TypeArray),
	"string":   availableJSONTypes(json.TypeString),
	"integer":  availableJSONTypes(json.TypeInteger),
	"float":    availableJSONTypes(json.TypeFloat),
	"boolean":  availableJSONTypes(json.TypeBoolean),
	"null":     availableJSONTypes(json.TypeNull),
}

func availableJSONTypes(tt ...json.Type) map[json.Type]struct{} {
	res := map[json.Type]struct{}{}
	for _, t := range tt {
		res[t] = struct{}{}
	}
	return res
}

func (n *baseNode) RealType() string {
	if n.realType == "" {
		return n.jsonType.String()
	}
	return n.realType
}

func (n *baseNode) setJsonType(t json.Type) {
	n.jsonType = t
}

func (n baseNode) Parent() Node {
	return n.parent
}

func (n *baseNode) SetParent(parent Node) {
	n.parent = parent
}

func (n baseNode) BasisLexEventOfSchemaForNode() lexeme.LexEvent {
	return n.schemaLexEvent
}

// Constraint returns requested Constraint if found.
func (n baseNode) Constraint(t constraint.Type) constraint.Constraint {
	if n.constraints == nil {
		return nil
	}
	c, ok := n.constraints.Get(t)
	if ok {
		return c
	}
	return nil
}

// AddConstraint adds new constraint to this node.
// Won't add if c is nil.
func (n *baseNode) AddConstraint(c constraint.Constraint) {
	if c == nil {
		return
	}

	if n.constraints.Has(c.Type()) { // find an existing constraint
		panic(errors.Format(errors.ErrDuplicateRule, c.Type().String()))
	}

	n.constraints.Set(c.Type(), c)
}

func (n *baseNode) DeleteConstraint(t constraint.Type) {
	n.constraints.Delete(t)
}

// ConstraintMap returns all constraints.
func (n baseNode) ConstraintMap() *Constraints {
	return n.constraints
}

func (n baseNode) NumberOfConstraints() int {
	return n.constraints.Len()
}

func (n baseNode) Value() bytes.Bytes {
	return n.schemaLexEvent.Value()
}

func (n *baseNode) SetComment(s string) {
	n.comment = s
}

func (n *baseNode) Comment() string {
	return n.comment
}

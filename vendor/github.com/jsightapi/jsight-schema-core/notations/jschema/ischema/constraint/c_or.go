package constraint

import (
	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/json"
)

// Or constraint.
// Used for compile-time checking.
type Or struct {
	source schema.RuleASTNodeSource
}

var (
	_ Constraint = Or{}
	_ Constraint = (*Or)(nil)
)

func NewOr(s schema.RuleASTNodeSource) *Or {
	return &Or{
		source: s,
	}
}

func (c Or) IsGenerated() bool {
	return c.source == schema.RuleASTNodeSourceGenerated
}

func (Or) IsJsonTypeCompatible(json.Type) bool {
	return true
}

func (Or) Type() Type {
	return OrConstraintType
}

func (Or) String() string {
	return "[ UNVERIFIABLE CONSTRAINT ] " + OrConstraintType.String()
}

func (Or) ASTNode() schema.RuleASTNode {
	// Check `collectASTRules` function for the actual logic.
	return newEmptyRuleASTNode()
}

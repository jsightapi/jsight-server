package constraint

import (
	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/json"
)

type AnyConstraint struct{}

var (
	_ Constraint = AnyConstraint{}
	_ Constraint = (*AnyConstraint)(nil)
)

func NewAny() *AnyConstraint {
	return &AnyConstraint{}
}

func (AnyConstraint) IsJsonTypeCompatible(json.Type) bool {
	return true
}

func (AnyConstraint) Type() Type {
	return AnyConstraintType
}

func (AnyConstraint) String() string {
	return AnyConstraintType.String()
}

func (AnyConstraint) ASTNode() schema.RuleASTNode {
	return newEmptyRuleASTNode()
}

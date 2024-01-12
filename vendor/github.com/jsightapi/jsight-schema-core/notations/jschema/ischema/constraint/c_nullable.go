package constraint

import (
	"strconv"

	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/bytes"
	"github.com/jsightapi/jsight-schema-core/errs"
	"github.com/jsightapi/jsight-schema-core/json"
)

type Nullable struct {
	value bool
}

var (
	_ Constraint = Nullable{}
	_ Constraint = (*Nullable)(nil)
	_ BoolKeeper = Nullable{}
	_ BoolKeeper = (*Nullable)(nil)
)

func NewNullable(ruleValue bytes.Bytes) *Nullable {
	c := Nullable{}

	var err error
	if c.value, err = ruleValue.ParseBool(); err != nil {
		panic(errs.ErrInvalidValueOfConstraint.F(NullableConstraintType.String()))
	}
	return &c
}

func (Nullable) IsJsonTypeCompatible(json.Type) bool {
	return true
}

func (Nullable) Type() Type {
	return NullableConstraintType
}

func (c Nullable) String() string {
	if c.value {
		return NullableConstraintType.String() + colonTrue
	}
	return NullableConstraintType.String() + colonFalse
}

func (c Nullable) Bool() bool {
	return c.value
}

func (c Nullable) ASTNode() schema.RuleASTNode {
	return newRuleASTNode(schema.TokenTypeBoolean, strconv.FormatBool(c.value), schema.RuleASTNodeSourceManual)
}

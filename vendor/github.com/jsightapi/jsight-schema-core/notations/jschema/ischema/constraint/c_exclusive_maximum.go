package constraint //nolint:dupl // Duplicates exclusive minimum with small differences.

import (
	"strconv"

	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/bytes"
	"github.com/jsightapi/jsight-schema-core/errs"
	"github.com/jsightapi/jsight-schema-core/json"
)

type ExclusiveMaximum struct {
	exclusive bool
}

var (
	_ Constraint = ExclusiveMaximum{}
	_ Constraint = (*ExclusiveMaximum)(nil)
)

func NewExclusiveMaximum(ruleValue bytes.Bytes) *ExclusiveMaximum {
	c := ExclusiveMaximum{}
	var err error
	if c.exclusive, err = ruleValue.ParseBool(); err != nil {
		panic(errs.ErrInvalidValueOfConstraint.F(ExclusiveMaximumConstraintType.String()))
	}
	return &c
}

func (ExclusiveMaximum) IsJsonTypeCompatible(t json.Type) bool {
	return t == json.TypeInteger || t == json.TypeFloat
}

func (ExclusiveMaximum) Type() Type {
	return ExclusiveMaximumConstraintType
}

func (c ExclusiveMaximum) String() string {
	str := "[ UNVERIFIABLE CONSTRAINT ] " + ExclusiveMaximumConstraintType.String()
	if c.exclusive {
		str += colonTrue
	} else {
		str += colonFalse
	}
	return str
}

func (c ExclusiveMaximum) IsExclusive() bool {
	return c.exclusive
}

func (c ExclusiveMaximum) ASTNode() schema.RuleASTNode {
	return newRuleASTNode(schema.TokenTypeBoolean, strconv.FormatBool(c.exclusive), schema.RuleASTNodeSourceManual)
}

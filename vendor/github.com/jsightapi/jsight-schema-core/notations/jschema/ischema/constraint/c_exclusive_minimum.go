package constraint //nolint:dupl // Duplicates exclusive minimum with small differences.

import (
	"strconv"

	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/bytes"
	"github.com/jsightapi/jsight-schema-core/errs"
	"github.com/jsightapi/jsight-schema-core/json"
)

type ExclusiveMinimum struct {
	exclusive bool
}

var (
	_ Constraint = ExclusiveMinimum{}
	_ Constraint = (*ExclusiveMinimum)(nil)
)

func NewExclusiveMinimum(ruleValue bytes.Bytes) *ExclusiveMinimum {
	c := ExclusiveMinimum{}

	var err error
	if c.exclusive, err = ruleValue.ParseBool(); err != nil {
		panic(errs.ErrInvalidValueOfConstraint.F(ExclusiveMinimumConstraintType.String()))
	}
	return &c
}

func (ExclusiveMinimum) IsJsonTypeCompatible(t json.Type) bool {
	return t == json.TypeInteger || t == json.TypeFloat
}

func (ExclusiveMinimum) Type() Type {
	return ExclusiveMinimumConstraintType
}

func (c ExclusiveMinimum) String() string {
	str := "[ UNVERIFIABLE CONSTRAINT ] " + ExclusiveMinimumConstraintType.String()
	if c.exclusive {
		str += colonTrue
	} else {
		str += colonFalse
	}
	return str
}

func (c ExclusiveMinimum) IsExclusive() bool {
	return c.exclusive
}

func (c ExclusiveMinimum) ASTNode() schema.RuleASTNode {
	return newRuleASTNode(schema.TokenTypeBoolean, strconv.FormatBool(c.exclusive), schema.RuleASTNodeSourceManual)
}

package constraint

import (
	"strconv"

	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/bytes"
	"github.com/jsightapi/jsight-schema-core/errs"
	"github.com/jsightapi/jsight-schema-core/json"
)

type Optional struct {
	value bool
}

var (
	_ Constraint = Optional{}
	_ Constraint = (*Optional)(nil)
	_ BoolKeeper = Optional{}
	_ BoolKeeper = (*Optional)(nil)
)

func NewOptional(ruleValue bytes.Bytes) *Optional {
	c := Optional{}

	var err error
	if c.value, err = ruleValue.ParseBool(); err != nil {
		panic(errs.ErrInvalidValueOfConstraint.F(OptionalConstraintType.String()))
	}
	return &c
}

func (Optional) IsJsonTypeCompatible(json.Type) bool {
	return true
}

func (Optional) Type() Type {
	return OptionalConstraintType
}

func (c Optional) String() string {
	str := "[ UNVERIFIABLE CONSTRAINT ] " + OptionalConstraintType.String()
	if c.value {
		str += colonTrue
	} else {
		str += colonFalse
	}
	return str
}

func (c Optional) Bool() bool {
	return c.value
}

func (c Optional) ASTNode() schema.RuleASTNode {
	return newRuleASTNode(schema.TokenTypeBoolean, strconv.FormatBool(c.value), schema.RuleASTNodeSourceManual)
}

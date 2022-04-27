package constraint

import (
	"strconv"

	jschema "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/bytes"
	"github.com/jsightapi/jsight-schema-go-library/errors"
	"github.com/jsightapi/jsight-schema-go-library/internal/json"
)

type Nullable struct {
	value bool
}

var (
	_ Constraint = Nullable{}
	_ BoolKeeper = Nullable{}
)

func NewNullable(ruleValue bytes.Bytes) *Nullable {
	c := Nullable{}

	var err error
	if c.value, err = ruleValue.ParseBool(); err != nil {
		panic(errors.Format(errors.ErrInvalidValueOfConstraint, NullableConstraintType.String()))
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
		return NullableConstraintType.String() + ": true"
	}
	return NullableConstraintType.String() + ": false"
}

func (c Nullable) Bool() bool {
	return c.value
}

func (c Nullable) ASTNode() jschema.RuleASTNode {
	return newRuleASTNode(jschema.JSONTypeBoolean, strconv.FormatBool(c.value), jschema.RuleASTNodeSourceManual)
}

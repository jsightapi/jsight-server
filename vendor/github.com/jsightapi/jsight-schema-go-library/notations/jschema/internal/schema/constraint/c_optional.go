package constraint

import (
	"strconv"

	jschema "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/bytes"
	"github.com/jsightapi/jsight-schema-go-library/errors"
	"github.com/jsightapi/jsight-schema-go-library/internal/json"
)

type Optional struct {
	value bool
}

var _ Constraint = Optional{}

func NewOptional(ruleValue bytes.Bytes) *Optional {
	c := Optional{}

	var err error
	if c.value, err = ruleValue.ParseBool(); err != nil {
		panic(errors.Format(errors.ErrInvalidValueOfConstraint, OptionalConstraintType.String()))
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
		str += ": true"
	} else {
		str += ": false"
	}
	return str
}

func (c Optional) Bool() bool {
	return c.value
}

func (c Optional) ASTNode() jschema.RuleASTNode {
	return newRuleASTNode(jschema.JSONTypeBoolean, strconv.FormatBool(c.value), jschema.RuleASTNodeSourceManual)
}

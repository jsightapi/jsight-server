package constraint //nolint:dupl // Duplicates exclusive minimum with small differences.

import (
	"strconv"

	jschema "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/bytes"
	"github.com/jsightapi/jsight-schema-go-library/errors"
	"github.com/jsightapi/jsight-schema-go-library/internal/json"
)

type ExclusiveMinimum struct {
	exclusive bool
}

var _ Constraint = ExclusiveMinimum{}

func NewExclusiveMinimum(ruleValue bytes.Bytes) *ExclusiveMinimum {
	c := ExclusiveMinimum{}

	var err error
	if c.exclusive, err = ruleValue.ParseBool(); err != nil {
		panic(errors.Format(errors.ErrInvalidValueOfConstraint, ExclusiveMinimumConstraintType.String()))
	}
	return &c
}

func (ExclusiveMinimum) IsJsonTypeCompatible(t json.Type) bool {
	if t == json.TypeInteger || t == json.TypeFloat {
		return true
	}
	return false
}

func (ExclusiveMinimum) Type() Type {
	return ExclusiveMinimumConstraintType
}

func (c ExclusiveMinimum) String() string {
	str := "[ UNVERIFIABLE CONSTRAINT ] " + ExclusiveMinimumConstraintType.String()
	if c.exclusive {
		str += ": true"
	} else {
		str += ": false"
	}
	return str
}

func (c ExclusiveMinimum) IsExclusive() bool {
	return c.exclusive
}

func (c ExclusiveMinimum) ASTNode() jschema.RuleASTNode {
	return newRuleASTNode(jschema.JSONTypeBoolean, strconv.FormatBool(c.exclusive), jschema.RuleASTNodeSourceManual)
}

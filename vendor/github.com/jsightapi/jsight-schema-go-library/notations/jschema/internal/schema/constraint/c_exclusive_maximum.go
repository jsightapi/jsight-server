package constraint //nolint:dupl // Duplicates exclusive minimum with small differences.

import (
	"strconv"

	jschema "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/bytes"
	"github.com/jsightapi/jsight-schema-go-library/errors"
	"github.com/jsightapi/jsight-schema-go-library/internal/json"
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
		panic(errors.Format(errors.ErrInvalidValueOfConstraint, ExclusiveMaximumConstraintType.String()))
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
		str += ": true"
	} else {
		str += ": false"
	}
	return str
}

func (c ExclusiveMaximum) IsExclusive() bool {
	return c.exclusive
}

func (c ExclusiveMaximum) ASTNode() jschema.RuleASTNode {
	return newRuleASTNode(jschema.JSONTypeBoolean, strconv.FormatBool(c.exclusive), jschema.RuleASTNodeSourceManual)
}

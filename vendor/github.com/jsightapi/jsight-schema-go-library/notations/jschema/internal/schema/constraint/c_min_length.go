package constraint

import (
	jschema "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/bytes"
	"github.com/jsightapi/jsight-schema-go-library/errors"
	"github.com/jsightapi/jsight-schema-go-library/internal/json"
)

type MinLength struct {
	value    *json.Number
	rawValue bytes.Bytes
}

var _ Constraint = MinLength{}

func NewMinLength(ruleValue bytes.Bytes) *MinLength {
	number, err := json.NewIntegerNumber(ruleValue)
	if err != nil {
		panic(err)
	}

	return &MinLength{
		rawValue: ruleValue,
		value:    number,
	}
}

func (MinLength) IsJsonTypeCompatible(t json.Type) bool {
	return t == json.TypeString
}

func (MinLength) Type() Type {
	return MinLengthConstraintType
}

func (c MinLength) String() string {
	return MinLengthConstraintType.String() + ": " + c.value.String()
}

func (c MinLength) Validate(value bytes.Bytes) {
	length := len(value.Unquote())
	jsonLength := json.NewNumberFromInt(length)
	if jsonLength.LessThan(c.value) {
		panic(errors.Format(
			errors.ErrConstraintStringLengthValidation,
			MinLengthConstraintType.String(),
			c.value.String(),
		))
	}
}

func (c MinLength) ASTNode() jschema.RuleASTNode {
	return newRuleASTNode(jschema.JSONTypeNumber, c.rawValue.String(), jschema.RuleASTNodeSourceManual)
}

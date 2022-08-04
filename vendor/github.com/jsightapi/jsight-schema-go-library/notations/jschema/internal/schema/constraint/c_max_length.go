package constraint

import (
	jschema "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/bytes"
	"github.com/jsightapi/jsight-schema-go-library/errors"
	"github.com/jsightapi/jsight-schema-go-library/internal/json"
)

type MaxLength struct {
	value    *json.Number
	rawValue bytes.Bytes
}

var (
	_ Constraint       = MaxLength{}
	_ Constraint       = (*MaxLength)(nil)
	_ LiteralValidator = MaxLength{}
	_ LiteralValidator = (*MaxLength)(nil)
)

func NewMaxLength(ruleValue bytes.Bytes) *MaxLength {
	number, err := json.NewIntegerNumber(ruleValue)
	if err != nil {
		panic(err)
	}

	return &MaxLength{
		rawValue: ruleValue,
		value:    number,
	}
}

func (MaxLength) IsJsonTypeCompatible(t json.Type) bool {
	return t == json.TypeString
}

func (MaxLength) Type() Type {
	return MaxLengthConstraintType
}

func (c MaxLength) String() string {
	return MaxLengthConstraintType.String() + ": " + c.value.String()
}

func (c MaxLength) Validate(value bytes.Bytes) {
	length := len(value.Unquote())
	jsonLength := json.NewNumberFromInt(length)
	if jsonLength.GreaterThan(c.value) {
		panic(errors.Format(
			errors.ErrConstraintStringLengthValidation,
			MaxLengthConstraintType.String(),
			c.value.String(),
		))
	}
}

func (c MaxLength) ASTNode() jschema.RuleASTNode {
	return newRuleASTNode(jschema.JSONTypeNumber, c.rawValue.String(), jschema.RuleASTNodeSourceManual)
}

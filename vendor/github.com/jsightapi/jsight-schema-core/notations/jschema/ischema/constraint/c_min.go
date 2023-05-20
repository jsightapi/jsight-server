package constraint

import (
	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/bytes"
	"github.com/jsightapi/jsight-schema-core/errs"
	"github.com/jsightapi/jsight-schema-core/json"
)

type Min struct {
	min       *json.Number
	rawValue  bytes.Bytes
	exclusive bool
}

var (
	_ Constraint       = Min{}
	_ Constraint       = (*Min)(nil)
	_ LiteralValidator = Min{}
	_ LiteralValidator = (*Min)(nil)
)

func NewMin(ruleValue bytes.Bytes) *Min {
	number, err := json.NewNumber(ruleValue)
	if err != nil {
		panic(err)
	}

	return &Min{
		rawValue: ruleValue,
		min:      number,
	}
}

func (Min) IsJsonTypeCompatible(t json.Type) bool {
	return t == json.TypeInteger || t == json.TypeFloat
}

func (Min) Type() Type {
	return MinConstraintType
}

func (c Min) String() string {
	str := MinConstraintType.String() + ": " + c.min.String()
	if c.exclusive {
		return str + " (exclusive: true)"
	}
	return str
}

func (c *Min) SetExclusive(exclusive bool) {
	c.exclusive = exclusive
}

func (c *Min) Exclusive() bool {
	return c.exclusive
}

func (c Min) Validate(value bytes.Bytes) {
	jsonNumber, err := json.NewNumber(value)
	if err != nil {
		panic(err)
	}
	if c.exclusive {
		if c.min.GreaterThanOrEqual(jsonNumber) {
			panic(errs.ErrConstraintValidation.F(
				MinConstraintType.String(),
				c.min.String(),
				"(exclusive)",
			))
		}
	} else {
		if c.min.GreaterThan(jsonNumber) {
			panic(errs.ErrConstraintValidation.F(MinConstraintType.String(), c.min.String(), ""))
		}
	}
}

func (c Min) ASTNode() schema.RuleASTNode {
	return newRuleASTNode(schema.TokenTypeNumber, c.rawValue.String(), schema.RuleASTNodeSourceManual)
}

func (c *Min) Value() *json.Number {
	return c.min
}

package constraint

import (
	jschema "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/bytes"
	"github.com/jsightapi/jsight-schema-go-library/errors"
	"github.com/jsightapi/jsight-schema-go-library/internal/json"
)

type Max struct {
	max       *json.Number
	rawValue  bytes.Bytes
	exclusive bool
}

var (
	_ Constraint       = Max{}
	_ Constraint       = (*Max)(nil)
	_ LiteralValidator = Max{}
	_ LiteralValidator = (*Max)(nil)
)

func NewMax(ruleValue bytes.Bytes) *Max {
	number, err := json.NewNumber(ruleValue)
	if err != nil {
		panic(err)
	}

	return &Max{
		rawValue: ruleValue,
		max:      number,
	}
}

func (Max) IsJsonTypeCompatible(t json.Type) bool {
	return t == json.TypeInteger || t == json.TypeFloat
}

func (Max) Type() Type {
	return MaxConstraintType
}

func (c Max) String() string {
	str := MaxConstraintType.String() + ": " + c.max.String()
	if c.exclusive {
		return str + " (exclusive: true)"
	}
	return str
}

func (c *Max) SetExclusive(exclusive bool) {
	c.exclusive = exclusive
}

func (c *Max) Exclusive() bool {
	return c.exclusive
}

func (c Max) Validate(value bytes.Bytes) {
	jsonNumber, err := json.NewNumber(value)
	if err != nil {
		panic(err)
	}
	if c.exclusive {
		if c.max.LessThanOrEqual(jsonNumber) {
			panic(errors.Format(errors.ErrConstraintValidation, MaxConstraintType.String(), c.max.String(), "(exclusive)")) //nolint:lll
		}
	} else {
		if c.max.LessThan(jsonNumber) {
			panic(errors.Format(errors.ErrConstraintValidation, MaxConstraintType.String(), c.max.String(), ""))
		}
	}
}

func (c Max) ASTNode() jschema.RuleASTNode {
	return newRuleASTNode(jschema.TokenTypeNumber, c.rawValue.String(), jschema.RuleASTNodeSourceManual)
}

func (c *Max) Value() *json.Number {
	return c.max
}

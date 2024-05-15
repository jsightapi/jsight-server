package constraint

import (
	"strconv"

	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/bytes"
	"github.com/jsightapi/jsight-schema-core/errs"
	"github.com/jsightapi/jsight-schema-core/json"
)

type MinLength struct {
	value uint
}

var (
	_ Constraint       = MinLength{}
	_ Constraint       = (*MinLength)(nil)
	_ LiteralValidator = MinLength{}
	_ LiteralValidator = (*MinLength)(nil)
)

func NewMinLength(ruleValue bytes.Bytes) *MinLength {
	return &MinLength{
		value: parseUint(ruleValue, MinLengthConstraintType),
	}
}

func (MinLength) IsJsonTypeCompatible(t json.Type) bool {
	return t == json.TypeString
}

func (MinLength) Type() Type {
	return MinLengthConstraintType
}

func (c MinLength) String() string {
	return MinLengthConstraintType.String() + ": " + strconv.FormatUint(uint64(c.value), 10)
}

func (c MinLength) Validate(value bytes.Bytes) {
	length := uint(value.Unquote().Len())
	if length < c.value {
		panic(errs.ErrConstraintStringLengthValidation.F(
			MinLengthConstraintType.String(),
			strconv.FormatUint(uint64(c.value), 10),
		))
	}
}

func (c MinLength) ASTNode() schema.RuleASTNode {
	return newRuleASTNode(
		schema.TokenTypeNumber,
		strconv.FormatUint(uint64(c.value), 10),
		schema.RuleASTNodeSourceManual,
	)
}

func (c MinLength) Value() uint {
	return c.value
}

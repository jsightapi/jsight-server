package constraint

import (
	"strconv"

	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/bytes"
	"github.com/jsightapi/jsight-schema-core/errs"
	"github.com/jsightapi/jsight-schema-core/json"
)

type MaxLength struct {
	value uint
}

var (
	_ Constraint       = MaxLength{}
	_ Constraint       = (*MaxLength)(nil)
	_ LiteralValidator = MaxLength{}
	_ LiteralValidator = (*MaxLength)(nil)
)

func NewMaxLength(ruleValue bytes.Bytes) *MaxLength {
	return &MaxLength{
		value: parseUint(ruleValue, MaxLengthConstraintType),
	}
}

func (MaxLength) IsJsonTypeCompatible(t json.Type) bool {
	return t == json.TypeString
}

func (MaxLength) Type() Type {
	return MaxLengthConstraintType
}

func (c MaxLength) String() string {
	return MaxLengthConstraintType.String() + ": " + strconv.FormatUint(uint64(c.value), 10)
}

func (c MaxLength) Validate(value bytes.Bytes) {
	length := uint(value.Unquote().Len())
	if length > c.value {
		panic(errs.ErrConstraintStringLengthValidation.F(
			MaxLengthConstraintType.String(),
			strconv.FormatUint(uint64(c.value), 10),
		))
	}
}

func (c MaxLength) ASTNode() schema.RuleASTNode {
	return newRuleASTNode(
		schema.TokenTypeNumber,
		strconv.FormatUint(uint64(c.value), 10),
		schema.RuleASTNodeSourceManual,
	)
}

func (c MaxLength) Value() uint {
	return c.value
}

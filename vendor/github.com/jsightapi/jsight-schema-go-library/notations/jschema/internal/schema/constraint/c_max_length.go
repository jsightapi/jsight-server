package constraint

import (
	"fmt"
	"strconv"

	jschema "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/bytes"
	"github.com/jsightapi/jsight-schema-go-library/errors"
	"github.com/jsightapi/jsight-schema-go-library/internal/json"
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
	return fmt.Sprintf("%s: %d", MaxLengthConstraintType, c.value)
}

func (c MaxLength) Validate(value bytes.Bytes) {
	length := uint(len(value.Unquote()))
	if length > c.value {
		panic(errors.Format(
			errors.ErrConstraintStringLengthValidation,
			MaxLengthConstraintType.String(),
			strconv.FormatUint(uint64(c.value), 10),
		))
	}
}

func (c MaxLength) ASTNode() jschema.RuleASTNode {
	return newRuleASTNode(
		jschema.JSONTypeNumber,
		strconv.FormatUint(uint64(c.value), 10),
		jschema.RuleASTNodeSourceManual,
	)
}

func (c MaxLength) Value() uint {
	return c.value
}

package constraint

import (
	"fmt"
	"strconv"

	jschema "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/bytes"
	"github.com/jsightapi/jsight-schema-go-library/errors"
	"github.com/jsightapi/jsight-schema-go-library/internal/json"
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
	return fmt.Sprintf("%s: %d", MinLengthConstraintType, c.value)
}

func (c MinLength) Validate(value bytes.Bytes) {
	length := uint(len(value.Unquote()))
	if length < c.value {
		panic(errors.Format(
			errors.ErrConstraintStringLengthValidation,
			MinLengthConstraintType.String(),
			strconv.FormatUint(uint64(c.value), 10),
		))
	}
}

func (c MinLength) ASTNode() jschema.RuleASTNode {
	return newRuleASTNode(
		jschema.JSONTypeNumber,
		strconv.FormatUint(uint64(c.value), 10),
		jschema.RuleASTNodeSourceManual,
	)
}

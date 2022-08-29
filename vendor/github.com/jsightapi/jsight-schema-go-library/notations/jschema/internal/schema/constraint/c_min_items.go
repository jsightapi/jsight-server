package constraint

import (
	"strconv"

	jschema "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/bytes"
	"github.com/jsightapi/jsight-schema-go-library/errors"
	"github.com/jsightapi/jsight-schema-go-library/internal/json"
)

type MinItems struct {
	value uint
}

var (
	_ Constraint     = MinItems{}
	_ Constraint     = (*MinItems)(nil)
	_ ArrayValidator = MinItems{}
	_ ArrayValidator = (*MinItems)(nil)
)

func NewMinItems(ruleValue bytes.Bytes) *MinItems {
	return &MinItems{
		value: parseUint(ruleValue, MinItemsConstraintType),
	}
}

func (MinItems) IsJsonTypeCompatible(t json.Type) bool {
	return t == json.TypeArray
}

func (MinItems) Type() Type {
	return MinItemsConstraintType
}

func (c MinItems) String() string {
	return MinItemsConstraintType.String() + ": " + strconv.FormatUint(uint64(c.value), 10)
}

func (c MinItems) ValidateTheArray(numberOfChildren uint) {
	if numberOfChildren < c.value {
		panic(errors.ErrConstraintMinItemsValidation)
	}
}

func (c MinItems) Value() uint {
	return c.value
}

func (c MinItems) ASTNode() jschema.RuleASTNode {
	return newRuleASTNode(
		jschema.TokenTypeNumber,
		strconv.FormatUint(uint64(c.value), 10),
		jschema.RuleASTNodeSourceManual,
	)
}

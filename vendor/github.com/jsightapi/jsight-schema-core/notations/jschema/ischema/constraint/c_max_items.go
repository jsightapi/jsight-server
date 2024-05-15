package constraint

import (
	"strconv"

	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/bytes"
	"github.com/jsightapi/jsight-schema-core/errs"
	"github.com/jsightapi/jsight-schema-core/json"
)

type MaxItems struct {
	value uint
}

var (
	_ Constraint     = MaxItems{}
	_ Constraint     = (*MaxItems)(nil)
	_ ArrayValidator = MaxItems{}
	_ ArrayValidator = (*MaxItems)(nil)
)

func NewMaxItems(ruleValue bytes.Bytes) *MaxItems {
	return &MaxItems{
		value: parseUint(ruleValue, MaxItemsConstraintType),
	}
}

func (MaxItems) IsJsonTypeCompatible(t json.Type) bool {
	return t == json.TypeArray
}

func (MaxItems) Type() Type {
	return MaxItemsConstraintType
}

func (c MaxItems) String() string {
	return MaxItemsConstraintType.String() + ": " + strconv.FormatUint(uint64(c.value), 10)
}

func (c MaxItems) ValidateTheArray(numberOfChildren uint) {
	if numberOfChildren > c.value {
		panic(errs.ErrConstraintMaxItemsValidation.F())
	}
}

func (c MaxItems) Value() uint {
	return c.value
}

func (c MaxItems) ASTNode() schema.RuleASTNode {
	return newRuleASTNode(
		schema.TokenTypeNumber,
		strconv.FormatUint(uint64(c.value), 10),
		schema.RuleASTNodeSourceManual,
	)
}

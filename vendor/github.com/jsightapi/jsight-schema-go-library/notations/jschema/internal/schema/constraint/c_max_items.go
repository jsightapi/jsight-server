package constraint

import (
	jschema "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/bytes"
	"github.com/jsightapi/jsight-schema-go-library/errors"
	"github.com/jsightapi/jsight-schema-go-library/internal/json"
)

type MaxItems struct {
	value    *json.Number
	rawValue bytes.Bytes
}

var (
	_ Constraint     = MaxItems{}
	_ Constraint     = (*MaxItems)(nil)
	_ ArrayValidator = MaxItems{}
	_ ArrayValidator = (*MaxItems)(nil)
)

func NewMaxItems(ruleValue bytes.Bytes) *MaxItems {
	number, err := json.NewIntegerNumber(ruleValue)
	if err != nil {
		panic(err)
	}

	return &MaxItems{
		rawValue: ruleValue,
		value:    number,
	}
}

func (MaxItems) IsJsonTypeCompatible(t json.Type) bool {
	return t == json.TypeArray
}

func (MaxItems) Type() Type {
	return MaxItemsConstraintType
}

func (c MaxItems) String() string {
	return MaxItemsConstraintType.String() + ": " + c.value.String()
}

func (c MaxItems) ValidateTheArray(numberOfChildren uint) {
	length := json.NewNumberFromUint(numberOfChildren)
	if length.GreaterThan(c.value) {
		panic(errors.ErrConstraintMaxItemsValidation)
	}
}

func (c MaxItems) Value() *json.Number {
	return c.value
}

func (c MaxItems) ASTNode() jschema.RuleASTNode {
	return newRuleASTNode(jschema.JSONTypeNumber, c.rawValue.String(), jschema.RuleASTNodeSourceManual)
}

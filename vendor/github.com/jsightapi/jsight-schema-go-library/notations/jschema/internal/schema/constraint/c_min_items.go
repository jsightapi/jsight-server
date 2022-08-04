package constraint

import (
	jschema "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/bytes"
	"github.com/jsightapi/jsight-schema-go-library/errors"
	"github.com/jsightapi/jsight-schema-go-library/internal/json"
)

type MinItems struct {
	value    *json.Number
	rawValue bytes.Bytes
}

var (
	_ Constraint     = MinItems{}
	_ Constraint     = (*MinItems)(nil)
	_ ArrayValidator = MinItems{}
	_ ArrayValidator = (*MinItems)(nil)
)

func NewMinItems(ruleValue bytes.Bytes) *MinItems {
	number, err := json.NewIntegerNumber(ruleValue)
	if err != nil {
		panic(err)
	}

	return &MinItems{
		rawValue: ruleValue,
		value:    number,
	}
}

func (MinItems) IsJsonTypeCompatible(t json.Type) bool {
	return t == json.TypeArray
}

func (MinItems) Type() Type {
	return MinItemsConstraintType
}

func (c MinItems) String() string {
	return MinItemsConstraintType.String() + ": " + c.value.String()
}

func (c MinItems) ValidateTheArray(numberOfChildren uint) {
	length := json.NewNumberFromUint(numberOfChildren)
	if length.LessThan(c.value) {
		panic(errors.ErrConstraintMinItemsValidation)
	}
}

func (c MinItems) Value() *json.Number {
	return c.value
}

func (c MinItems) ASTNode() jschema.RuleASTNode {
	return newRuleASTNode(jschema.JSONTypeNumber, c.rawValue.String(), jschema.RuleASTNodeSourceManual)
}

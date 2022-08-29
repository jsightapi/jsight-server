package constraint

import (
	"strconv"

	jschema "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/bytes"
	"github.com/jsightapi/jsight-schema-go-library/errors"
	"github.com/jsightapi/jsight-schema-go-library/internal/json"
)

type Const struct {
	nodeValue bytes.Bytes
	apply     bool
}

var (
	_ Constraint       = Const{}
	_ Constraint       = (*Const)(nil)
	_ BoolKeeper       = Const{}
	_ BoolKeeper       = (*Const)(nil)
	_ LiteralValidator = Const{}
	_ LiteralValidator = (*Const)(nil)
)

func NewConst(value, nodeValue bytes.Bytes) *Const {
	c := Const{
		nodeValue: nodeValue,
	}

	var err error
	if c.apply, err = value.ParseBool(); err != nil {
		panic(errors.Format(errors.ErrInvalidValueOfConstraint, ConstConstraintType.String()))
	}
	return &c
}

func (Const) IsJsonTypeCompatible(t json.Type) bool {
	return t != json.TypeObject && t != json.TypeArray
}

func (Const) Type() Type {
	return ConstConstraintType
}

func (c Const) String() string {
	if c.apply {
		return ConstConstraintType.String() + ": true"
	}
	return ConstConstraintType.String() + ": false"
}

func (c Const) Bool() bool {
	return c.apply
}

func (c Const) Validate(v bytes.Bytes) {
	if !c.apply {
		return
	}

	if v.String() != c.nodeValue.String() {
		panic(errors.Format(errors.ErrInvalidConst, c.nodeValue.String()))
	}
}

func (c Const) ASTNode() jschema.RuleASTNode {
	return newRuleASTNode(jschema.TokenTypeBoolean, strconv.FormatBool(c.apply), jschema.RuleASTNodeSourceManual)
}

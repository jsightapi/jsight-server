package constraint

import (
	"strconv"

	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/bytes"
	"github.com/jsightapi/jsight-schema-core/errs"
	"github.com/jsightapi/jsight-schema-core/json"
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
		panic(errs.ErrInvalidValueOfConstraint.F(ConstConstraintType.String()))
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
		return ConstConstraintType.String() + colonTrue
	}
	return ConstConstraintType.String() + colonFalse
}

func (c Const) Bool() bool {
	return c.apply
}

func (c Const) Validate(v bytes.Bytes) {
	if !c.apply {
		return
	}

	if v.Unquote().String() != c.nodeValue.Unquote().String() {
		panic(errs.ErrInvalidConst.F(c.nodeValue.String()))
	}
}

func (c Const) ASTNode() schema.RuleASTNode {
	return newRuleASTNode(schema.TokenTypeBoolean, strconv.FormatBool(c.apply), schema.RuleASTNodeSourceManual)
}

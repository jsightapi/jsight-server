package constraint

import (
	"time"

	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/bytes"
	"github.com/jsightapi/jsight-schema-core/errs"
	"github.com/jsightapi/jsight-schema-core/json"
)

type Date struct{}

var (
	_ Constraint       = Date{}
	_ Constraint       = (*Date)(nil)
	_ LiteralValidator = Date{}
	_ LiteralValidator = (*Date)(nil)
)

func NewDate() *Date {
	return &Date{}
}

func (Date) IsJsonTypeCompatible(t json.Type) bool {
	return t == json.TypeString
}

func (Date) Type() Type {
	return DateConstraintType
}

func (Date) String() string {
	return DateConstraintType.String()
}

func (Date) Validate(value bytes.Bytes) {
	str := value.Unquote().String()
	_, err := time.Parse("2006-01-02", str)
	if err != nil {
		panic(errs.ErrInvalidDate.F(err))
	}
}

func (Date) ASTNode() schema.RuleASTNode {
	return newEmptyRuleASTNode()
}

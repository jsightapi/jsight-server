package constraint

import (
	"time"

	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/bytes"
	"github.com/jsightapi/jsight-schema-core/errs"
	"github.com/jsightapi/jsight-schema-core/json"
)

type DateTime struct{}

var (
	_ Constraint       = DateTime{}
	_ Constraint       = (*DateTime)(nil)
	_ LiteralValidator = DateTime{}
	_ LiteralValidator = (*DateTime)(nil)
)

func NewDateTime() *DateTime {
	return &DateTime{}
}

func (DateTime) IsJsonTypeCompatible(t json.Type) bool {
	return t == json.TypeString
}

func (DateTime) Type() Type {
	return DateTimeConstraintType
}

func (DateTime) String() string {
	return DateTimeConstraintType.String()
}

func (DateTime) Validate(value bytes.Bytes) {
	str := value.Unquote().String()
	_, err := time.Parse(time.RFC3339, str)
	if err != nil {
		panic(errs.ErrInvalidDateTime.F())
	}
}

func (DateTime) ASTNode() schema.RuleASTNode {
	return newEmptyRuleASTNode()
}

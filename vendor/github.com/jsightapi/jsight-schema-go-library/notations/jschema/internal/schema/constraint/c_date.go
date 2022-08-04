package constraint

import (
	"time"

	jschema "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/bytes"
	"github.com/jsightapi/jsight-schema-go-library/errors"
	"github.com/jsightapi/jsight-schema-go-library/internal/json"
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
		panic(errors.Format(errors.ErrInvalidDate, err))
	}
}

func (Date) ASTNode() jschema.RuleASTNode {
	return newEmptyRuleASTNode()
}

package constraint

import (
	"time"

	jschema "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/bytes"
	"github.com/jsightapi/jsight-schema-go-library/errors"
	"github.com/jsightapi/jsight-schema-go-library/internal/json"
)

type DateTime struct{}

var _ Constraint = DateTime{}

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
		panic(errors.ErrInvalidDateTime)
	}
}

func (DateTime) ASTNode() jschema.RuleASTNode {
	return newEmptyRuleASTNode()
}

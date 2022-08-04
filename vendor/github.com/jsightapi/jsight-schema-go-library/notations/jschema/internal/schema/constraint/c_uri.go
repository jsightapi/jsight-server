package constraint

import (
	"net/url"

	jschema "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/bytes"
	"github.com/jsightapi/jsight-schema-go-library/errors"
	"github.com/jsightapi/jsight-schema-go-library/internal/json"
)

type Uri struct{}

var (
	_ Constraint       = Uri{}
	_ Constraint       = (*Uri)(nil)
	_ LiteralValidator = Uri{}
	_ LiteralValidator = (*Uri)(nil)
)

func NewUri() *Uri {
	return &Uri{}
}

func (Uri) IsJsonTypeCompatible(t json.Type) bool {
	return t == json.TypeString
}

func (Uri) Type() Type {
	return UriConstraintType
}

func (Uri) String() string {
	return UriConstraintType.String()
}

func (Uri) Validate(value bytes.Bytes) {
	val := value.Unquote().String()
	u, err := url.ParseRequestURI(val)
	if err != nil || !u.IsAbs() || u.Hostname() == "" {
		panic(errors.Format(errors.ErrInvalidUri, val))
	}
}

func (Uri) ASTNode() jschema.RuleASTNode {
	return newEmptyRuleASTNode()
}

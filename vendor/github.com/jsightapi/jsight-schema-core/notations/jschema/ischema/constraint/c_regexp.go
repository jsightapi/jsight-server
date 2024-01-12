package constraint

import (
	"encoding/json"
	"regexp"

	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/bytes"
	"github.com/jsightapi/jsight-schema-core/errs"
	internalJSON "github.com/jsightapi/jsight-schema-core/json"
)

type Regex struct {
	re         *regexp.Regexp
	expression string
}

var (
	_ Constraint       = Regex{}
	_ Constraint       = (*Regex)(nil)
	_ LiteralValidator = Regex{}
	_ LiteralValidator = (*Regex)(nil)
)

func NewRegex(value bytes.Bytes) *Regex {
	var str string // decoded json string. JSON "aaa\\bbb" to string "aaa\bbb".
	err := json.Unmarshal(value.Data(), &str)
	if err != nil {
		panic(err)
	}

	return &Regex{
		expression: str,
		re:         regexp.MustCompile(str), // can panic
	}
}

func (Regex) IsJsonTypeCompatible(t internalJSON.Type) bool {
	return t == internalJSON.TypeString
}

func (Regex) Type() Type {
	return RegexConstraintType
}

func (c Regex) String() string {
	return RegexConstraintType.String() + ": " + c.expression
}

func (c Regex) Validate(value bytes.Bytes) {
	if !c.re.Match(value.Unquote().Data()) {
		panic(errs.ErrDoesNotMatchRegularExpression.F())
	}
}

func (c Regex) ASTNode() schema.RuleASTNode {
	return newRuleASTNode(schema.TokenTypeString, c.expression, schema.RuleASTNodeSourceManual)
}

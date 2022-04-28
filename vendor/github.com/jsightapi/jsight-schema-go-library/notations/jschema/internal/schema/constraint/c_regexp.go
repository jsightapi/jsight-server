package constraint

import (
	"encoding/json"
	"regexp"

	jschema "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/bytes"
	"github.com/jsightapi/jsight-schema-go-library/errors"
	internalJSON "github.com/jsightapi/jsight-schema-go-library/internal/json"
)

type Regex struct {
	re         *regexp.Regexp
	expression string
}

var _ Constraint = Regex{}

func NewRegex(value bytes.Bytes) *Regex {
	var str string // decoded json string. JSON "aaa\\bbb" to string "aaa\bbb".
	err := json.Unmarshal(value, &str)
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
	if !c.re.Match(value.Unquote()) {
		panic(errors.ErrDoesNotMatchRegularExpression)
	}
}

func (c Regex) ASTNode() jschema.RuleASTNode {
	return newRuleASTNode(jschema.JSONTypeString, c.expression, jschema.RuleASTNodeSourceManual)
}

package constraint

import (
	"bytes"
	"strings"

	jschema "github.com/jsightapi/jsight-schema-go-library"
	jbytes "github.com/jsightapi/jsight-schema-go-library/bytes"
	"github.com/jsightapi/jsight-schema-go-library/errors"
	"github.com/jsightapi/jsight-schema-go-library/internal/json"
)

type Enum struct {
	uniqueIdx map[string]struct{}
	ruleName  string
	items     []enumItem
}

type enumItem struct {
	comment string
	value   jbytes.Bytes
}

var _ Constraint = Enum{}

func NewEnum() *Enum {
	return &Enum{
		uniqueIdx: make(map[string]struct{}),
		items:     make([]enumItem, 0, 5),
	}
}

func (Enum) IsJsonTypeCompatible(t json.Type) bool {
	return t.IsLiteralType()
}

func (Enum) Type() Type {
	return EnumConstraintType
}

func (c Enum) String() string {
	var str strings.Builder
	str.WriteString(EnumConstraintType.String())
	str.WriteString(": [")
	for i, v := range c.items {
		str.WriteString(v.value.String())
		if len(c.items)-1 != i {
			str.WriteString(", ")
		}
	}
	str.WriteString("]")
	return str.String()
}

func (c *Enum) Append(b jbytes.Bytes) int {
	key := b.TrimSpaces().String()
	if _, ok := c.uniqueIdx[key]; ok {
		panic(errors.Format(errors.ErrDuplicationInEnumRule, b.String()))
	}
	idx := len(c.items)
	c.items = append(c.items, enumItem{value: b})
	c.uniqueIdx[key] = struct{}{}
	return idx
}

func (c *Enum) SetComment(idx int, comment string) {
	c.items[idx].comment = comment
}

func (c *Enum) SetRuleName(s string) {
	c.ruleName = s
}

func (c *Enum) RuleName() string {
	return c.ruleName
}

func (c Enum) Validate(a jbytes.Bytes) {
	for _, b := range c.items {
		if bytes.Equal(a, b.value) {
			return
		}
	}
	panic(errors.Format(errors.ErrDoesNotMatchAnyOfTheEnumValues))
}

func (c Enum) ASTNode() jschema.RuleASTNode {
	const source = jschema.RuleASTNodeSourceManual

	if c.ruleName != "" {
		return newRuleASTNode(jschema.JSONTypeShortcut, c.ruleName, source)
	}

	n := newRuleASTNode(jschema.JSONTypeArray, "", source)
	n.Items = make([]jschema.RuleASTNode, 0, len(c.items))

	for _, b := range c.items {
		an := newRuleASTNode(
			json.Guess(b.value).JsonType().ToTokenType(),
			b.value.Unquote().String(),
			source,
		)
		an.Comment = b.comment

		n.Items = append(n.Items, an)
	}

	return n
}

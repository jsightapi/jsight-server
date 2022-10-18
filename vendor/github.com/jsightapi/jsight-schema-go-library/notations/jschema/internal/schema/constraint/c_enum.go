package constraint

import (
	"encoding/json"
	"strings"

	jschema "github.com/jsightapi/jsight-schema-go-library"
	jbytes "github.com/jsightapi/jsight-schema-go-library/bytes"
	"github.com/jsightapi/jsight-schema-go-library/errors"
	jjson "github.com/jsightapi/jsight-schema-go-library/internal/json"
)

type Enum struct {
	uniqueIdx map[enumItemValue]struct{}
	ruleName  string
	items     []EnumItem
}

type EnumItem struct {
	src     jbytes.Bytes
	comment string
	enumItemValue
}

type enumItemValue struct {
	value    string
	jsonType jjson.Type
}

func (v enumItemValue) String() string {
	if v.jsonType == jjson.TypeString {
		b, err := json.Marshal(v.value)
		if err != nil {
			panic(errors.ErrImpossible)
		}
		return string(b)
	} else {
		return v.value
	}
}

func NewEnumItem(b jbytes.Bytes, c string) EnumItem {
	i := EnumItem{src: b, comment: c}
	b = b.TrimSpaces()
	i.jsonType = jjson.Guess(b).JsonType()
	if i.jsonType == jjson.TypeString {
		b = b.Unquote()
	}
	i.value = b.String()
	return i
}

var (
	_ Constraint       = Enum{}
	_ Constraint       = (*Enum)(nil)
	_ LiteralValidator = Enum{}
	_ LiteralValidator = (*Enum)(nil)
)

func NewEnum() *Enum {
	return &Enum{
		uniqueIdx: make(map[enumItemValue]struct{}),
		items:     make([]EnumItem, 0, 5),
	}
}

func (Enum) IsJsonTypeCompatible(t jjson.Type) bool {
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
		str.WriteString(v.enumItemValue.String())
		if len(c.items)-1 != i {
			str.WriteString(", ")
		}
	}
	str.WriteString("]")
	return str.String()
}

func (c *Enum) Append(i EnumItem) int {
	if _, ok := c.uniqueIdx[i.enumItemValue]; ok {
		panic(errors.Format(errors.ErrDuplicationInEnumRule, i.src.String()))
	}
	idx := len(c.items)
	c.items = append(c.items, i)
	c.uniqueIdx[i.enumItemValue] = struct{}{}
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
	aa := NewEnumItem(a, "")
	for _, b := range c.items {
		if aa.enumItemValue == b.enumItemValue {
			return
		}
	}
	panic(errors.ErrDoesNotMatchAnyOfTheEnumValues)
}

func (c Enum) ASTNode() jschema.RuleASTNode {
	const source = jschema.RuleASTNodeSourceManual

	if c.ruleName != "" {
		return newRuleASTNode(jschema.TokenTypeShortcut, c.ruleName, source)
	}

	n := newRuleASTNode(jschema.TokenTypeArray, "", source)
	n.Items = make([]jschema.RuleASTNode, 0, len(c.items))

	for _, b := range c.items {
		an := newRuleASTNode(
			b.jsonType.ToTokenType(),
			b.value,
			source,
		)
		an.Comment = b.comment

		n.Items = append(n.Items, an)
	}

	return n
}

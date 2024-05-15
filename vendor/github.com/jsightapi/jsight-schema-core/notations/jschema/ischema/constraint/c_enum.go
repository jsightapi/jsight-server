package constraint

import (
	stdJson "encoding/json"
	"strings"

	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/bytes"
	"github.com/jsightapi/jsight-schema-core/errs"
	"github.com/jsightapi/jsight-schema-core/json"
)

type Enum struct {
	uniqueIdx map[enumItemValue]struct{}
	ruleName  string
	items     []EnumItem
}

type EnumItem struct {
	comment string
	enumItemValue
	src bytes.Bytes
}

type enumItemValue struct {
	value    string
	jsonType json.Type
}

func (v enumItemValue) String() string {
	if v.jsonType == json.TypeString {
		b, err := stdJson.Marshal(v.value)
		if err != nil {
			panic(errs.ErrRuntimeFailure.F())
		}
		return string(b)
	} else {
		return v.value
	}
}

func NewEnumItem(b bytes.Bytes, c string) EnumItem {
	i := EnumItem{src: b, comment: c}
	b = b.TrimSpaces()
	i.jsonType = json.Guess(b).JsonType()
	if i.jsonType == json.TypeString {
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
		panic(errs.ErrDuplicationInEnumRule.F(i.src.String()))
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

func (c Enum) Validate(a bytes.Bytes) {
	aa := NewEnumItem(a, "")
	for _, b := range c.items {
		if aa.enumItemValue == b.enumItemValue {
			return
		}
	}
	panic(errs.ErrDoesNotMatchAnyOfTheEnumValues.F())
}

func (c Enum) ASTNode() schema.RuleASTNode {
	const source = schema.RuleASTNodeSourceManual

	if c.ruleName != "" {
		return newRuleASTNode(schema.TokenTypeShortcut, c.ruleName, source)
	}

	n := newRuleASTNode(schema.TokenTypeArray, "", source)
	n.Items = make([]schema.RuleASTNode, 0, len(c.items))

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

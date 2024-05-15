package constraint

import (
	"strings"

	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/bytes"
	"github.com/jsightapi/jsight-schema-core/errs"
	"github.com/jsightapi/jsight-schema-core/json"
)

type AllOf struct {
	schemaName []string
}

var (
	_ Constraint = AllOf{}
	_ Constraint = (*AllOf)(nil)
)

func NewAllOf() *AllOf {
	return &AllOf{
		schemaName: make([]string, 0, 3),
	}
}

func (AllOf) IsJsonTypeCompatible(t json.Type) bool {
	return t == json.TypeObject
}

func (AllOf) Type() Type {
	return AllOfConstraintType
}

func (c AllOf) String() string {
	return AllOfConstraintType.String() + ": " + strings.Join(c.schemaName, ", ")
}

func (c *AllOf) Append(scalar bytes.Bytes) {
	if !json.Guess(scalar).IsString() {
		panic(errs.ErrUnacceptableValueInAllOfRule.F())
	}

	s := scalar.Unquote()

	if !s.IsUserTypeName() {
		panic(errs.ErrInvalidSchemaNameInAllOfRule.F(s))
	}
	c.schemaName = append(c.schemaName, s.String())
}

func (c AllOf) SchemaNames() []string {
	return c.schemaName
}

func (c AllOf) ASTNode() schema.RuleASTNode {
	const source = schema.RuleASTNodeSourceManual

	if len(c.schemaName) == 1 {
		return newRuleASTNode(schema.TokenTypeShortcut, c.schemaName[0], source)
	}

	n := newRuleASTNode(schema.TokenTypeArray, "", source)
	n.Items = make([]schema.RuleASTNode, 0, len(c.schemaName))

	for _, sn := range c.schemaName {
		n.Items = append(n.Items, newRuleASTNode(schema.TokenTypeShortcut, sn, source))
	}

	return n
}

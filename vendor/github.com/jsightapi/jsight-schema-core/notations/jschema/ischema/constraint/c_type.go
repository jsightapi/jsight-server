package constraint

import (
	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/bytes"
	"github.com/jsightapi/jsight-schema-core/json"
)

type TypeConstraint struct {
	value  bytes.Bytes
	source schema.RuleASTNodeSource
}

var (
	_ Constraint  = TypeConstraint{}
	_ Constraint  = (*TypeConstraint)(nil)
	_ BytesKeeper = TypeConstraint{}
	_ BytesKeeper = (*TypeConstraint)(nil)
)

func NewType(ruleValue bytes.Bytes, source schema.RuleASTNodeSource) *TypeConstraint {
	return &TypeConstraint{
		value:  ruleValue,
		source: source,
	}
}

func (c TypeConstraint) IsGenerated() bool {
	return c.source == schema.RuleASTNodeSourceGenerated
}

func (TypeConstraint) IsJsonTypeCompatible(json.Type) bool {
	return true
}

func (TypeConstraint) Type() Type {
	return TypeConstraintType
}

func (c TypeConstraint) String() string {
	return TypeConstraintType.String() + ": " + c.value.String()
}

func (c TypeConstraint) Bytes() bytes.Bytes {
	return c.value
}

func (c TypeConstraint) ASTNode() schema.RuleASTNode {
	t := schema.TokenTypeString
	if c.value.Unquote().IsUserTypeName() {
		t = schema.TokenTypeShortcut
	}
	return newRuleASTNode(t, c.value.Unquote().String(), c.source)
}

func (c TypeConstraint) Source() schema.RuleASTNodeSource { return c.source }

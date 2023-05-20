package constraint

import (
	"strings"

	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/bytes"
	"github.com/jsightapi/jsight-schema-core/errs"
	"github.com/jsightapi/jsight-schema-core/json"
)

type AdditionalPropertiesMode int

const (
	AdditionalPropertiesCanBeAny AdditionalPropertiesMode = iota
	AdditionalPropertiesMustBeSchemaType
	AdditionalPropertiesMustBeUserType
	AdditionalPropertiesNotAllowed
)

type AdditionalProperties struct {
	schemaType schema.SchemaType // only for AdditionalPropertiesMustBeSchemaType
	typeName   bytes.Bytes       // only for AdditionalPropertiesMustBeUserType
	astNode    schema.RuleASTNode
	mode       AdditionalPropertiesMode
}

var (
	_ Constraint = AdditionalProperties{}
	_ Constraint = (*AdditionalProperties)(nil)
)

// NewAdditionalProperties create an additional properties constraint.
// Depends on `ruleValue` value, might return nil.
// Might panic if got unknown JSON type.
//
// Handle next cases:
//
//	{additionalProperties: "any"}
//	{additionalProperties: true}
//	{additionalProperties: false} - in that case this function will return nil.
//	{additionalProperties: "@Foo"}
//	{additionalProperties: "string"}
func NewAdditionalProperties(ruleValue bytes.Bytes) *AdditionalProperties {
	c := &AdditionalProperties{}

	c.astNode = newEmptyRuleASTNode()
	c.astNode.Source = schema.RuleASTNodeSourceManual

	txt := ruleValue.Unquote()
	txtStr := txt.String()
	switch {
	case txt.OneOf("any", "true"):
		if txt.String() == "true" {
			c.astNode.TokenType = schema.TokenTypeBoolean
			c.astNode.Value = "true"
		} else {
			c.astNode.TokenType = schema.TokenTypeString
			c.astNode.Value = txtStr
		}
		c.mode = AdditionalPropertiesCanBeAny

	case txt.String() == "false":
		c.astNode.TokenType = schema.TokenTypeBoolean
		c.astNode.Value = "false"
		c.mode = AdditionalPropertiesNotAllowed

	case txt.IsUserTypeName():
		c.astNode.TokenType = schema.TokenTypeString
		c.astNode.Value = txtStr
		c.mode = AdditionalPropertiesMustBeUserType
		c.typeName = txt

	case schema.IsValidType(txtStr):
		c.astNode.TokenType = schema.TokenTypeString
		c.astNode.Value = txtStr
		c.mode = AdditionalPropertiesMustBeSchemaType
		c.schemaType = schema.SchemaType(txtStr)

	default:
		panic(errs.ErrUnknownJSchemaType.F(txtStr))
	}
	return c
}

func (AdditionalProperties) IsJsonTypeCompatible(t json.Type) bool {
	return t == json.TypeObject
}

func (AdditionalProperties) Type() Type {
	return AdditionalPropertiesConstraintType
}

func (c AdditionalProperties) String() string {
	buf := strings.Builder{}
	buf.WriteString(AdditionalPropertiesConstraintType.String() + ": ")

	switch c.mode {
	case AdditionalPropertiesCanBeAny:
		buf.WriteString("any")

	case AdditionalPropertiesMustBeSchemaType:
		buf.WriteString(string(c.schemaType))

	case AdditionalPropertiesMustBeUserType:
		buf.WriteString(c.typeName.String())

	case AdditionalPropertiesNotAllowed:
		buf.WriteString("false")

	default:
		panic(errs.ErrRuntimeFailure.F())
	}

	return buf.String()
}

func (c AdditionalProperties) Mode() AdditionalPropertiesMode {
	return c.mode
}

func (c AdditionalProperties) SchemaType() schema.SchemaType {
	return c.schemaType
}

func (c AdditionalProperties) TypeName() bytes.Bytes {
	return c.typeName
}

func (c AdditionalProperties) IsEqual(c2 AdditionalProperties) bool {
	return c.schemaType == c2.schemaType && c.typeName.String() == c2.typeName.String()
}

func (c AdditionalProperties) ASTNode() schema.RuleASTNode {
	return c.astNode
}

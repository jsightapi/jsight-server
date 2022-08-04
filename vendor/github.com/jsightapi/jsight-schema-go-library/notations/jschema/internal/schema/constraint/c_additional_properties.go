package constraint

import (
	"strings"

	jschema "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/bytes"
	"github.com/jsightapi/jsight-schema-go-library/errors"
	"github.com/jsightapi/jsight-schema-go-library/internal/json"
)

type AdditionalPropertiesMode int

const (
	AdditionalPropertiesCanBeAny AdditionalPropertiesMode = iota
	AdditionalPropertiesMustBeSchemaType
	AdditionalPropertiesMustBeType
	AdditionalPropertiesNotAllowed
)

type AdditionalProperties struct {
	schemaType jschema.SchemaType // only for AdditionalPropertiesMustBeSchemaType
	typeName   bytes.Bytes        // only for AdditionalPropertiesMustBeType
	astNode    jschema.RuleASTNode
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
//  {additionalProperties: "any"}
//  {additionalProperties: true}
//  {additionalProperties: false} - in that case this function will return nil.
//  {additionalProperties: "@Foo"}
//  {additionalProperties: "string"}
func NewAdditionalProperties(ruleValue bytes.Bytes) *AdditionalProperties {
	c := &AdditionalProperties{}

	c.astNode = newEmptyRuleASTNode()
	c.astNode.Source = jschema.RuleASTNodeSourceManual

	txt := ruleValue.Unquote()
	txtStr := txt.String()
	switch {
	case txt.OneOf("any", "true"):
		if txt.String() == "true" {
			c.astNode.JSONType = jschema.JSONTypeBoolean
			c.astNode.Value = "true"
		} else {
			c.astNode.JSONType = jschema.JSONTypeString
			c.astNode.Value = txtStr
		}
		c.mode = AdditionalPropertiesCanBeAny

	case txt.String() == "false":
		c.astNode.JSONType = jschema.JSONTypeBoolean
		c.astNode.Value = "false"
		c.mode = AdditionalPropertiesNotAllowed

	case txt.IsUserTypeName():
		c.astNode.JSONType = jschema.JSONTypeString
		c.astNode.Value = txtStr
		c.mode = AdditionalPropertiesMustBeType
		c.typeName = txt

	case jschema.IsValidType(txtStr):
		c.astNode.JSONType = jschema.JSONTypeString
		c.astNode.Value = txtStr
		c.mode = AdditionalPropertiesMustBeSchemaType
		c.schemaType = jschema.SchemaType(txtStr)

	default:
		panic(errors.Format(errors.ErrUnknownJSchemaType, txtStr))
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

	case AdditionalPropertiesMustBeType:
		buf.WriteString(c.typeName.String())

	case AdditionalPropertiesNotAllowed:
		buf.WriteString("false")

	default:
		panic(errors.Format(errors.ErrGeneric, "Constraint error"))
	}

	return buf.String()
}

func (c AdditionalProperties) Mode() AdditionalPropertiesMode {
	return c.mode
}

func (c AdditionalProperties) SchemaType() jschema.SchemaType {
	return c.schemaType
}

func (c AdditionalProperties) TypeName() bytes.Bytes {
	return c.typeName
}

func (c AdditionalProperties) IsEqual(c2 AdditionalProperties) bool {
	return c.schemaType == c2.schemaType && c.typeName.String() == c2.typeName.String()
}

func (c AdditionalProperties) ASTNode() jschema.RuleASTNode {
	return c.astNode
}

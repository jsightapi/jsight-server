package internal

import (
	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/errs"
)

func RuleToASTNode(r schema.RuleASTNode) schema.ASTNode {
	switch r.TokenType {
	case schema.TokenTypeString, schema.TokenTypeShortcut:
		return stringRuleToASTNode(r)
	case schema.TokenTypeObject:
		return objectRuleToASTNode(r)
	default:
		panic(errs.ErrRuntimeFailure.F())
	}
}

// stringRuleToASTNode returns the ASTNode for "OR" rule elements. JSight example: // {or: [ "email", "integer" ]}
func stringRuleToASTNode(r schema.RuleASTNode) schema.ASTNode {
	a := schema.ASTNode{
		Rules: schema.MakeRuleASTNodes(1),
	}

	return stringRuleToASTNodeType(a, r.Value)
}

func stringRuleToASTNodeType(a schema.ASTNode, s string) schema.ASTNode {
	if s == "any" {
		a.TokenType = schema.TokenTypeString
		a.SchemaType = s // any
	} else if format := FormatFromSchemaType(s); format != nil { // JSight example: // {or: [ "email"... ]}
		a.TokenType = schema.TokenTypeString
		a.Rules.Set("type", schema.RuleASTNode{
			TokenType: schema.TokenTypeString,
			Value:     *format,
		})
	} else { // JSight example: // {or: [ "integer", "@cat" ]}
		a.TokenType = TokenType(s)
		a.SchemaType = s
		a.Value = s
	}

	return a
}

// objectRuleToASTNode returns the ASTNode for "OR" rule elements. JSight example: // {or: [ {type: "integer"}, ... ]}
func objectRuleToASTNode(r schema.RuleASTNode) schema.ASTNode {
	a := schema.ASTNode{
		Rules: schema.MakeRuleASTNodes(r.Properties.Len()),
	}

	if r.Properties != nil { // or: [ {...} ]
		_ = r.Properties.Each(func(k string, v schema.RuleASTNode) error {
			a.Rules.Set(k, v)
			return nil
		})
	}

	if typeRule, ok := r.Properties.Get("type"); ok { // or: [ { type: ...} ]
		a = stringRuleToASTNodeType(a, typeRule.Value)
	}

	return a
}

func strRef(s string) *string {
	return &s
}

func FormatFromSchemaType(s string) *string {
	switch s {
	case string(schema.SchemaTypeEmail), string(schema.SchemaTypeURI),
		string(schema.SchemaTypeUUID), string(schema.SchemaTypeDate):
		return strRef(s)
	case string(schema.SchemaTypeDateTime):
		return strRef("date-time")
	default:
		return nil
	}
}

func TokenType(s string) schema.TokenType {
	if s[0] == '@' {
		return schema.TokenTypeShortcut
	}

	switch s {
	case "string":
		return schema.TokenTypeString
	case "boolean":
		return schema.TokenTypeBoolean
	case "float", "integer", "decimal":
		return schema.TokenTypeNumber
	case "object":
		return schema.TokenTypeObject
	case "array":
		return schema.TokenTypeArray
	case "null":
		return schema.TokenTypeNull
	default:
		panic(errs.ErrRuntimeFailure.F())
	}
}

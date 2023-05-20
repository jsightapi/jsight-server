package constraint

import (
	schema "github.com/jsightapi/jsight-schema-core"
)

func newEmptyRuleASTNode() schema.RuleASTNode {
	return schema.RuleASTNode{
		Properties: &schema.RuleASTNodes{},
		Source:     schema.RuleASTNodeSourceManual,
	}
}

func newRuleASTNode(t schema.TokenType, v string, s schema.RuleASTNodeSource) schema.RuleASTNode {
	an := newEmptyRuleASTNode()

	an.TokenType = t
	an.Value = v
	an.Source = s

	return an
}

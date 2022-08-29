package constraint

import (
	jschema "github.com/jsightapi/jsight-schema-go-library"
)

func newEmptyRuleASTNode() jschema.RuleASTNode {
	return jschema.RuleASTNode{
		Properties: &jschema.RuleASTNodes{},
		Source:     jschema.RuleASTNodeSourceManual,
	}
}

func newRuleASTNode(t jschema.TokenType, v string, s jschema.RuleASTNodeSource) jschema.RuleASTNode {
	an := newEmptyRuleASTNode()

	an.TokenType = t
	an.Value = v
	an.Source = s

	return an
}

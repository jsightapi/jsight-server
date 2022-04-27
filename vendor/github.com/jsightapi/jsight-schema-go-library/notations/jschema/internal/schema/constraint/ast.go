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

func newRuleASTNode(t jschema.JSONType, v string, s jschema.RuleASTNodeSource) jschema.RuleASTNode {
	an := newEmptyRuleASTNode()

	an.JSONType = t
	an.Value = v
	an.Source = s

	return an
}

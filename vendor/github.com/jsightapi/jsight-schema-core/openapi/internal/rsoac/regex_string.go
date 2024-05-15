package rsoac

import (
	schema "github.com/jsightapi/jsight-schema-core"
)

type RegexString struct {
	Pattern *Pattern `json:"pattern"`
	Type    string   `json:"type"`
}

func newRegexString(astNode schema.ASTNode) *RegexString {
	var p = RegexString{
		Pattern: newPattern(astNode.Value),
		Type:    "string",
	}
	return &p
}

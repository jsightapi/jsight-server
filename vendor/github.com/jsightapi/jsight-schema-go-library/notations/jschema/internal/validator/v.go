package validator

import (
	"github.com/jsightapi/jsight-schema-go-library/internal/lexeme"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema/internal/schema"
)

type validator interface {
	parent() validator
	setParent(validator)

	// feed returns array (pointers to validators, or nil if not found), bool
	// (true if validator of node is completed), panic on error.
	feed(jsonLexeme lexeme.LexEvent) ([]validator, bool)

	// node returns this validator node.
	// For debug/log only.
	node() schema.Node
}

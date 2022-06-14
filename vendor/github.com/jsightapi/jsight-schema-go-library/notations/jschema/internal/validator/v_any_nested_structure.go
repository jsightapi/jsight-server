package validator

import (
	"fmt"

	"github.com/jsightapi/jsight-schema-go-library/internal/lexeme"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema/internal/schema"
)

// validator to process any nested structures.

type anyNestedStructure struct {
	node_   schema.Node
	parent_ validator
	depth   uint
}

func newAnyNestedStructureValidator(node schema.Node, parent validator) *anyNestedStructure {
	v := anyNestedStructure{
		node_:   node,
		parent_: parent,
	}
	return &v
}

func (v anyNestedStructure) node() schema.Node {
	return v.node_
}

func (v anyNestedStructure) parent() validator {
	return v.parent_
}

func (v *anyNestedStructure) setParent(parent validator) {
	v.parent_ = parent
}

// return nil (empty list pointers to validators) and bool (true if validator is done)
func (v *anyNestedStructure) feed(jsonLexeme lexeme.LexEvent) ([]validator, bool) {
	if jsonLexeme.Type().IsOpening() {
		v.depth++
	} else {
		v.depth--
	}

	if v.depth == 0 {
		return nil, true
	}

	return nil, false
}

func (v anyNestedStructure) log() string {
	return fmt.Sprintf("%s [%p]", v.node_.Type().String(), v.node_)
}

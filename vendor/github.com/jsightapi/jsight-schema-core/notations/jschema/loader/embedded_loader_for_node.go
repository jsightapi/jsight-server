package loader

import (
	"github.com/jsightapi/jsight-schema-core/lexeme"
	"github.com/jsightapi/jsight-schema-core/notations/jschema/ischema"
)

// nodeLoader is responsible for loading the JSON elements in the nodes of the
// internal representation of the SCHEMA.
type nodeLoader struct {
	// The schema.
	// Parameter from the main loader.
	schema *ischema.ISchema

	// nodesPerCurrentLineCount counts the number of nodes in a line. To check
	// because the rule cannot be added if there is more than one node suitable
	// for this in the row.
	// Parameter from the main loader.
	nodesPerCurrentLineCount *uint

	// leaf a schema node which is processed at the moment. A schema is a tree of
	// nodes. Starting from the root node, we fill the nodes with lexeme events.
	// If necessary, new nodes are added. Thus, the complete scheme tree is constructed.
	leaf ischema.Node
}

func newNodeLoader(
	schem *ischema.ISchema,
	nodesPerCurrentLineCount *uint,
) *nodeLoader {
	return &nodeLoader{
		schema:                   schem,
		nodesPerCurrentLineCount: nodesPerCurrentLineCount,
	}
}

// Load returns the newly added node or nil.
func (nl *nodeLoader) Load(lex lexeme.LexEvent) ischema.Node {
	defer lexeme.CatchLexEventError(lex)

	switch lex.Type() {
	case lexeme.NewLine:
		*(nl.nodesPerCurrentLineCount) = 0
		return nil
	case lexeme.EndTop:
		return nil
	}

	if nl.leaf == nil {
		node := ischema.NewNode(lex)
		nl.schema.SetRootNode(node)
		nl.leaf = node
		*(nl.nodesPerCurrentLineCount)++
		return node
	}

	var isNewChildNode bool
	nl.leaf, isNewChildNode = nl.leaf.Grow(lex)
	if isNewChildNode {
		*(nl.nodesPerCurrentLineCount)++
		return nl.leaf
	}

	return nil
}

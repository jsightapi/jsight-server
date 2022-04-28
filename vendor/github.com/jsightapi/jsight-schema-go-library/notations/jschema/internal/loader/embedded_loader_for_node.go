package loader

import (
	"github.com/jsightapi/jsight-schema-go-library/internal/lexeme"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema/internal/schema"
)

// This class is responsible for loading the JSON elements in the nodes of the internal representation of the SCHEMA.

type nodeLoader struct {
	schema                   *schema.Schema // Parameter from the main loader. The resulting schema
	nodesPerCurrentLineCount *uint          // Parameter from the main loader. Counts the number of nodes in a line. To check because the rule cannot be added if there is more than one nodes suitable for this in the row.
	leaf                     schema.Node    // The schema node which is processed in the moment. A schema is a tree of nodes. Starting from the root node, we fill the nodes with lexeme events. If necessary, new nodes are added. Thus, the complete scheme tree is constructed.
}

func newNodeLoader(schem *schema.Schema, nodesPerCurrentLineCount *uint) *nodeLoader {
	return &nodeLoader{
		schema:                   schem,
		nodesPerCurrentLineCount: nodesPerCurrentLineCount,
	}
}

// returns the newly added node or nil
func (nl *nodeLoader) load(lex lexeme.LexEvent) schema.Node {
	defer lexeme.CatchLexEventError(lex)

	switch lex.Type() { //nolint:exhaustive // We will handle over cases below.
	case lexeme.NewLine:
		*(nl.nodesPerCurrentLineCount) = 0
		return nil
	case lexeme.EndTop:
		return nil
	}

	if nl.leaf == nil {
		node := schema.NewNode(lex)
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

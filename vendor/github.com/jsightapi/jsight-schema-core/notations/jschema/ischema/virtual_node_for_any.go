package ischema

import (
	"sync"

	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/bytes"
	"github.com/jsightapi/jsight-schema-core/fs"
	"github.com/jsightapi/jsight-schema-core/json"
	"github.com/jsightapi/jsight-schema-core/lexeme"
	"github.com/jsightapi/jsight-schema-core/notations/jschema/ischema/constraint"
)

var once sync.Once
var virtualAnyNode Node

func VirtualNodeForAny() Node {
	once.Do(func() {
		virtualAnyNode = makeVirtualNodeForAny()
	})
	return virtualAnyNode
}

func makeVirtualNodeForAny() Node {
	f := fs.NewFile("virtual", `"" // {type: "any"}`)
	lex := lexeme.NewLexEvent(lexeme.LiteralEnd, 0, 1, f)

	node := newLiteralNode(lex)
	node.jsonType = json.TypeString
	node.AddConstraint(constraint.NewType(bytes.NewBytes(`"any"`), schema.RuleASTNodeSourceUnknown))

	return node
}

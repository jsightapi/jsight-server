package loader

import (
	"sync"

	"github.com/jsightapi/jsight-schema-go-library/internal/lexeme"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema/internal/scanner"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema/internal/schema"
)

// Contains information about the mode in which the loader is located.
// It affects how the received lexical events will be interpreted depending on whether they are in the comments or not.

type mode int

const (
	readDefault mode = iota
	readInlineComment
	readMultiLineComment
)

// Loads the schema from the scanner into the internal view.
// Does not check for the correctness of the branch because it deals with the scanner.
type loader struct {
	// The schema resulting.
	schema schema.Schema

	// rootSchema a scheme into which types can be added from the "or" rule.
	rootSchema *schema.Schema

	// scanner a tool to search for lexical events in a byte sequence containing
	// a schema.
	scanner *scanner.Scanner

	// lastAddedNode the last node added to the internal Schema.
	lastAddedNode schema.Node

	// The rule is responsible for creating constraints for SCHEMA internal representation
	// nodes from the RULES described in the SCHEMA file.
	rule *ruleLoader

	// The node class is responsible for loading the JSON elements in the nodes
	// of the internal representation of the SCHEMA.
	node *nodeLoader

	// mode used for processing inline comment, multi-line comment, or no comment
	// section.
	mode mode

	// nodesPerCurrentLineCount the number of nodes in a line. To check because
	// the rule cannot be added if there is more than one nodes suitable for this
	// in the row.
	nodesPerCurrentLineCount uint
}

func LoadSchema(scan *scanner.Scanner, rootSchema *schema.Schema, areKeysOptionalByDefault bool) *schema.Schema {
	s := LoadSchemaWithoutCompile(scan, rootSchema)
	CompileBasic(&s, areKeysOptionalByDefault)
	return &s
}

var loaderPool = sync.Pool{
	New: func() interface{} {
		return &loader{
			schema: schema.New(),
		}
	},
}

func LoadSchemaWithoutCompile(scan *scanner.Scanner, rootSchema *schema.Schema) schema.Schema {
	l := loaderPool.Get().(*loader) //nolint:errcheck // We're sure about this type.
	defer func() {
		l.reset()
		loaderPool.Put(l)
	}()

	l.scanner = scan

	l.rootSchema = rootSchema
	if rootSchema == nil {
		l.rootSchema = &l.schema
	}

	l.node = newNodeLoader(&l.schema, &l.nodesPerCurrentLineCount)
	l.doLoad()

	return l.schema
}

func (l *loader) reset() {
	l.schema = schema.New()
	l.rootSchema = nil
	l.scanner = nil
	l.lastAddedNode = nil
	l.rule = nil
	l.node = nil
	l.mode = readDefault
	l.nodesPerCurrentLineCount = 0
}

// doLoad the main function, in which there is a cycle of scanning and loading schemas.
func (l *loader) doLoad() {
	for {
		lex, ok := l.scanner.Next()
		if !ok {
			break
		}

		skip, err := l.handleLex(lex)
		if err != nil {
			panic(err)
		}

		if skip {
			continue
		}

		switch l.mode {
		case readMultiLineComment, readInlineComment:
			l.rule.load(lex)
		default:
			if node := l.node.load(lex); node != nil {
				l.lastAddedNode = node
			}
		}
	}
}

func (l *loader) handleLex(lex lexeme.LexEvent) (bool, error) {
	switch lex.Type() { //nolint:exhaustive // It's okay here.
	case lexeme.TypesShortcutBegin, lexeme.KeyShortcutBegin:
		return true, nil

	case lexeme.TypesShortcutEnd:
		l.mode = readDefault
		if err := addShortcutConstraint(l.lastAddedNode, l.rootSchema, lex); err != nil {
			return false, err
		}
		return true, nil

	case lexeme.MultiLineAnnotationBegin:
		l.mode = readMultiLineComment
		l.rule = newRuleLoader(l.lastAddedNode, l.nodesPerCurrentLineCount, l.rootSchema)
		return true, nil

	case lexeme.MultiLineAnnotationEnd:
		l.mode = readDefault
		return true, nil

	case lexeme.InlineAnnotationBegin:
		if l.mode == readDefault { // not multiLine comment
			l.mode = readInlineComment
			l.rule = newRuleLoader(l.lastAddedNode, l.nodesPerCurrentLineCount, l.rootSchema)
			return true, nil
		}

	case lexeme.InlineAnnotationEnd:
		if l.mode == readInlineComment { // not multiLine comment
			l.mode = readDefault
			return true, nil
		}
	}

	return false, nil
}

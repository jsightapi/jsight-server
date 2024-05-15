package loader

import (
	"sync"

	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/lexeme"
	"github.com/jsightapi/jsight-schema-core/notations/jschema/ischema"
	"github.com/jsightapi/jsight-schema-core/notations/jschema/scanner"
)

// mode contains information about the mode in which the loader is located.
// It affects how the received lexical events will be interpreted depending on
// whether they are in the comments or not.
type mode int

const (
	readDefault mode = iota
	readInlineComment
	readMultiLineComment
)

// loader loads the schema from the scanner into the internal view.
// Does not check for the correctness of the branch because it deals with the scanner.
type loader struct {
	// The schema resulting.
	schema ischema.ISchema

	// rootSchema a scheme into which types can be added from the "or" rule.
	rootSchema *ischema.ISchema

	// scanner a tool to search for lexical events in a byte sequence containing
	// a schema.
	scanner *scanner.Scanner

	// lastAddedNode the last node added to the internal Schema.
	lastAddedNode ischema.Node

	// rules all available rules.
	rules map[string]schema.Rule

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

func LoadSchema(scan *scanner.Scanner, rootSchema *ischema.ISchema) *ischema.ISchema {
	s := LoadSchemaWithoutCompile(scan, rootSchema, nil)
	CompileBasic(&s, false)
	return &s
}

var loaderPool = sync.Pool{
	New: func() interface{} {
		return &loader{
			schema: ischema.New(),
		}
	},
}

func LoadSchemaWithoutCompile(
	scan *scanner.Scanner,
	rootSchema *ischema.ISchema,
	rules map[string]schema.Rule,
) ischema.ISchema {
	l := loaderPool.Get().(*loader)
	defer func() {
		l.reset()
		loaderPool.Put(l)
	}()

	l.scanner = scan
	l.rules = rules

	l.rootSchema = rootSchema
	if rootSchema == nil {
		l.rootSchema = &l.schema
	}

	l.node = newNodeLoader(&l.schema, &l.nodesPerCurrentLineCount)
	l.doLoad()

	return l.schema
}

func (l *loader) reset() {
	l.schema = ischema.New()
	l.rootSchema = nil
	l.scanner = nil
	l.lastAddedNode = nil
	l.rules = nil
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
			if node := l.node.Load(lex); node != nil {
				l.lastAddedNode = node
			}
		}
	}
}

func (l *loader) handleLex(lex lexeme.LexEvent) (bool, error) {
	switch lex.Type() {
	case lexeme.TypesShortcutBegin, lexeme.KeyShortcutBegin:
		return l.mode != readMultiLineComment && l.mode != readInlineComment, nil

	case lexeme.TypesShortcutEnd:
		if l.mode == readMultiLineComment || l.mode == readInlineComment {
			return false, nil
		}

		l.mode = readDefault
		if err := addShortcutConstraint(l.lastAddedNode, l.rootSchema, lex); err != nil {
			return false, err
		}
		return true, nil

	case lexeme.MultiLineAnnotationBegin:
		l.mode = readMultiLineComment
		l.rule = newRuleLoader(l.lastAddedNode, l.nodesPerCurrentLineCount, l.rootSchema, l.rules)
		return true, nil

	case lexeme.MultiLineAnnotationEnd:
		l.mode = readDefault
		return true, nil

	case lexeme.InlineAnnotationBegin:
		if l.mode == readDefault { // not multiLine comment
			l.mode = readInlineComment
			l.rule = newRuleLoader(l.lastAddedNode, l.nodesPerCurrentLineCount, l.rootSchema, l.rules)
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

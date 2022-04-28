package loader

import (
	"strings"

	jschema "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/bytes"
	"github.com/jsightapi/jsight-schema-go-library/errors"
	"github.com/jsightapi/jsight-schema-go-library/internal/lexeme"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema/internal/schema"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema/internal/schema/constraint"
)

func addShortcutConstraint(node schema.Node, rootSchema *schema.Schema, lex lexeme.LexEvent) error {
	if lex.Type() != lexeme.TypesShortcutEnd {
		return errors.ErrLoader
	}

	// At this point lexeme value is valid, and we can safely use it.
	// Lexeme value examples:
	// - @foo
	// - @foo | @bar

	// Determines which constraint should be used.
	val := lex.Value().String()

	if strings.ContainsRune(val, '|') {
		addORShortcut(node, rootSchema, val)
	} else {
		addTypeShortcut(node, val)
	}
	return nil
}

func addORShortcut(node schema.Node, rootSchema *schema.Schema, val string) {
	// At this point lexeme value is valid, and we can safely use it.
	// Lexeme value example: "@foo | @bar"

	ss := constraint.NewTypesList(jschema.RuleASTNodeSourceGenerated)
	for _, s := range strings.Split(val, "|") {
		typ := schema.New()
		typ.SetRootNode(node)

		CompileBasic(&typ, true)

		lex := node.BasisLexEventOfSchemaForNode()
		rootSchema.AddUnnamedType(&typ, lex.File(), lex.Begin())

		s = strings.TrimSpace(s)
		ss.AddName(s, s, jschema.RuleASTNodeSourceGenerated)
	}

	node.AddConstraint(ss)
	node.AddConstraint(constraint.NewOr(jschema.RuleASTNodeSourceGenerated))
}

func addTypeShortcut(node schema.Node, val string) {
	node.AddConstraint(constraint.NewType(
		bytes.Bytes(strings.TrimSpace(val)),
		jschema.RuleASTNodeSourceGenerated,
	))
}

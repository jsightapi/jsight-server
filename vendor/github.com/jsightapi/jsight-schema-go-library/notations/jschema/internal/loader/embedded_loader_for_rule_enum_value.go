package loader

import (
	stdErrors "errors"
	"fmt"
	"strings"

	jschemaLib "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/errors"
	"github.com/jsightapi/jsight-schema-go-library/internal/lexeme"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema/internal/schema/constraint"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema/rules"
)

// Loader for "enum" rule value (array of literals). Ex: [123, 45.67, "abc", true, null]

type enumValueLoader struct {
	enumConstraint *constraint.Enum

	// stateFunc a function for running a state machine (the current state of the
	// state machine).
	stateFunc func(lexeme.LexEvent)

	// rules a set of all available rules.
	// Will be used for creating enum from one of rule.
	rules map[string]jschemaLib.Rule

	lastIdx int

	// inProgress true - if loading in progress, false - if loading finisher.
	inProgress bool
}

func newEnumValueLoader(
	enumConstraint *constraint.Enum,
	rules map[string]jschemaLib.Rule,
) embeddedLoader {
	l := &enumValueLoader{
		enumConstraint: enumConstraint,
		inProgress:     true,
		rules:          rules,
	}
	l.stateFunc = l.begin
	return l
}

// Returns false when the load is complete.
func (l *enumValueLoader) load(lex lexeme.LexEvent) bool {
	defer lexeme.CatchLexEventError(lex)
	l.stateFunc(lex)
	return l.inProgress
}

// begin of array "[", or "@"
func (l *enumValueLoader) begin(lex lexeme.LexEvent) {
	switch lex.Type() {
	case lexeme.ArrayBegin:
		l.stateFunc = l.arrayItemBeginOrArrayEnd
	case lexeme.MixedValueBegin:
		l.stateFunc = l.ruleNameBegin
	default:
		panic(errors.ErrInvalidValueInEnumRule)
	}
}

// begin of array item begin or array end
// ex: [1 <--
// ex: [" <--
// ex: ] <--
func (l *enumValueLoader) arrayItemBeginOrArrayEnd(lex lexeme.LexEvent) {
	switch lex.Type() {
	case lexeme.ArrayItemBegin:
		l.stateFunc = l.literal
	case lexeme.ArrayEnd:
		l.stateFunc = l.endOfLoading
		l.inProgress = false
	case lexeme.InlineAnnotationBegin:
		l.stateFunc = l.commentStart
	default:
		panic(errors.ErrLoader)
	}
}

func (l *enumValueLoader) commentStart(lex lexeme.LexEvent) {
	if lex.Type() != lexeme.InlineAnnotationTextBegin {
		panic(errors.ErrLoader)
	}
	l.stateFunc = l.commentEnd
}

func (l *enumValueLoader) commentEnd(lex lexeme.LexEvent) {
	if lex.Type() != lexeme.InlineAnnotationTextEnd {
		panic(errors.ErrLoader)
	}

	l.enumConstraint.SetComment(l.lastIdx, lex.Value().String())
	l.stateFunc = l.annotationEnd
}

func (l *enumValueLoader) annotationEnd(lex lexeme.LexEvent) {
	if lex.Type() != lexeme.InlineAnnotationEnd {
		panic(errors.ErrLoader)
	}
	l.stateFunc = l.arrayItemBeginOrArrayEnd
}

// array item value (literal)
func (l *enumValueLoader) literal(lex lexeme.LexEvent) {
	switch lex.Type() {
	case lexeme.LiteralBegin:
		return
	case lexeme.LiteralEnd:
		l.lastIdx = l.enumConstraint.Append(lex.Value())
		l.stateFunc = l.arrayItemEnd
	default:
		panic(errors.ErrIncorrectArrayItemTypeInEnumRule)
	}
}

// array item end
func (l *enumValueLoader) arrayItemEnd(lex lexeme.LexEvent) {
	if lex.Type() != lexeme.ArrayItemEnd {
		panic(errors.ErrLoader)
	}
	l.stateFunc = l.arrayItemBeginOrArrayEnd
}

// shortcutBeginOrArrayEnd process expected rule name.
// ex: @ <--
func (l *enumValueLoader) ruleNameBegin(lex lexeme.LexEvent) {
	if lex.Type() != lexeme.TypesShortcutBegin {
		panic(errors.ErrLoader)
	}
	l.stateFunc = l.ruleName
}

// ruleName process rule name
func (l *enumValueLoader) ruleName(lex lexeme.LexEvent) {
	if lex.Type() != lexeme.TypesShortcutEnd {
		panic(errors.ErrLoader)
	}

	v := strings.TrimSpace(string(lex.Value()))

	r, ok := l.rules[v]
	if !ok {
		panic(errors.Format(errors.ErrEnumRuleNotFound, v))
	}

	e, ok := r.(*rules.Enum)
	if !ok {
		panic(errors.Format(errors.ErrNotAnEnumRule, v))
	}

	vv, err := e.Values()
	if err != nil {
		panic(fmt.Errorf("Invalid enum %q: %s", v, getDetailsFromEnumError(err))) //nolint:stylecheck // It's expected.
	}

	l.enumConstraint.SetRuleName(v)
	for _, v := range vv {
		l.enumConstraint.Append(v)
	}
	l.stateFunc = l.endOfLoading
	l.inProgress = false
}

func getDetailsFromEnumError(err error) string {
	var de interface{ Message() string }
	if stdErrors.As(err, &de) {
		return de.Message()
	}
	return err.Error()
}

// The endOfLoading method should not be called during normal operation. Ensures
// that the loader will not continue to work after the load is complete.
func (*enumValueLoader) endOfLoading(lexeme.LexEvent) {
	panic(errors.ErrLoader)
}

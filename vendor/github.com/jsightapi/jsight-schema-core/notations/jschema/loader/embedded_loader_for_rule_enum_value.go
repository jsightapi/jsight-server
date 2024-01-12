package loader

import (
	stdErrors "errors"

	jschemaLib "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/errs"
	"github.com/jsightapi/jsight-schema-core/lexeme"
	"github.com/jsightapi/jsight-schema-core/notations/jschema/ischema/constraint"
	"github.com/jsightapi/jsight-schema-core/rules/enum"
)

// enumValueLoader loader for "enum" rule value (array of literals).
// Ex: [123, 45.67, "abc", true, null]
type enumValueLoader struct {
	enumConstraint *constraint.Enum

	// stateFunc a function for running a state machine (the current state of the
	// state machine).
	stateFunc func(lexeme.LexEvent)

	// rules a set of all available rules.
	// Will be used for creating enum from one of rule.
	rules map[string]jschemaLib.Rule

	// lastIdx index of last added enum value.
	lastIdx int

	// inProgress true - if loading in progress, false - if loading finisher.
	inProgress bool
}

var _ embeddedLoader = (*enumValueLoader)(nil)

func newEnumValueLoader(
	enumConstraint *constraint.Enum,
	rules map[string]jschemaLib.Rule,
) *enumValueLoader {
	l := &enumValueLoader{
		enumConstraint: enumConstraint,
		inProgress:     true,
		rules:          rules,
	}
	l.stateFunc = l.begin
	return l
}

func (l *enumValueLoader) Load(lex lexeme.LexEvent) bool {
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
		panic(errs.ErrInvalidValueInEnumRule.F())
	}
}

// arrayItemBeginOrArrayEnd begin of array item begin or array end
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
		panic(errs.ErrLoader.F())
	}
}

func (l *enumValueLoader) commentStart(lex lexeme.LexEvent) {
	if lex.Type() != lexeme.InlineAnnotationTextBegin {
		panic(errs.ErrLoader.F())
	}
	l.stateFunc = l.commentEnd
}

func (l *enumValueLoader) commentEnd(lex lexeme.LexEvent) {
	if lex.Type() != lexeme.InlineAnnotationTextEnd {
		panic(errs.ErrLoader.F())
	}

	l.enumConstraint.SetComment(l.lastIdx, lex.Value().String())
	l.stateFunc = l.annotationEnd
}

func (l *enumValueLoader) annotationEnd(lex lexeme.LexEvent) {
	if lex.Type() != lexeme.InlineAnnotationEnd {
		panic(errs.ErrLoader.F())
	}
	l.stateFunc = l.arrayItemBeginOrArrayEnd
}

// array item value (literal)
func (l *enumValueLoader) literal(lex lexeme.LexEvent) {
	switch lex.Type() {
	case lexeme.LiteralBegin:
	case lexeme.LiteralEnd:
		l.lastIdx = l.enumConstraint.Append(constraint.NewEnumItem(lex.Value(), ""))
		l.stateFunc = l.arrayItemEnd
	default:
		panic(errs.ErrIncorrectArrayItemTypeInEnumRule.F())
	}
}

func (l *enumValueLoader) arrayItemEnd(lex lexeme.LexEvent) {
	if lex.Type() != lexeme.ArrayItemEnd {
		panic(errs.ErrLoader.F())
	}
	l.stateFunc = l.arrayItemBeginOrArrayEnd
}

// ruleNameBegin process expected rule name.
// ex: @ <--
func (l *enumValueLoader) ruleNameBegin(lex lexeme.LexEvent) {
	if lex.Type() != lexeme.TypesShortcutBegin {
		panic(errs.ErrLoader.F())
	}
	l.stateFunc = l.ruleName
}

// ruleName process rule name
func (l *enumValueLoader) ruleName(lex lexeme.LexEvent) {
	if lex.Type() != lexeme.TypesShortcutEnd {
		panic(errs.ErrLoader.F())
	}

	v := lex.Value().TrimSpaces().String()

	r, ok := l.rules[v]
	if !ok {
		panic(errs.ErrEnumRuleNotFound.F(v))
	}

	e, ok := r.(*enum.Enum)
	if !ok {
		panic(errs.ErrNotAnEnumRule.F(v))
	}

	vv, err := e.Values()
	if err != nil {
		panic(errs.ErrInvalidEnumValues.F(v, getDetailsFromEnumError(err)))
	}

	l.enumConstraint.SetRuleName(v)
	for _, v := range vv {
		if v.Type == jschemaLib.SchemaTypeComment {
			continue
		}
		l.enumConstraint.Append(constraint.NewEnumItem(v.Value, v.Comment))
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
	panic(errs.ErrLoader.F())
}

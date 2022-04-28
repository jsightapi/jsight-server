package loader

import (
	"github.com/jsightapi/jsight-schema-go-library/errors"
	"github.com/jsightapi/jsight-schema-go-library/internal/lexeme"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema/internal/schema/constraint"
)

// Loader for "enum" rule value (array of literals). Ex: [123, 45.67, "abc", true, null]

type enumValueLoader struct {
	enumConstraint *constraint.Enum

	// stateFunc a function for running a state machine (the current state of the
	// state machine).
	stateFunc func(lexeme.LexEvent)

	lastIdx int

	// inProgress true - if loading in progress, false - if loading finisher.
	inProgress bool
}

func newEnumValueLoader(enumConstraint *constraint.Enum) embeddedLoader {
	l := new(enumValueLoader)
	l.enumConstraint = enumConstraint
	l.stateFunc = l.begin
	l.inProgress = true
	return l
}

// Returns false when the load is complete.
func (l *enumValueLoader) load(lex lexeme.LexEvent) bool {
	defer lexeme.CatchLexEventError(lex)
	l.stateFunc(lex)
	return l.inProgress
}

// begin of array "["
func (l *enumValueLoader) begin(lex lexeme.LexEvent) {
	switch lex.Type() {
	case lexeme.ArrayBegin:
		l.stateFunc = l.arrayItemBeginOrArrayEnd
	default:
		panic(errors.ErrArrayWasExpectedInEnumRule)
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
		// switch l.nodeTypesListConstraint().Len() {
		// case 0:
		// 	panic(common.ErrEmptyArrayInOrRule)
		// case 1:
		// 	panic(common.ErrOneElementInArrayInOrRule)
		// }
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
	switch lex.Type() {
	case lexeme.ArrayItemEnd:
		l.stateFunc = l.arrayItemBeginOrArrayEnd
	default:
		panic(errors.ErrLoader)
	}
}

// The method should not be called during normal operation. Ensures that the loader will not continue to work after the load is complete.
func (*enumValueLoader) endOfLoading(lexeme.LexEvent) {
	panic(errors.ErrLoader)
}

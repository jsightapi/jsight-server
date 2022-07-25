package loader

import (
	"github.com/jsightapi/jsight-schema-go-library/errors"
	"github.com/jsightapi/jsight-schema-go-library/internal/lexeme"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema/internal/schema/constraint"
)

// Loader for "allOf" rule value (string or array)
// example: "@name"
// example: ["@name1", "@name2"]

type allOfValueLoader struct {
	allOfConstraint *constraint.AllOf

	// stateFunc a function for running a state machine (the current state of the
	// state machine).
	stateFunc func(lexeme.LexEvent)

	// inProgress indicates loading finished.
	inProgress bool
}

func newAllOfValueLoader(allOfConstraint *constraint.AllOf) embeddedLoader {
	l := new(allOfValueLoader)
	l.allOfConstraint = allOfConstraint
	l.stateFunc = l.begin
	l.inProgress = true
	return l
}

// Returns false when the load is complete.
func (l *allOfValueLoader) load(lex lexeme.LexEvent) bool {
	defer lexeme.CatchLexEventError(lex)
	l.stateFunc(lex)
	return l.inProgress
}

// begin of array "[" or scalar '"'
func (l *allOfValueLoader) begin(lex lexeme.LexEvent) {
	switch lex.Type() {
	case lexeme.ArrayBegin:
		l.stateFunc = l.arrayItemBeginOrArrayEnd
	case lexeme.LiteralBegin:
		l.stateFunc = l.scalarValue
	default:
		panic(errors.ErrUnacceptableValueInAllOfRule)
	}
}

// begin of array item or array end
func (l *allOfValueLoader) arrayItemBeginOrArrayEnd(lex lexeme.LexEvent) {
	switch lex.Type() {
	case lexeme.ArrayItemBegin:
		l.stateFunc = l.arrayItemValue
	case lexeme.ArrayEnd:
		l.stateFunc = l.endOfLoading
		l.inProgress = false
	default:
		panic(errors.ErrLoader)
	}
}

// array item value
func (l *allOfValueLoader) arrayItemValue(lex lexeme.LexEvent) {
	switch lex.Type() {
	case lexeme.LiteralBegin:
		return
	case lexeme.LiteralEnd:
		l.allOfConstraint.Append(lex.Value())
		l.stateFunc = l.arrayItemEnd
	default:
		panic(errors.ErrUnacceptableValueInAllOfRule)
	}
}

// array item end
func (l *allOfValueLoader) arrayItemEnd(lex lexeme.LexEvent) {
	switch lex.Type() {
	case lexeme.ArrayItemEnd:
		l.stateFunc = l.arrayItemBeginOrArrayEnd
	default:
		panic(errors.ErrLoader)
	}
}

// scalar value
func (l *allOfValueLoader) scalarValue(lex lexeme.LexEvent) {
	switch lex.Type() {
	case lexeme.LiteralEnd:
		l.allOfConstraint.Append(lex.Value())
		l.stateFunc = l.endOfLoading
		l.inProgress = false
	default:
		panic(errors.ErrUnacceptableValueInAllOfRule)
	}
}

// The method should not be called during normal operation. Ensures that the loader will not continue to work after the load is complete.
func (*allOfValueLoader) endOfLoading(lexeme.LexEvent) {
	panic(errors.ErrLoader)
}

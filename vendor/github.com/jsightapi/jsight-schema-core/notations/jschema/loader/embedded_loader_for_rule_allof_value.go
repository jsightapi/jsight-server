package loader

import (
	"github.com/jsightapi/jsight-schema-core/errs"
	"github.com/jsightapi/jsight-schema-core/lexeme"
	"github.com/jsightapi/jsight-schema-core/notations/jschema/ischema/constraint"
)

// allOfValueLoader loader for "allOf" rule value (string or array).
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

var _ embeddedLoader = (*allOfValueLoader)(nil)

func newAllOfValueLoader(allOfConstraint *constraint.AllOf) *allOfValueLoader {
	l := &allOfValueLoader{
		allOfConstraint: allOfConstraint,
		inProgress:      true,
	}
	l.stateFunc = l.begin
	return l
}

func (l *allOfValueLoader) Load(lex lexeme.LexEvent) bool {
	defer lexeme.CatchLexEventError(lex)
	l.stateFunc(lex)
	return l.inProgress
}

// begin of array "[" or scalar '"'.
func (l *allOfValueLoader) begin(lex lexeme.LexEvent) {
	switch lex.Type() {
	case lexeme.ArrayBegin:
		l.stateFunc = l.arrayItemBeginOrArrayEnd
	case lexeme.LiteralBegin:
		l.stateFunc = l.scalarValue
	default:
		panic(errs.ErrUnacceptableValueInAllOfRule.F())
	}
}

// arrayItemBeginOrArrayEnd begin of array item or array end.
func (l *allOfValueLoader) arrayItemBeginOrArrayEnd(lex lexeme.LexEvent) {
	switch lex.Type() {
	case lexeme.ArrayItemBegin:
		l.stateFunc = l.arrayItemValue
	case lexeme.ArrayEnd:
		l.stateFunc = l.endOfLoading
		l.inProgress = false
	default:
		panic(errs.ErrLoader.F())
	}
}

func (l *allOfValueLoader) arrayItemValue(lex lexeme.LexEvent) {
	switch lex.Type() {
	case lexeme.LiteralBegin:
		return
	case lexeme.LiteralEnd:
		l.allOfConstraint.Append(lex.Value())
		l.stateFunc = l.arrayItemEnd
	default:
		panic(errs.ErrUnacceptableValueInAllOfRule.F())
	}
}

func (l *allOfValueLoader) arrayItemEnd(lex lexeme.LexEvent) {
	if lex.Type() != lexeme.ArrayItemEnd {
		panic(errs.ErrLoader.F())
	}
	l.stateFunc = l.arrayItemBeginOrArrayEnd
}

func (l *allOfValueLoader) scalarValue(lex lexeme.LexEvent) {
	if lex.Type() != lexeme.LiteralEnd {
		panic(errs.ErrUnacceptableValueInAllOfRule.F())
	}
	l.allOfConstraint.Append(lex.Value())
	l.stateFunc = l.endOfLoading
	l.inProgress = false
}

// endOfLoading the method should not be called during normal operation. Ensures
// that the loader will not continue to work after the load is complete.
func (*allOfValueLoader) endOfLoading(lexeme.LexEvent) {
	panic(errs.ErrLoader.F())
}

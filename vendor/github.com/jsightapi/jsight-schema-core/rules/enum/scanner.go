package enum

import (
	stdErrors "errors"

	"github.com/jsightapi/jsight-schema-core/bytes"
	"github.com/jsightapi/jsight-schema-core/errs"
	"github.com/jsightapi/jsight-schema-core/fs"
	"github.com/jsightapi/jsight-schema-core/internal/ds"
	"github.com/jsightapi/jsight-schema-core/kit"
	"github.com/jsightapi/jsight-schema-core/lexeme"
)

type stepFunc func(byte) (state, error)

// state values are returned by the state transition functions assigned to
// scanner.state and the method scanner.eof.
// They give details about the current state of the scan that callers might be
// interested to know about.
// It is okay to ignore the return value of any particular call to scanner.state.
type state uint8

const (
	// scanSkip indicates an uninteresting byte, so we can keep scanning forward.
	scanSkip state = iota

	// scanBeginLiteral indicates beginning of any value outside an array or object.
	scanBeginLiteral
)

// scanner represents a scanner is a JSchema scanning state machine.
// Callers call scan.reset() and then pass bytes in one at a time
// by calling scan.step(&scan, c) for each byte.
// The return value, referred to as an opcode, tells the
// caller about significant parsing events like beginning
// and ending literals, objects, and arrays, so that the
// caller can follow along if it wishes.
// The return value scanEnd indicates that a single top-level
// JSON value has been completed, *before* the byte that
// just got passed in.  (The indication must be delayed in order
// to recognize the end of numbers: is 123 a whole value or
// the beginning of 12345e+6?).
type scanner struct {
	// step is a func to be called to execute the next transition.
	// Also tried using an integer constant and a single func
	// with a switch, but using the func directly was 10% faster
	// on a 64-bit Mac Mini, and it's nicer to read.
	step stepFunc

	// returnToStep a stack of step functions, to preserve the sequence of steps
	// (and return to them) in some cases.
	returnToStep *ds.Stack[stepFunc]

	// stack a stack of found lexical event. The stack is needed for the scanner
	// to take into account the nesting of SCHEME elements.
	stack *ds.Stack[lexeme.LexEvent]

	// uniqueValues represent a map of found values.
	// Useful for duplication tracking.
	uniqueValues map[enumItemValue]struct{}

	// file a structure containing jSchema data.
	file *fs.File

	// data jSchema content.
	data bytes.Bytes

	// finds a list of found types of lexical event for the current step. Several
	// lexical events can be found in one step (example: ArrayItemBegin and LiteralBegin).
	finds []lexeme.LexEventType

	// index scanned byte index.
	index bytes.Index

	// dataSize a size of schema data in bytes. Count once for optimization.
	dataSize bytes.Index

	// annotation one of the possible States of annotation processing (annotationNone,
	// annotationInline).
	annotation bool

	// unfinishedLiteral a sign that a literal has been started but not completed.
	unfinishedLiteral bool

	// lengthComputing used when a file contains data after the schema (for example,
	// in jApi).
	lengthComputing bool

	hasTrailingCharacters bool
}

func newScanner(file *fs.File, oo ...scannerOption) *scanner {
	content := file.Content()

	s := &scanner{
		file:         file,
		data:         content,
		dataSize:     content.LenIndex(),
		returnToStep: &ds.Stack[stepFunc]{},
		stack:        &ds.Stack[lexeme.LexEvent]{},
		uniqueValues: map[enumItemValue]struct{}{},
		finds:        make([]lexeme.LexEventType, 0, 3),
	}

	s.step = s.stateBegin

	for _, o := range oo {
		o(s)
	}

	return s
}

type scannerOption func(*scanner)

// scannerComputeLength switch scanner in length computing mode.
// scanner in this mode shouldn't be used for parsing.
func scannerComputeLength(s *scanner) {
	s.lengthComputing = true
}

func (s *scanner) Length() (uint, error) {
	if !s.lengthComputing {
		return 0, errs.ErrRuntimeFailure.F()
	}
	var length uint
	for {
		lex, err := s.Next()
		if stdErrors.Is(err, errEOS) {
			break
		}
		if err != nil {
			return 0, err
		}

		if lex.Type() == lexeme.EndTop {
			// Found character after the end of the schema and spaces.
			// Example: char "s" in "{} some text"
			length = uint(lex.End()) - 1
			break
		}

		length = uint(lex.End()) + 1
		if lex.End() == s.dataSize {
			length--
		}
	}
	for ; length > 0; length-- {
		c := s.data.Byte(length - 1)
		if !bytes.IsBlank(c) {
			break
		}
	}
	return length, nil
}

var errEOS = stdErrors.New("end of stream")

// Next reads schema byte by byte.
// Stops if it detects lexical events.
// Returns pointer to found lexeme event, or nil if you have complete reading.
func (s *scanner) Next() (lexeme.LexEvent, error) {
	if len(s.finds) != 0 {
		lex, err := s.shiftFound()
		if err != nil {
			return lexeme.LexEvent{}, err
		}
		return s.processingFoundLexeme(lex)
	}

	for s.index < s.dataSize {
		c := s.data.Byte(s.index)
		s.index++

		_, err := s.step(c)
		if err != nil {
			return lexeme.LexEvent{}, err
		}

		if len(s.finds) != 0 {
			lex, err := s.shiftFound()
			if err != nil {
				return lexeme.LexEvent{}, err
			}
			return s.processingFoundLexeme(lex)
		}
	}

	return s.processTail()
}

func (s *scanner) processTail() (lexeme.LexEvent, error) {
	if s.stack.Len() == 0 {
		return lexeme.LexEvent{}, errEOS
	}

	s.index++
	switch s.stack.Peek().Type() {
	case lexeme.LiteralBegin:
		if s.unfinishedLiteral {
			break
		}
		return s.processingFoundLexeme(lexeme.LiteralEnd)

	case lexeme.InlineAnnotationBegin:
		return s.processingFoundLexeme(lexeme.InlineAnnotationEnd)

	case lexeme.InlineAnnotationTextBegin:
		return s.processingFoundLexeme(lexeme.InlineAnnotationTextEnd)

	case lexeme.MultiLineAnnotationBegin:
		return s.processingFoundLexeme(lexeme.MultiLineAnnotationEnd)

	case lexeme.MultiLineAnnotationTextBegin:
		return s.processingFoundLexeme(lexeme.MultiLineAnnotationTextEnd)
	}

	err := kit.NewJSchemaError(s.file, errs.ErrUnexpectedEOF.F())
	err.SetIndex(s.dataSize - 1)
	return lexeme.LexEvent{}, err
}

// stateBegin first state of the scanner.
// Expects open square brace as the start of the enum values.
func (s *scanner) stateBegin(c byte) (state, error) {
	if bytes.IsBlank(c) {
		return scanSkip, nil
	}

	if c != '[' {
		err := kit.NewJSchemaError(s.file, errs.ErrEnumArrayExpected.F())
		err.SetIndex(s.index - 1)
		return scanSkip, err
	}

	s.found(lexeme.ArrayBegin)
	s.step = s.stateFoundArrayItemBeginOrEmpty
	return scanSkip, nil
}

func (s *scanner) stateFoundArrayItemBeginOrEmpty(c byte) (state, error) {
	if bytes.IsNewLine(c) {
		if s.annotation {
			return scanSkip, s.newJSchemaErrorAtCharacter("inside inline annotation")
		}
		s.found(lexeme.NewLine)
		return scanSkip, nil
	}

	if c == ']' {
		return s.stateFoundArrayEnd()
	}

	r, err := s.stateBeginArrayItemOrEmpty(c)
	if err != nil {
		return scanSkip, err
	}
	if r == scanBeginLiteral {
		s.found(lexeme.ArrayItemBegin)
		s.found(lexeme.LiteralBegin)
	}
	return r, nil
}

func (s *scanner) stateFoundArrayItemBegin(c byte) (state, error) {
	r, err := s.stateBeginValue(c)
	if err != nil {
		return scanSkip, err
	}

	if r == scanBeginLiteral {
		s.found(lexeme.ArrayItemBegin)
		s.found(lexeme.LiteralBegin)
	}
	return r, nil
}

func (s *scanner) stateBeginValue(c byte) (state, error) {
	if bytes.IsNewLine(c) {
		if s.annotation {
			return scanSkip, s.newJSchemaErrorAtCharacter("inside inline annotation")
		}
		s.found(lexeme.NewLine)
		return scanSkip, nil
	}
	if bytes.IsBlank(c) {
		return scanSkip, nil
	}
	if s.isAnnotationStart(c) {
		return scanSkip, s.switchToAnnotation()
	}
	switch c {
	case '"':
		s.step = s.stateInString
		s.unfinishedLiteral = true
		return scanBeginLiteral, nil
	case '-':
		s.step = s.stateNeg
		s.unfinishedLiteral = true
		return scanBeginLiteral, nil
	case '0': // beginning of 0.123
		s.step = s.state0
		return scanBeginLiteral, nil
	case 't': // beginning of true
		s.step = s.stateT
		s.unfinishedLiteral = true
		return scanBeginLiteral, nil
	case 'f': // beginning of false
		s.step = s.stateF
		s.unfinishedLiteral = true
		return scanBeginLiteral, nil
	case 'n': // beginning of null
		s.step = s.stateN
		s.unfinishedLiteral = true
		return scanBeginLiteral, nil
	}
	if '1' <= c && c <= '9' { // beginning of 1234.5
		s.step = s.state1
		return scanBeginLiteral, nil
	}
	return scanSkip, s.newJSchemaErrorAtCharacter("looking for beginning of value")
}

// After reading `[`.
func (s *scanner) stateBeginArrayItemOrEmpty(c byte) (state, error) {
	if c == ']' {
		return s.stateFoundArrayEnd()
	}
	return s.stateBeginValue(c)
}

func (s *scanner) stateEndValue(c byte) (state, error) {
	length := s.stack.Len()

	if length == 0 { // json ex `{} `
		s.step = s.stateEndTop
		return s.step(c)
	}

	t := s.stack.Peek().Type()

	if t == lexeme.LiteralBegin {
		s.found(lexeme.LiteralEnd)

		if err := s.validateValue(); err != nil {
			return scanSkip, err
		}

		if length == 1 { // json ex `123 `
			s.step = s.stateEndTop
			return s.step(c)
		}

		t = s.stack.Get(length - 2).Type()
	}

	if t == lexeme.ArrayItemBegin {
		s.found(lexeme.ArrayItemEnd)
		s.step = s.stateAfterArrayItem
		return s.step(c)
	}
	if s.lengthComputing && t == lexeme.InlineAnnotationBegin {
		s.annotation = false
		_ = s.stack.Pop()
		s.step = s.returnToStep.Pop()
		return s.step(c)
	}

	return scanSkip, s.newJSchemaErrorAtCharacter("at the end of value")
}

func (s *scanner) validateValue() error {
	begin := s.stack.Peek().Begin()

	v := s.file.Content().Sub(begin, s.index-1)
	key := newEnumItem(v)
	if _, ok := s.uniqueValues[key]; ok {
		e := errs.ErrDuplicationInEnumRule.F(v.String())
		err := kit.NewJSchemaError(s.file, e)
		err.SetIndex(begin)
		return err
	}
	s.uniqueValues[key] = struct{}{}
	return nil
}

func (s *scanner) stateAfterArrayItem(c byte) (state, error) {
	if bytes.IsNewLine(c) {
		if s.annotation {
			return scanSkip, s.newJSchemaErrorAtCharacter("inside inline annotation")
		}
		s.found(lexeme.NewLine)
		return scanSkip, nil
	}
	if bytes.IsBlank(c) {
		return scanSkip, nil
	}
	if s.isAnnotationStart(c) {
		return scanSkip, s.switchToAnnotation()
	}
	if c == ',' {
		s.step = s.stateFoundArrayItemBegin
		return scanSkip, nil
	}
	if c == ']' {
		return s.stateFoundArrayEnd()
	}
	return scanSkip, s.newJSchemaErrorAtCharacter("after array item")
}

func (s *scanner) stateFoundArrayEnd() (state, error) {
	s.found(lexeme.ArrayEnd)
	if s.stack.Len() == 0 {
		s.step = s.stateEndTop
	} else {
		s.step = s.stateEndValue
	}
	return scanSkip, nil
}

// stateEndTop is the state after finishing the top-level value,
// such as after reading `{}` or `[1,2,3]`.
// Only space characters should be seen now.
func (s *scanner) stateEndTop(c byte) (state, error) {
	switch {
	case bytes.IsNewLine(c):
		if s.annotation {
			return scanSkip, s.newJSchemaErrorAtCharacter("inside inline annotation")
		}
		s.found(lexeme.NewLine)
		return scanSkip, nil

	case s.isAnnotationStart(c):
		return scanSkip, s.switchToAnnotation()

	case !bytes.IsBlank(c):
		if s.lengthComputing {
			if s.stack.Len() > 0 {
				// Looks like we have invalid schema, and we should keep scanning.
				s.hasTrailingCharacters = true
				return scanSkip, nil
			}
			s.found(lexeme.EndTop)
			return scanSkip, errEOS
		} else if !s.annotation {
			return scanSkip, s.newJSchemaErrorAtCharacter("non-space byte after top-level value")
		}
	}

	if s.hasTrailingCharacters {
		s.found(lexeme.EndTop)
		return scanSkip, errEOS
	}
	return scanSkip, nil
}

// After reading `"`.
func (s *scanner) stateInString(c byte) (state, error) {
	switch c {
	case '"':
		s.step = s.stateEndValue
		s.unfinishedLiteral = false
		return scanSkip, nil
	case '\\':
		s.step = s.stateInStringEsc
		return scanSkip, nil
	}
	if c < 0x20 {
		return scanSkip, s.newJSchemaErrorAtCharacter("in string literal")
	}
	return scanSkip, nil
}

// After reading `"\` during a quoted string.
func (s *scanner) stateInStringEsc(c byte) (state, error) {
	switch c {
	case 'b', 'f', 'n', 'r', 't', '\\', '/', '"':
		s.step = s.stateInString
		return scanSkip, nil
	case 'u':
		s.returnToStep.Push(s.stateInString)
		s.step = s.stateInStringEscU
		return scanSkip, nil
	}
	return scanSkip, s.newJSchemaErrorAtCharacter("in string escape code")
}

// After reading `"\u` during a quoted string.
func (s *scanner) stateInStringEscU(c byte) (state, error) {
	if bytes.IsHexDigit(c) {
		s.step = s.stateInStringEscU1
		return scanSkip, nil
	}
	return scanSkip, s.newJSchemaErrorAtCharacter("in \\u hexadecimal character escape")
}

// After reading `"\u1` during a quoted string.
func (s *scanner) stateInStringEscU1(c byte) (state, error) {
	if bytes.IsHexDigit(c) {
		s.step = s.stateInStringEscU12
		return scanSkip, nil
	}
	return scanSkip, s.newJSchemaErrorAtCharacter("in \\u hexadecimal character escape")
}

// After reading `"\u12` during a quoted string.
func (s *scanner) stateInStringEscU12(c byte) (state, error) {
	if bytes.IsHexDigit(c) {
		s.step = s.stateInStringEscU123
		return scanSkip, nil
	}
	return scanSkip, s.newJSchemaErrorAtCharacter("in \\u hexadecimal character escape")
}

// After reading `"\u123` during a quoted string.
func (s *scanner) stateInStringEscU123(c byte) (state, error) {
	if bytes.IsHexDigit(c) {
		s.step = s.returnToStep.Pop()
		return scanSkip, nil
	}
	return scanSkip, s.newJSchemaErrorAtCharacter("in \\u hexadecimal character escape")
}

// After reading `-` during a number.
func (s *scanner) stateNeg(c byte) (state, error) {
	if c == '0' {
		s.step = s.state0
		s.unfinishedLiteral = false
		return scanSkip, nil
	}
	if '1' <= c && c <= '9' {
		s.step = s.state1
		s.unfinishedLiteral = false
		return scanSkip, nil
	}
	return scanSkip, s.newJSchemaErrorAtCharacter("in numeric literal")
}

// After reading a non-zero integer during a number, such as after reading `1` or
// `100` but not `0`.
func (s *scanner) state1(c byte) (state, error) {
	if bytes.IsDigit(c) {
		s.step = s.state1
		return scanSkip, nil
	}
	return s.state0(c)
}

// After reading `0` during a number.
func (s *scanner) state0(c byte) (state, error) {
	if c == '.' {
		s.unfinishedLiteral = true
		s.step = s.stateDot
		return scanSkip, nil
	}
	if c == 'e' || c == 'E' {
		return scanSkip, s.newJSchemaErrorAtCharacter(messageEIsNotAllowed)
	}
	return s.stateEndValue(c)
}

// After reading the integer and decimal point in a number, such as after reading `1.`.
func (s *scanner) stateDot(c byte) (state, error) {
	if bytes.IsDigit(c) {
		s.unfinishedLiteral = false
		s.step = s.stateDot0
		return scanSkip, nil
	}
	return scanSkip, s.newJSchemaErrorAtCharacter("after decimal point in numeric literal")
}

// After reading the integer, decimal point, and subsequent digits of a number,
// such as after reading `3.14`.
func (s *scanner) stateDot0(c byte) (state, error) {
	if bytes.IsDigit(c) {
		return scanSkip, nil
	}
	if c == 'e' || c == 'E' {
		return scanSkip, s.newJSchemaErrorAtCharacter(messageEIsNotAllowed)
	}
	return s.stateEndValue(c)
}

// After reading `t`.
func (s *scanner) stateT(c byte) (state, error) {
	if c == 'r' {
		s.step = s.stateTr
		return scanSkip, nil
	}
	return scanSkip, s.newJSchemaErrorAtCharacter("in literal true (expecting 'r')")
}

// After reading `tr`.
func (s *scanner) stateTr(c byte) (state, error) {
	if c == 'u' {
		s.step = s.stateTru
		return scanSkip, nil
	}
	return scanSkip, s.newJSchemaErrorAtCharacter("in literal true (expecting 'u')")
}

// After reading `tru`.
func (s *scanner) stateTru(c byte) (state, error) {
	if c == 'e' {
		s.step = s.stateEndValue
		s.unfinishedLiteral = false
		return scanSkip, nil
	}
	return scanSkip, s.newJSchemaErrorAtCharacter("in literal true (expecting 'e')")
}

// After reading `f`.
func (s *scanner) stateF(c byte) (state, error) {
	if c == 'a' {
		s.step = s.stateFa
		return scanSkip, nil
	}
	return scanSkip, s.newJSchemaErrorAtCharacter("in literal false (expecting 'a')")
}

// After reading `fa`.
func (s *scanner) stateFa(c byte) (state, error) {
	if c == 'l' {
		s.step = s.stateFal
		return scanSkip, nil
	}
	return scanSkip, s.newJSchemaErrorAtCharacter("in literal false (expecting 'l')")
}

// After reading `fal`.
func (s *scanner) stateFal(c byte) (state, error) {
	if c == 's' {
		s.step = s.stateFals
		return scanSkip, nil
	}
	return scanSkip, s.newJSchemaErrorAtCharacter("in literal false (expecting 's')")
}

// After reading `fals`.
func (s *scanner) stateFals(c byte) (state, error) {
	if c == 'e' {
		s.step = s.stateEndValue
		s.unfinishedLiteral = false
		return scanSkip, nil
	}
	return scanSkip, s.newJSchemaErrorAtCharacter("in literal false (expecting 'e')")
}

// After reading `n`.
func (s *scanner) stateN(c byte) (state, error) {
	if c == 'u' {
		s.step = s.stateNu
		return scanSkip, nil
	}
	return scanSkip, s.newJSchemaErrorAtCharacter("in literal null (expecting 'u')")
}

// After reading `nu`.
func (s *scanner) stateNu(c byte) (state, error) {
	if c == 'l' {
		s.step = s.stateNul
		return scanSkip, nil
	}
	return scanSkip, s.newJSchemaErrorAtCharacter("in literal null (expecting 'l')")
}

// After reading `nul`.
func (s *scanner) stateNul(c byte) (state, error) {
	if c == 'l' {
		s.step = s.stateEndValue
		s.unfinishedLiteral = false
		return scanSkip, nil
	}
	return scanSkip, s.newJSchemaErrorAtCharacter("in literal null (expecting 'l')")
}

func (s *scanner) stateAnyAnnotationStart(c byte) (st state, err error) {
	switch c {
	case '/':
		s.annotation = true
		s.found(lexeme.InlineAnnotationBegin)
		s.step = s.stateInlineAnnotation
	case '*':
		s.annotation = true
		s.found(lexeme.MultiLineAnnotationBegin)
		s.step = s.stateMultiLineAnnotation
	default:
		err = s.newJSchemaErrorAtCharacter("after first slash")
	}
	return scanSkip, err
}

func (s *scanner) stateInlineAnnotation(c byte) (state, error) {
	if bytes.IsBlank(c) {
		return scanSkip, nil
	}

	s.found(lexeme.InlineAnnotationTextBegin)
	s.step = s.stateInlineAnnotationText
	return s.step(c)
}

func (s *scanner) stateMultiLineAnnotation(c byte) (state, error) {
	if bytes.IsNewLine(c) {
		s.found(lexeme.NewLine)
		return scanSkip, nil
	}
	if bytes.IsBlank(c) {
		return scanSkip, nil
	}
	s.found(lexeme.MultiLineAnnotationTextBegin)
	s.step = s.stateMultiLineAnnotationText
	return s.step(c)
}

func (s *scanner) stateMultiLineAnnotationText(c byte) (state, error) {
	if c == '*' && s.data.Byte(s.index) == '/' {
		s.found(lexeme.MultiLineAnnotationTextEnd)
		s.step = s.stateMultiLineAnnotationEnd
	}
	return scanSkip, nil
}

func (s *scanner) stateMultiLineAnnotationEnd(c byte) (state, error) {
	if c != '/' {
		return scanSkip, s.newJSchemaErrorAtCharacter("in multi-line annotation after \"*\" character")
	}
	// after *
	s.found(lexeme.MultiLineAnnotationEnd)
	s.step = s.returnToStep.Pop()
	s.annotation = false
	return scanSkip, nil
}

func (s *scanner) stateInlineAnnotationText(c byte) (state, error) {
	if bytes.IsNewLine(c) {
		s.found(lexeme.InlineAnnotationTextEnd)
		s.found(lexeme.InlineAnnotationEnd)
		s.found(lexeme.NewLine)
		s.step = s.returnToStep.Pop()
		s.annotation = false
	}
	return scanSkip, nil
}

const messageEIsNotAllowed = "isn't allowed 'cause not obvious it's a float or an integer"

func (s *scanner) found(lexType lexeme.LexEventType) {
	s.finds = append(s.finds, lexType)
}

func (s *scanner) shiftFound() (lexeme.LexEventType, error) {
	length := len(s.finds)
	if length == 0 {
		return 0, stdErrors.New("empty set of found lexical event")
	}
	lexType := s.finds[0]
	copy(s.finds[0:], s.finds[1:])
	s.finds = s.finds[:length-1]
	return lexType, nil
}

func (s *scanner) newJSchemaErrorAtCharacter(context string) kit.JSchemaError {
	// Make runes (utf8 symbols) from current index to last of slice s.data.
	// Get first rune. Then make string with format ' symbol '
	r := s.data.SubLow(s.index - 1).DecodeRune()
	e := errs.ErrInvalidCharacter.F(string(r), context)
	err := kit.NewJSchemaError(s.file, e)
	err.SetIndex(s.index - 1)
	return err
}

func (s *scanner) processingFoundLexeme(lexType lexeme.LexEventType) (lexeme.LexEvent, error) {
	i := s.index - 1
	switch {
	case lexType == lexeme.NewLine || lexType == lexeme.EndTop:
		return lexeme.NewLexEvent(lexType, i, i, s.file), nil

	case lexType.IsOpening():
		// `[`, `"` or literal first character (ex: `1` in `123`).
		lex := lexeme.NewLexEvent(lexType, i, i, s.file)
		s.stack.Push(lex)
		return lex, nil
	}

	return s.processingFoundLexemeClosingTag(lexType, i)
}

func (s *scanner) processingFoundLexemeClosingTag(lexType lexeme.LexEventType, i bytes.Index) (lexeme.LexEvent, error) {
	pair := s.stack.Pop()
	pairType := pair.Type()

	switch {
	case isNonScalarPair(pairType, lexType):
		return lexeme.NewLexEvent(lexType, pair.Begin(), i, s.file), nil

	case isScalarPair(pairType, lexType):
		if lexType == lexeme.MixedValueEnd && s.data.Byte(i-1) == ' ' {
			i--
		}
		return lexeme.NewLexEvent(lexType, pair.Begin(), i-1, s.file), nil
	}
	return lexeme.LexEvent{}, stdErrors.New("incorrect ending of the lexical event")
}

func isNonScalarPair(pairType, lexType lexeme.LexEventType) bool {
	return (pairType == lexeme.ArrayBegin && lexType == lexeme.ArrayEnd) ||
		(pairType == lexeme.MultiLineAnnotationBegin && lexType == lexeme.MultiLineAnnotationEnd)
}

func isScalarPair(pairType, lexType lexeme.LexEventType) bool {
	return (pairType == lexeme.LiteralBegin && lexType == lexeme.LiteralEnd) ||
		(pairType == lexeme.ArrayItemBegin && lexType == lexeme.ArrayItemEnd) ||
		(pairType == lexeme.MultiLineAnnotationTextBegin && lexType == lexeme.MultiLineAnnotationTextEnd) ||
		(pairType == lexeme.InlineAnnotationTextBegin && lexType == lexeme.InlineAnnotationTextEnd) ||
		(pairType == lexeme.InlineAnnotationBegin && lexType == lexeme.InlineAnnotationEnd) ||
		(pairType == lexeme.MixedValueBegin && lexType == lexeme.MixedValueEnd)
}

func (*scanner) isAnnotationStart(c byte) bool {
	return c == '/'
}

func (s *scanner) switchToAnnotation() error {
	if s.annotation {
		return s.newJSchemaErrorAtCharacter("inside inline annotation")
	}
	s.returnToStep.Push(s.step)
	s.step = s.stateAnyAnnotationStart
	return nil
}

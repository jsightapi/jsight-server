package json

import (
	"github.com/jsightapi/jsight-schema-go-library/bytes"
	"github.com/jsightapi/jsight-schema-go-library/errors"
	"github.com/jsightapi/jsight-schema-go-library/fs"
	"github.com/jsightapi/jsight-schema-go-library/internal/ds"
	"github.com/jsightapi/jsight-schema-go-library/internal/lexeme"
)

type (
	state    uint8
	stepFunc func(*scanner, byte) state
)

// These values are returned by the state transition functions
// assigned to scanner.state and the method scanner.eof.
// They give details about the current state of the scan that
// callers might be interested to know about.
// It is okay to ignore the return value of any particular
// call to scanner.state.
const (
	// scanContinue indicates an uninteresting byte, so we can keep scanning forward.
	scanContinue state = iota // uninteresting byte

	// scanSkipSpace indicates a space byte, can be skipped.
	scanSkipSpace

	// scanBeginObject indicates beginning of an object.
	scanBeginObject

	// scanObjectKey indicates finished object key (string)
	scanObjectKey

	// scanObjectValue indicates finished non-last value in an object.
	scanObjectValue

	// scanEndObject indicates the end of object (implies scanObjectValue if
	// possible).
	scanEndObject

	// scanBeginArray indicates beginning of an array.
	scanBeginArray

	// scanArrayValue indicates finished array value.
	scanArrayValue

	// scanEndArray indicates the end of array (implies scanArrayValue if possible).
	scanEndArray

	// scanBeginLiteral indicates beginning of any value outside an array or object.
	scanBeginLiteral

	// scanEnd indicates the end of the scanning. Top-level value ended *before*
	// this byte.
	scanEnd
)

// scanner represents a scanner is a JSON scanning state machine.
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

	// file a structure containing JSON data.
	file *fs.File

	// data a JSON content.
	data bytes.Bytes

	// stack a stack of found lexical event. The stack is needed for the scanner
	// to take into account the nesting of JSON or SCHEME elements.
	stack *ds.Stack[lexeme.LexEvent]

	// finds a list of found types of lexical event for the current step. Several
	// lexical events can be found in one step (example: ArrayItemBegin and LiteralBegin).
	finds []lexeme.LexEventType

	// index scanned byte index.
	index bytes.Index

	// dataSize a size of JSON data in bytes. Count once for optimization.
	dataSize bytes.Index

	// unfinishedLiteral a sign that a literal has been started but not completed.
	unfinishedLiteral bool

	// allowTrailingNonSpaceCharacters allows to have non-empty characters at the
	// end of the JSON.
	allowTrailingNonSpaceCharacters bool
}

func newScanner(file *fs.File) *scanner {
	return &scanner{
		step:         stateFoundRootValue,
		file:         file,
		data:         file.Content(),
		dataSize:     bytes.Index(len(file.Content())),
		returnToStep: &ds.Stack[stepFunc]{},
		stack:        &ds.Stack[lexeme.LexEvent]{},
		finds:        make([]lexeme.LexEventType, 0, 3),
	}
}

func (s *scanner) Length() uint {
	var length uint
	for {
		lex, ok := s.Next()
		if !ok {
			break
		}

		if lex.Type() == lexeme.EndTop {
			// Found character after the end of the schema and spaces. Ex: char "s" in "{} some text"
			length = uint(lex.End()) - 1
			break
		}
		length = uint(lex.End()) + 1
	}
	for {
		if length == 0 {
			break
		}
		c := s.data[length-1]
		if bytes.IsBlank(c) {
			length--
		} else {
			break
		}
	}
	return length
}

// Next reads JSON byte by byte.
// Panic if an invalid JSON structure is found.
// Stops if it detects lexical events.
// Returns pointer to found lexeme event, or nil if you have complete JSON reading.
func (s *scanner) Next() (lexeme.LexEvent, bool) {
	if len(s.finds) != 0 {
		return s.processingFoundLexeme(s.shiftFound()), true
	}

	for s.index < s.dataSize {
		c := s.data[s.index]
		s.index++

		// useful for debugging comment below 1 line for release
		// fmt.Println("Next->step", runtime.FuncForPC(reflect.ValueOf(s.step).Pointer()).Name())

		s.step(s, c)

		if len(s.finds) != 0 {
			return s.processingFoundLexeme(s.shiftFound()), true
		}
	}

	if s.stack.Len() != 0 {
		s.index++
		switch s.stack.Peek().Type() { //nolint:exhaustive // We handle all cases.
		case lexeme.LiteralBegin:
			if s.unfinishedLiteral {
				break
			}
			return s.processingFoundLexeme(lexeme.LiteralEnd), true
		case lexeme.InlineAnnotationBegin:
			return s.processingFoundLexeme(lexeme.InlineAnnotationEnd), true
		case lexeme.InlineAnnotationTextBegin:
			return s.processingFoundLexeme(lexeme.InlineAnnotationTextEnd), true
		}
		err := errors.NewDocumentError(s.file, errors.ErrUnexpectedEOF)
		err.SetIndex(s.dataSize - 1)
		panic(err)
	}

	return lexeme.LexEvent{}, false
}

func (s *scanner) found(lexType lexeme.LexEventType) {
	s.finds = append(s.finds, lexType)
}

func (s *scanner) shiftFound() lexeme.LexEventType {
	length := len(s.finds)
	if length == 0 {
		panic("Empty set of found lexical event")
	}
	lexType := s.finds[0]
	copy(s.finds[0:], s.finds[1:])
	s.finds = s.finds[:length-1]
	return lexType
}

func (s *scanner) processingFoundLexeme(lexType lexeme.LexEventType) lexeme.LexEvent { //nolint:gocyclo // todo try to make this more readable
	i := s.index - 1
	if lexType == lexeme.NewLine || lexType == lexeme.EndTop {
		return lexeme.NewLexEvent(lexType, i, i, s.file)
	}

	if lexType.IsOpening() {
		var lex lexeme.LexEvent
		if lexType == lexeme.InlineAnnotationBegin || lexType == lexeme.MultiLineAnnotationBegin {
			lex = lexeme.NewLexEvent(lexType, i-1, i, s.file) // `//` or `/*`
		} else {
			// `{`, `[`, `"` or literal first character (ex: `1` in `123`).
			lex = lexeme.NewLexEvent(lexType, i, i, s.file)
		}
		s.stack.Push(lex)
		return lex
	}

	// closing tag
	pair := s.stack.Pop()
	pairType := pair.Type()
	if (pairType == lexeme.ObjectBegin && lexType == lexeme.ObjectEnd) ||
		(pairType == lexeme.ArrayBegin && lexType == lexeme.ArrayEnd) {
		return lexeme.NewLexEvent(lexType, pair.Begin(), i, s.file)
	}

	if (pairType == lexeme.LiteralBegin && lexType == lexeme.LiteralEnd) ||
		(pairType == lexeme.ArrayItemBegin && lexType == lexeme.ArrayItemEnd) ||
		(pairType == lexeme.ObjectKeyBegin && lexType == lexeme.ObjectKeyEnd) ||
		(pairType == lexeme.ObjectValueBegin && lexType == lexeme.ObjectValueEnd) {
		return lexeme.NewLexEvent(lexType, pair.Begin(), i-1, s.file)
	}
	panic("Incorrect ending of the lexical event")
}

func stateFoundRootValue(s *scanner, c byte) state {
	r := stateBeginValue(s, c)
	switch r { //nolint:exhaustive // It's okay.
	case scanBeginObject:
		s.found(lexeme.ObjectBegin)

	case scanBeginArray:
		s.found(lexeme.ArrayBegin)

	case scanBeginLiteral:
		s.found(lexeme.LiteralBegin)
	}
	return r
}

func stateFoundObjectKeyBeginOrEmpty(s *scanner, c byte) state {
	if bytes.IsBlank(c) {
		return scanSkipSpace
	}

	return stateBeginKeyOrEmpty(s, c)
}

func stateFoundObjectKeyBegin(s *scanner, c byte) state {
	if bytes.IsBlank(c) {
		return scanSkipSpace
	}

	r := stateBeginString(s, c)
	s.found(lexeme.ObjectKeyBegin)
	return r
}

func stateFoundObjectValueBegin(s *scanner, c byte) state {
	r := stateBeginValue(s, c)
	switch r { //nolint:exhaustive // It's okay.
	case scanBeginLiteral:
		s.found(lexeme.ObjectValueBegin)
		s.found(lexeme.LiteralBegin)

	case scanBeginObject:
		s.found(lexeme.ObjectValueBegin)
		s.found(lexeme.ObjectBegin)

	case scanBeginArray:
		s.found(lexeme.ObjectValueBegin)
		s.found(lexeme.ArrayBegin)
	}
	return r
}

func stateFoundArrayItemBeginOrEmpty(s *scanner, c byte) state {
	r := stateBeginArrayItemOrEmpty(s, c)
	switch r { //nolint:exhaustive // It's okay.
	case scanBeginLiteral:
		s.found(lexeme.ArrayItemBegin)
		s.found(lexeme.LiteralBegin)

	case scanBeginObject:
		s.found(lexeme.ArrayItemBegin)
		s.found(lexeme.ObjectBegin)

	case scanBeginArray:
		s.found(lexeme.ArrayItemBegin)
		s.found(lexeme.ArrayBegin)
	}
	return r
}

func stateFoundArrayItemBegin(s *scanner, c byte) state {
	r := stateBeginValue(s, c)
	switch r { //nolint:exhaustive // It's okay.
	case scanBeginLiteral:
		s.found(lexeme.ArrayItemBegin)
		s.found(lexeme.LiteralBegin)

	case scanBeginObject:
		s.found(lexeme.ArrayItemBegin)
		s.found(lexeme.ObjectBegin)

	case scanBeginArray:
		s.found(lexeme.ArrayItemBegin)
		s.found(lexeme.ArrayBegin)
	}
	return r
}

func stateBeginValue(s *scanner, c byte) state { //nolint:gocyclo // It's okay.
	if bytes.IsBlank(c) {
		return scanSkipSpace
	}
	switch c {
	case '{':
		s.step = stateFoundObjectKeyBeginOrEmpty
		return scanBeginObject
	case '[':
		s.step = stateFoundArrayItemBeginOrEmpty
		return scanBeginArray
	case '"':
		s.step = stateInString
		s.unfinishedLiteral = true
		return scanBeginLiteral
	case '-':
		s.step = stateNeg
		s.unfinishedLiteral = true
		return scanBeginLiteral
	case '0': // beginning of 0.123
		s.step = state0
		return scanBeginLiteral
	case 't': // beginning of true
		s.step = stateT
		s.unfinishedLiteral = true
		return scanBeginLiteral
	case 'f': // beginning of false
		s.step = stateF
		s.unfinishedLiteral = true
		return scanBeginLiteral
	case 'n': // beginning of null
		s.step = stateN
		s.unfinishedLiteral = true
		return scanBeginLiteral
	}
	if '1' <= c && c <= '9' { // beginning of 1234.5
		s.step = state1
		return scanBeginLiteral
	}
	panic(s.newDocumentErrorAtCharacter("looking for beginning of value"))
}

// after reading `[`
func stateBeginArrayItemOrEmpty(s *scanner, c byte) state {
	if c == ']' {
		return stateFoundArrayEnd(s)
	}
	return stateBeginValue(s, c)
}

// after reading `{`
func stateBeginKeyOrEmpty(s *scanner, c byte) state {
	if c == '}' {
		return stateFoundObjectEnd(s)
	}
	s.found(lexeme.ObjectKeyBegin)
	return stateBeginString(s, c)
}

// after reading `{"key": value,`
func stateBeginString(s *scanner, c byte) state {
	if c == '"' {
		s.step = stateInString
		return scanBeginLiteral
	}
	panic(s.newDocumentErrorAtCharacter("looking for beginning of string"))
}

func stateEndValue(s *scanner, c byte) state {
	length := s.stack.Len()

	if length == 0 { // json ex `{} `
		s.step = stateEndTop
		return s.step(s, c)
	}

	t := s.stack.Peek().Type()

	if t == lexeme.LiteralBegin {
		s.found(lexeme.LiteralEnd)

		if length == 1 { // json ex `123 `
			s.step = stateEndTop
			return s.step(s, c)
		}

		t = s.stack.Get(length - 2).Type()
	}

	switch t { //nolint:exhaustive // We will throw a panic in over cases.
	case lexeme.ObjectKeyBegin:
		s.found(lexeme.ObjectKeyEnd)
		s.step = stateAfterObjectKey
		return s.step(s, c)
	case lexeme.ObjectValueBegin:
		s.found(lexeme.ObjectValueEnd)
		s.step = stateAfterObjectValue
		return s.step(s, c)
	case lexeme.ArrayItemBegin:
		s.found(lexeme.ArrayItemEnd)
		s.step = stateAfterArrayItem
		return s.step(s, c)
	}
	panic(s.newDocumentErrorAtCharacter("at the end of value"))
}

func stateAfterObjectKey(s *scanner, c byte) state {
	if bytes.IsBlank(c) {
		return scanSkipSpace
	}

	if c == ':' {
		s.step = stateFoundObjectValueBegin
		return scanObjectKey
	}
	panic(s.newDocumentErrorAtCharacter("after object key"))
}

func stateAfterObjectValue(s *scanner, c byte) state {
	if bytes.IsBlank(c) {
		return scanSkipSpace
	}
	if c == ',' {
		s.step = stateFoundObjectKeyBegin
		return scanObjectValue
	}
	if c == '}' {
		return stateFoundObjectEnd(s)
	}
	panic(s.newDocumentErrorAtCharacter("after object key:value pair"))
}

func stateAfterArrayItem(s *scanner, c byte) state {
	if bytes.IsBlank(c) {
		return scanSkipSpace
	}
	if c == ',' {
		s.step = stateFoundArrayItemBegin
		return scanArrayValue
	}
	if c == ']' {
		return stateFoundArrayEnd(s)
	}
	panic(s.newDocumentErrorAtCharacter("after array item"))
}

func stateFoundObjectEnd(s *scanner) state {
	s.found(lexeme.ObjectEnd)
	s.step = stateEndValue
	return scanEndObject
}

func stateFoundArrayEnd(s *scanner) state {
	s.found(lexeme.ArrayEnd)
	if s.stack.Len() == 0 {
		s.step = stateEndTop
	} else {
		s.step = stateEndValue
	}
	return scanEndArray
}

// stateEndTop is the state after finishing the top-level value,
// such as after reading `{}` or `[1,2,3]`.
// Only space characters should be seen now.
func stateEndTop(s *scanner, c byte) state {
	if !bytes.IsBlank(c) {
		if !s.allowTrailingNonSpaceCharacters {
			panic(s.newDocumentErrorAtCharacter("non-space byte after top-level value"))
		}
		s.found(lexeme.EndTop)
	}
	return scanEnd
}

// after reading `"`
func stateInString(s *scanner, c byte) state {
	switch c {
	case '"':
		s.step = stateEndValue
		s.unfinishedLiteral = false
		return scanContinue
	case '\\':
		s.step = stateInStringEsc
		return scanContinue
	}
	if c < 0x20 {
		panic(s.newDocumentErrorAtCharacter("in string literal"))
	}
	return scanContinue
}

// after reading `"\` during a quoted string
func stateInStringEsc(s *scanner, c byte) state {
	switch c {
	case 'b', 'f', 'n', 'r', 't', '\\', '/', '"':
		s.step = stateInString
		return scanContinue
	case 'u':
		s.returnToStep.Push(stateInString)
		s.step = stateInStringEscU
		return scanContinue
	}
	panic(s.newDocumentErrorAtCharacter("in string escape code"))
}

// after reading `"\u` during a quoted string
func stateInStringEscU(s *scanner, c byte) state {
	if '0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F' {
		s.step = stateInStringEscU1
		return scanContinue
	}
	panic(s.newDocumentErrorAtCharacter("in \\u hexadecimal character escape"))
}

// after reading `"\u1` during a quoted string
func stateInStringEscU1(s *scanner, c byte) state {
	if '0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F' {
		s.step = stateInStringEscU12
		return scanContinue
	}
	panic(s.newDocumentErrorAtCharacter("in \\u hexadecimal character escape"))
}

// after reading `"\u12` during a quoted string
func stateInStringEscU12(s *scanner, c byte) state {
	if '0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F' {
		s.step = stateInStringEscU123
		return scanContinue
	}
	panic(s.newDocumentErrorAtCharacter("in \\u hexadecimal character escape"))
}

// after reading `"\u123` during a quoted string
func stateInStringEscU123(s *scanner, c byte) state {
	if '0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F' {
		s.step = s.returnToStep.Pop() // = stateInString for JSON, = stateInAnnotationObjectKey for AnnotationObject
		return scanContinue
	}
	panic(s.newDocumentErrorAtCharacter("in \\u hexadecimal character escape"))
}

// after reading `-` during a number
func stateNeg(s *scanner, c byte) state {
	if c == '0' {
		s.step = state0
		s.unfinishedLiteral = false
		return scanContinue
	}
	if '1' <= c && c <= '9' {
		s.step = state1
		s.unfinishedLiteral = false
		return scanContinue
	}
	panic(s.newDocumentErrorAtCharacter("in numeric literal"))
}

// after reading a non-zero integer during a number,
// such as after reading `1` or `100` but not `0`
func state1(s *scanner, c byte) state {
	if '0' <= c && c <= '9' {
		s.step = state1
		return scanContinue
	}
	return state0(s, c)
}

// after reading `0` during a number
func state0(s *scanner, c byte) state {
	if c == '.' {
		s.step = stateDot
		return scanContinue
	}
	if c == 'e' || c == 'E' {
		s.step = stateE
		return scanContinue
	}
	return stateEndValue(s, c)
}

// after reading the integer and decimal point in a number, such as after reading `1.`
func stateDot(s *scanner, c byte) state {
	if '0' <= c && c <= '9' {
		s.step = stateDot0
		return scanContinue
	}
	panic(s.newDocumentErrorAtCharacter("after decimal point in numeric literal"))
}

// after reading the integer, decimal point, and subsequent
// digits of a number, such as after reading `3.14`
func stateDot0(s *scanner, c byte) state {
	if '0' <= c && c <= '9' {
		return scanContinue
	}
	if c == 'e' || c == 'E' {
		s.step = stateE
		return scanContinue
	}
	return stateEndValue(s, c)
}

// after reading the mantissa and e in a number,
// such as after reading `314e` or `0.314e`
func stateE(s *scanner, c byte) state {
	if c == '+' || c == '-' {
		s.step = stateESign
		return scanContinue
	}
	return stateESign(s, c)
}

// after reading the mantissa, e, and sign in a number,
// such as after reading `314e-` or `0.314e+`
func stateESign(s *scanner, c byte) state {
	if '0' <= c && c <= '9' {
		s.step = stateE0
		return scanContinue
	}
	panic(s.newDocumentErrorAtCharacter("in exponent of numeric literal"))
}

// after reading the mantissa, e, optional sign,
// and at least one digit of the exponent in a number,
// such as after reading `314e-2` or `0.314e+1` or `3.14e0`
func stateE0(s *scanner, c byte) state {
	if '0' <= c && c <= '9' {
		return scanContinue
	}
	return stateEndValue(s, c)
}

// after reading `t`
func stateT(s *scanner, c byte) state {
	if c == 'r' {
		s.step = stateTr
		return scanContinue
	}
	panic(s.newDocumentErrorAtCharacter("in literal true (expecting 'r')"))
}

// after reading `tr`
func stateTr(s *scanner, c byte) state {
	if c == 'u' {
		s.step = stateTru
		return scanContinue
	}
	panic(s.newDocumentErrorAtCharacter("in literal true (expecting 'u')"))
}

// after reading `tru`
func stateTru(s *scanner, c byte) state {
	if c == 'e' {
		s.step = stateEndValue
		s.unfinishedLiteral = false
		return scanContinue
	}
	panic(s.newDocumentErrorAtCharacter("in literal true (expecting 'e')"))
}

// after reading `f`
func stateF(s *scanner, c byte) state {
	if c == 'a' {
		s.step = stateFa
		return scanContinue
	}
	panic(s.newDocumentErrorAtCharacter("in literal false (expecting 'a')"))
}

// after reading `fa`
func stateFa(s *scanner, c byte) state {
	if c == 'l' {
		s.step = stateFal
		return scanContinue
	}
	panic(s.newDocumentErrorAtCharacter("in literal false (expecting 'l')"))
}

// after reading `fal`
func stateFal(s *scanner, c byte) state {
	if c == 's' {
		s.step = stateFals
		return scanContinue
	}
	panic(s.newDocumentErrorAtCharacter("in literal false (expecting 's')"))
}

// after reading `fals`
func stateFals(s *scanner, c byte) state {
	if c == 'e' {
		s.step = stateEndValue
		s.unfinishedLiteral = false
		return scanContinue
	}
	panic(s.newDocumentErrorAtCharacter("in literal false (expecting 'e')"))
}

// after reading `n`
func stateN(s *scanner, c byte) state {
	if c == 'u' {
		s.step = stateNu
		return scanContinue
	}
	panic(s.newDocumentErrorAtCharacter("in literal null (expecting 'u')"))
}

// after reading `nu`
func stateNu(s *scanner, c byte) state {
	if c == 'l' {
		s.step = stateNul
		return scanContinue
	}
	panic(s.newDocumentErrorAtCharacter("in literal null (expecting 'l')"))
}

// after reading `nul`
func stateNul(s *scanner, c byte) state {
	if c == 'l' {
		s.step = stateEndValue
		s.unfinishedLiteral = false
		return scanContinue
	}
	panic(s.newDocumentErrorAtCharacter("in literal null (expecting 'l')"))
}

func (s *scanner) newDocumentErrorAtCharacter(context string) errors.DocumentError {
	// Make runes (utf8 symbols) from current index to last of slice s.data.
	// Get first rune. Then make string with format ' symbol '
	runes := []rune(string(s.data[(s.index - 1):])) // TODO is memory allocation optimization required?
	e := errors.Format(errors.ErrInvalidCharacter, string(runes[0]), context)
	err := errors.NewDocumentError(s.file, e)
	err.SetIndex(s.index - 1)
	return err
}

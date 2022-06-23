package scanner

import (
	"fmt"

	"github.com/jsightapi/jsight-schema-go-library/bytes"
	"github.com/jsightapi/jsight-schema-go-library/errors"
	"github.com/jsightapi/jsight-schema-go-library/fs"
	"github.com/jsightapi/jsight-schema-go-library/internal/lexeme"
	"github.com/jsightapi/jsight-schema-go-library/internal/scanner"
)

type state uint8
type annotation uint8
type stepFunc func(*Scanner, byte) state

const (
	annotationNone      annotation = iota // not inside the annotation
	annotationInline                      // the pointer inside the inline annotation
	annotationMultiLine                   // the pointer inside the multi-line annotation
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

	// scanBeginTypesShortcut indicates beginning of "TYPE" or "OR" shortcut with
	// user defined types.
	//
	// Examples:
	// {
	//   "foo": @Fizz | @Buzz,
	//   "bar": @Fizz
	// }
	scanBeginTypesShortcut

	// scanKeyShortcutBegin indicates beginning of shortcut with user defined
	// types like:
	// Example
	// {
	//   @cat: "Tom"
	// }
	scanKeyShortcutBegin

	// scanEnd indicates the end of the scanning. Top-level value ended *before*
	// this byte.
	scanEnd
)

// Scanner represents a scanner is a JSchema scanning state machine.
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
type Scanner struct {
	// step is a func to be called to execute the next transition.
	// Also tried using an integer constant and a single func
	// with a switch, but using the func directly was 10% faster
	// on a 64-bit Mac Mini, and it's nicer to read.
	step stepFunc

	// returnToStep a stack of step functions, to preserve the sequence of steps
	// (and return to them) in some cases.
	returnToStep stepFuncStack

	// file a structure containing jSchema data.
	file *fs.File

	// data jSchema content.
	data bytes.Bytes

	// stack a stack of found lexical event. The stack is needed for the scanner
	// to take into account the nesting of SCHEME elements.
	stack scanner.LexemesStack

	// finds a list of found types of lexical event for the current step. Several
	// lexical events can be found in one step (example: ArrayItemBegin and LiteralBegin).
	finds []lexeme.LexEventType

	// prevContextsStack a stack of previous scanner contexts.
	// Used for restoring a previous context after finishing current one.
	prevContextsStack contextStack

	// context indicates which type of entity we process right now.
	context context

	// index scanned byte index.
	index bytes.Index

	// dataSize a size of schema data in bytes. Count once for optimization.
	dataSize bytes.Index

	// annotation one of the possible States of annotation processing (annotationNone,
	// annotationInline, annotationMultiLine).
	annotation annotation

	// unfinishedLiteral a sign that a literal has been started but not completed.
	unfinishedLiteral bool

	// lengthComputing used when a file contains data after the schema (for example,
	// in jApi).
	lengthComputing bool

	// boundary the character of the bounding lines.
	boundary byte

	// allowAnnotation indicates is annotation is allowed or not.
	allowAnnotation bool

	hasTrailingCharacters bool
}

type context struct {
	Type         contextType
	ArrayHasItem bool
}

func newContext(t contextType) context {
	return context{
		Type: t,
	}
}

type contextType int

const (
	contextTypeInitial contextType = iota
	contextTypeObject
	contextTypeArray
	contextTypeShortcut
)

func New(file *fs.File, oo ...Option) *Scanner {
	content := file.Content()

	s := &Scanner{
		step:            stateFoundRootValue,
		file:            file,
		data:            content,
		dataSize:        bytes.Index(len(content)),
		returnToStep:    make(stepFuncStack, 0, 2),
		stack:           scanner.NewLexemesStack(),
		finds:           make([]lexeme.LexEventType, 0, 3),
		context:         newContext(contextTypeInitial),
		allowAnnotation: true,
	}

	for _, o := range oo {
		o(s)
	}

	return s
}

type Option func(*Scanner)

// ComputeLength switch scanner in length computing mode.
// Scanner in this mode shouldn't be used for parsing.
func ComputeLength(s *Scanner) {
	s.lengthComputing = true
}

func (s *Scanner) Length() uint {
	if !s.lengthComputing {
		panic("Method not allowed")
	}
	var length uint
	for {
		lex, ok := s.Next()
		if !ok {
			break
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
		c := s.data[length-1]
		if !bytes.IsSpace(c) {
			break
		}
	}
	return length
}

func (s *Scanner) newDocumentError(code errors.ErrorCode, c byte) errors.DocumentError {
	e := errors.Format(code, scanner.QuoteChar(c))
	err := errors.NewDocumentError(s.file, e)
	err.SetIndex(s.index - 1)
	return err
}

func (s *Scanner) newDocumentErrorAtCharacter(context string) errors.DocumentError {
	// Make runes (utf8 symbols) from current index to last of slice s.data.
	// Get first rune. Then make string with format ' symbol '
	runes := []rune(string(s.data[(s.index - 1):])) // TODO is memory allocation optimization required?
	e := errors.Format(errors.ErrInvalidCharacter, string(runes[0]), context)
	err := errors.NewDocumentError(s.file, e)
	err.SetIndex(s.index - 1)
	return err
}

// Next reads schema byte by byte.
// Panic if an invalid jSchema structure is found.
// Stops if it detects lexical events.
// Returns pointer to found lexeme event, or nil if you have complete reading.
func (s *Scanner) Next() (lexeme.LexEvent, bool) {
	if len(s.finds) != 0 {
		return s.processingFoundLexeme(s.shiftFound()), true
	}

	for s.index < s.dataSize {
		c := s.data[s.index]
		s.index++

		// useful for debugging comment below 1 line for release
		// fmt.Printf("Schema-Next->step %s %c\n", runtime.FuncForPC(reflect.ValueOf(s.step).Pointer()).Name(), c)

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
		case lexeme.TypesShortcutBegin:
			s.found(lexeme.MixedValueEnd)
			return s.processingFoundLexeme(lexeme.TypesShortcutEnd), true
		}
		err := errors.NewDocumentError(s.file, errors.ErrUnexpectedEOF)
		err.SetIndex(s.dataSize - 1)
		panic(err)
	}

	return lexeme.LexEvent{}, false
}

func (s *Scanner) isFoundLastObjectEndOnAnnotation() (bool, lexeme.LexEventType) { //nolint:gocyclo // todo try to make this more readable
	length := s.stack.Len()

	switch {
	case length >= 5 &&
		s.stack.Get(length-1).Type() == lexeme.TypesShortcutBegin &&
		s.stack.Get(length-2).Type() == lexeme.MixedValueBegin &&
		s.stack.Get(length-3).Type() == lexeme.ObjectValueBegin &&
		s.stack.Get(length-4).Type() == lexeme.ObjectBegin &&
		(s.stack.Get(length-5).Type() == lexeme.InlineAnnotationBegin || s.stack.Get(length-5).Type() == lexeme.MultiLineAnnotationBegin):
		return true, s.stack.Get(length - 5).Type()

	case length >= 4 &&
		s.stack.Get(length-1).Type() == lexeme.LiteralBegin &&
		s.stack.Get(length-2).Type() == lexeme.ObjectValueBegin &&
		s.stack.Get(length-3).Type() == lexeme.ObjectBegin &&
		(s.stack.Get(length-4).Type() == lexeme.InlineAnnotationBegin || s.stack.Get(length-4).Type() == lexeme.MultiLineAnnotationBegin):
		return true, s.stack.Get(length - 4).Type()

	case length >= 3 &&
		s.stack.Get(length-1).Type() == lexeme.ObjectValueBegin &&
		s.stack.Get(length-2).Type() == lexeme.ObjectBegin &&
		(s.stack.Get(length-3).Type() == lexeme.InlineAnnotationBegin || s.stack.Get(length-3).Type() == lexeme.MultiLineAnnotationBegin):
		return true, s.stack.Get(length - 3).Type()

	case length >= 2 &&
		s.stack.Get(length-1).Type() == lexeme.ObjectBegin &&
		(s.stack.Get(length-2).Type() == lexeme.InlineAnnotationBegin || s.stack.Get(length-2).Type() == lexeme.MultiLineAnnotationBegin):
		return true, s.stack.Get(length - 2).Type()
	}
	return false, lexeme.InlineAnnotationBegin
}

func (s *Scanner) isInsideMultiLineAnnotation() bool {
	for i := s.stack.Len() - 1; i >= 0; i-- {
		if s.stack.Get(i).Type() == lexeme.MultiLineAnnotationBegin {
			return true
		}
	}
	return false
}

func (s *Scanner) found(lexType lexeme.LexEventType) {
	s.finds = append(s.finds, lexType)
}

func (s *Scanner) shiftFound() lexeme.LexEventType {
	length := len(s.finds)
	if length == 0 {
		panic("Empty set of found lexical event")
	}
	lexType := s.finds[0]
	copy(s.finds[0:], s.finds[1:])
	s.finds = s.finds[:length-1]
	return lexType
}

func (s *Scanner) processingFoundLexeme(lexType lexeme.LexEventType) lexeme.LexEvent { //nolint:gocyclo // todo try to make this more readable
	i := s.index - 1
	if lexType == lexeme.NewLine || lexType == lexeme.EndTop { //nolint:gocritic // todo rewrite this logic to switch
		return lexeme.NewLexEvent(lexType, i, i, s.file)
	} else if lexType.IsOpening() {
		var lex lexeme.LexEvent
		if lexType == lexeme.InlineAnnotationBegin || lexType == lexeme.MultiLineAnnotationBegin {
			lex = lexeme.NewLexEvent(lexType, i-1, i, s.file) // `//` or `/*`
		} else {
			// `{`, `[`, `"` or literal first character (ex: `1` in `123`).
			lex = lexeme.NewLexEvent(lexType, i, i, s.file)
		}
		s.stack.Push(lex)
		return lex
	} else { // closing tag
		pair := s.stack.Pop()
		pairType := pair.Type()
		if (pairType == lexeme.ObjectBegin && lexType == lexeme.ObjectEnd) ||
			(pairType == lexeme.ArrayBegin && lexType == lexeme.ArrayEnd) ||
			(pairType == lexeme.MultiLineAnnotationBegin && lexType == lexeme.MultiLineAnnotationEnd) {
			return lexeme.NewLexEvent(lexType, pair.Begin(), i, s.file)
		} else if (pairType == lexeme.LiteralBegin && lexType == lexeme.LiteralEnd) ||
			(pairType == lexeme.ArrayItemBegin && lexType == lexeme.ArrayItemEnd) ||
			(pairType == lexeme.ObjectKeyBegin && lexType == lexeme.ObjectKeyEnd) ||
			(pairType == lexeme.ObjectValueBegin && lexType == lexeme.ObjectValueEnd) ||
			(pairType == lexeme.InlineAnnotationTextBegin && lexType == lexeme.InlineAnnotationTextEnd) ||
			(pairType == lexeme.MultiLineAnnotationTextBegin && lexType == lexeme.MultiLineAnnotationTextEnd) ||
			(pairType == lexeme.InlineAnnotationBegin && lexType == lexeme.InlineAnnotationEnd) ||
			(pairType == lexeme.KeyShortcutBegin && lexType == lexeme.KeyShortcutEnd) ||
			(pairType == lexeme.TypesShortcutBegin && lexType == lexeme.TypesShortcutEnd) ||
			(pairType == lexeme.MixedValueBegin && lexType == lexeme.MixedValueEnd) {
			if lexType == lexeme.MixedValueEnd && s.data[i-1] == ' ' {
				i--
			}
			return lexeme.NewLexEvent(lexType, pair.Begin(), i-1, s.file)
		}
	}
	panic("Incorrect ending of the lexical event")
}

func (s *Scanner) isNewLine(c byte) bool {
	if c != '\n' && c != '\r' {
		return false
	}

	if s.annotation == annotationInline {
		panic(s.newDocumentErrorAtCharacter("inside inline annotation"))
	}
	return true
}

func (*Scanner) isAnnotationStart(c byte) bool {
	return c == '/'
}

func (s *Scanner) switchToAnnotation() {
	if !s.allowAnnotation {
		err := errors.NewDocumentError(s.file, errors.Format(errors.ErrAnnotationNotAllowed))
		err.SetIndex(s.index - 1)
		panic(err)
	}

	switch s.annotation {
	case annotationNone:
		s.returnToStep.Push(s.step)
		s.step = stateAnyAnnotationStart
	case annotationMultiLine:
		s.returnToStep.Push(s.step)
		s.step = stateInlineAnnotationStart
	default:
		// panic("Incorrect annotation mode")
		panic(s.newDocumentErrorAtCharacter("inside inline annotation"))
	}
}

func (s *Scanner) setContext(c context) {
	s.prevContextsStack.Push(s.context)
	s.context = c
}

func (s *Scanner) restoreContext() {
	s.context = s.prevContextsStack.Pop()
}

func stateFoundRootValue(s *Scanner, c byte) state {
	if s.isAnnotationStart(c) {
		s.switchToAnnotation()
		return scanContinue
	}
	if s.isCommentStart(c) {
		s.switchToComment()
		return scanContinue
	}

	r := stateBeginValue(s, c)
	switch r { //nolint:exhaustive // It's okay.
	case scanBeginObject:
		s.found(lexeme.ObjectBegin)
		s.setContext(newContext(contextTypeObject))

	case scanBeginArray:
		s.found(lexeme.ArrayBegin)
		s.setContext(newContext(contextTypeArray))

	case scanBeginLiteral:
		s.found(lexeme.LiteralBegin)

	case scanBeginTypesShortcut:
		s.found(lexeme.MixedValueBegin)
		s.found(lexeme.TypesShortcutBegin)
		s.setContext(newContext(contextTypeShortcut))
	}
	return r
}

func stateFoundObjectKeyBeginOrEmpty(s *Scanner, c byte) state {
	if s.isNewLine(c) {
		s.found(lexeme.NewLine)
		return scanSkipSpace
	}
	if bytes.IsSpace(c) {
		return scanSkipSpace
	}
	if s.isAnnotationStart(c) {
		s.switchToAnnotation()
		return scanContinue
	}
	if s.isCommentStart(c) {
		s.switchToComment()
		return scanContinue
	}
	if c == '@' {
		return beginKeyShortcut(s)
	}

	var r state
	if s.annotation == annotationNone {
		r = stateBeginKeyOrEmpty(s, c)
	} else {
		r = stateBeginAnnotationObjectKeyOrEmpty(s, c)
	}
	return r
}

func stateFoundObjectKeyBegin(s *Scanner, c byte) state {
	if s.isNewLine(c) {
		s.found(lexeme.NewLine)
		if s.annotation == annotationNone {
			s.allowAnnotation = true
		}
		s.step = stateFoundObjectKeyBeginAfterNewLine
		return scanSkipSpace
	}
	if bytes.IsSpace(c) {
		return scanSkipSpace
	}
	if s.isAnnotationStart(c) {
		s.switchToAnnotation()
		return scanContinue
	}
	if s.isCommentStart(c) {
		s.switchToComment()
		return scanContinue
	}
	if c == '@' {
		return beginKeyShortcut(s)
	}

	var r state
	if s.annotation == annotationNone {
		r = stateBeginString(s, c)
		s.found(lexeme.ObjectKeyBegin)
	} else {
		// ...OrEmpty because a comma before the closing parenthesis is allowed. Ex: {k:1,}
		r = stateBeginAnnotationObjectKeyOrEmpty(s, c)
	}
	return r
}

func stateFoundObjectKeyBeginAfterNewLine(s *Scanner, c byte) state {
	if s.isNewLine(c) {
		s.found(lexeme.NewLine)
		return scanSkipSpace
	}
	if bytes.IsSpace(c) {
		return scanSkipSpace
	}
	if s.isCommentStart(c) {
		s.switchToComment()
		return scanContinue
	}
	if c == '@' {
		return beginKeyShortcut(s)
	}

	var r state
	if s.annotation == annotationNone {
		r = stateBeginString(s, c)
		s.found(lexeme.ObjectKeyBegin)
	} else {
		// ...OrEmpty because a comma before the closing parenthesis is allowed. Ex: {k:1,}
		r = stateBeginAnnotationObjectKeyOrEmpty(s, c)
	}
	return r
}

func stateFoundObjectValueBegin(s *Scanner, c byte) state {
	r := stateBeginValue(s, c)
	switch r { //nolint:exhaustive // It's okay.
	case scanBeginLiteral:
		s.found(lexeme.ObjectValueBegin)
		s.found(lexeme.LiteralBegin)

	case scanBeginObject:
		s.found(lexeme.ObjectValueBegin)
		s.found(lexeme.ObjectBegin)
		s.setContext(newContext(contextTypeObject))

	case scanBeginArray:
		s.found(lexeme.ObjectValueBegin)
		s.found(lexeme.ArrayBegin)
		s.setContext(newContext(contextTypeArray))

	case scanBeginTypesShortcut:
		s.found(lexeme.ObjectValueBegin)
		s.found(lexeme.MixedValueBegin)
		s.found(lexeme.TypesShortcutBegin)
	}
	return r
}

func stateFoundArrayItemBeginOrEmpty(s *Scanner, c byte) state {
	if s.isNewLine(c) {
		s.found(lexeme.NewLine)
		return scanContinue
	}
	if s.isCommentStart(c) {
		s.switchToComment()
		return scanContinue
	}

	r := stateBeginArrayItemOrEmpty(s, c)
	switch r { //nolint:exhaustive // It's okay.
	case scanBeginLiteral:
		s.found(lexeme.ArrayItemBegin)
		s.found(lexeme.LiteralBegin)

	case scanBeginObject:
		s.found(lexeme.ArrayItemBegin)
		s.found(lexeme.ObjectBegin)
		s.setContext(newContext(contextTypeObject))

	case scanBeginArray:
		s.found(lexeme.ArrayItemBegin)
		s.found(lexeme.ArrayBegin)
		s.setContext(newContext(contextTypeArray))

	case scanBeginTypesShortcut:
		s.found(lexeme.ArrayItemBegin)
		s.found(lexeme.MixedValueBegin)
		s.found(lexeme.TypesShortcutBegin)
	}
	return r
}

func stateFoundArrayItemBegin(s *Scanner, c byte) state {
	if s.isCommentStart(c) {
		s.switchToComment()
		return scanContinue
	}

	r := stateBeginValue(s, c)
	switch r { //nolint:exhaustive // It's okay.
	case scanBeginLiteral:
		s.found(lexeme.ArrayItemBegin)
		s.found(lexeme.LiteralBegin)

	case scanBeginObject:
		s.found(lexeme.ArrayItemBegin)
		s.found(lexeme.ObjectBegin)
		s.setContext(newContext(contextTypeObject))

	case scanBeginArray:
		s.found(lexeme.ArrayItemBegin)
		s.found(lexeme.ArrayBegin)
		s.setContext(newContext(contextTypeArray))

	case scanBeginTypesShortcut:
		s.found(lexeme.ArrayItemBegin)
		s.found(lexeme.MixedValueBegin)
		s.found(lexeme.TypesShortcutBegin)
	}
	return r
}

func beginKeyShortcut(s *Scanner) state {
	if s.annotation != annotationNone {
		panic(s.newDocumentErrorAtCharacter("key shortcut not allowed in annotation"))
	}
	s.found(lexeme.KeyShortcutBegin)
	s.step = stateKeyShortcut
	return scanKeyShortcutBegin
}

func stateBeginValue(s *Scanner, c byte) state { //nolint:gocyclo // It's okay.
	if s.isNewLine(c) {
		s.found(lexeme.NewLine)
		return scanSkipSpace
	}
	if bytes.IsSpace(c) {
		return scanSkipSpace
	}
	if s.isAnnotationStart(c) {
		s.switchToAnnotation()
		return scanContinue
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
	case '@': // beginning of OR shortcut
		s.step = stateTypesShortcutBeginOfSchemaName
		s.unfinishedLiteral = true
		return scanBeginTypesShortcut
	}
	if '1' <= c && c <= '9' { // beginning of 1234.5
		s.step = state1
		return scanBeginLiteral
	}
	panic(s.newDocumentErrorAtCharacter("looking for beginning of value"))
}

// after reading `[`
func stateBeginArrayItemOrEmpty(s *Scanner, c byte) state {
	if c == ']' {
		return stateFoundArrayEnd(s)
	}
	if s.annotation == annotationNone {
		s.context.ArrayHasItem = true
	}
	return stateBeginValue(s, c)
}

// after reading `{`
func stateBeginKeyOrEmpty(s *Scanner, c byte) state {
	if s.annotation == annotationNone {
		s.allowAnnotation = true
	}
	if c == '}' {
		return stateFoundObjectEnd(s)
	}
	s.found(lexeme.ObjectKeyBegin)
	return stateBeginString(s, c)
}

// after reading `{"key": value,`
func stateBeginString(s *Scanner, c byte) state {
	if c == '"' {
		s.step = stateInString
		return scanBeginLiteral
	}
	panic(s.newDocumentErrorAtCharacter("looking for beginning of string"))
}

func stateEndValue(s *Scanner, c byte) state { //nolint:gocyclo // Pretty readable though.
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
	case lexeme.KeyShortcutBegin:
		s.found(lexeme.KeyShortcutEnd)
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
	case lexeme.TypesShortcutBegin:
		finishShortcut(s)
		return s.step(s, c)
	}
	if s.lengthComputing && t == lexeme.InlineAnnotationBegin {
		s.annotation = annotationNone
		_ = s.stack.Pop()
		s.step = s.returnToStep.Pop()
		return s.step(s, c)
	}
	panic(s.newDocumentErrorAtCharacter("at the end of value"))
}

func finishShortcut(s *Scanner) {
	s.found(lexeme.TypesShortcutEnd)
	switch s.context.Type {
	case contextTypeObject:
		s.found(lexeme.MixedValueEnd)
		s.found(lexeme.ObjectValueEnd)
		s.step = stateAfterObjectValue

	case contextTypeArray:
		s.found(lexeme.MixedValueEnd)
		s.found(lexeme.ArrayItemEnd)
		s.step = stateAfterArrayItem

	case contextTypeShortcut:
		s.found(lexeme.MixedValueEnd)
		s.step = stateEndTop
		s.restoreContext()

	default:
		panic(fmt.Sprintf("Unexpected context %q", s.context.Type))
	}
}

func stateAfterObjectKey(s *Scanner, c byte) state {
	if s.isNewLine(c) {
		s.found(lexeme.NewLine)
	}
	if bytes.IsSpace(c) {
		return scanSkipSpace
	}
	if s.isAnnotationStart(c) {
		s.switchToAnnotation()
		return scanContinue
	}

	if c == ':' {
		s.step = stateFoundObjectValueBegin
		return scanObjectKey
	}
	panic(s.newDocumentErrorAtCharacter("after object key"))
}

func stateAfterObjectValue(s *Scanner, c byte) state {
	if s.isNewLine(c) {
		s.found(lexeme.NewLine)
		return scanSkipSpace
	}
	if bytes.IsSpace(c) {
		return scanSkipSpace
	}
	if s.isAnnotationStart(c) {
		s.switchToAnnotation()
		return scanContinue
	}
	if s.isCommentStart(c) {
		s.switchToComment()
		return scanContinue
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

func stateAfterArrayItem(s *Scanner, c byte) state {
	if s.isNewLine(c) {
		s.found(lexeme.NewLine)
		return scanSkipSpace
	}
	if bytes.IsSpace(c) {
		return scanSkipSpace
	}
	if s.isAnnotationStart(c) {
		s.switchToAnnotation()
		return scanContinue
	}
	if s.isCommentStart(c) {
		s.switchToComment()
		return scanContinue
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

func stateFoundObjectEnd(s *Scanner) state {
	s.found(lexeme.ObjectEnd)
	s.restoreContext()
	s.step = stateEndValue
	if s.annotation != annotationNone {
		if ok, annotationType := s.isFoundLastObjectEndOnAnnotation(); ok {
			switch annotationType {
			case lexeme.InlineAnnotationBegin:
				s.step = stateInlineAnnotationTextPrefix
			case lexeme.MultiLineAnnotationBegin:
				s.step = stateMultiLineAnnotationEndAfterObject
			default:
				panic("Incorrect annotation begin in stack")
			}
		}
	}
	return scanEndObject
}

func stateFoundArrayEnd(s *Scanner) state {
	if s.annotation == annotationNone {
		s.allowAnnotation = !s.context.ArrayHasItem
	}
	s.found(lexeme.ArrayEnd)
	s.restoreContext()
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
func stateEndTop(s *Scanner, c byte) state {
	switch {
	case s.isNewLine(c):
		s.found(lexeme.NewLine)
		return scanSkipSpace

	case s.isAnnotationStart(c):
		s.switchToAnnotation()
		return scanContinue

	case s.isCommentStart(c):
		s.switchToComment()
		return scanContinue

	case !bytes.IsSpace(c):
		if s.lengthComputing {
			if s.stack.Len() > 0 {
				// Looks like we have invalid schema, and we should keep scanning.
				s.hasTrailingCharacters = true
				return scanContinue
			}
			s.found(lexeme.EndTop)
			return scanEnd
		} else if s.annotation == annotationNone {
			panic(s.newDocumentErrorAtCharacter("non-space byte after top-level value"))
		}
	}

	if s.hasTrailingCharacters {
		s.found(lexeme.EndTop)
	}
	return scanEnd
}

// after reading `"`
func stateInString(s *Scanner, c byte) state {
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
func stateInStringEsc(s *Scanner, c byte) state {
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
func stateInStringEscU(s *Scanner, c byte) state {
	if '0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F' {
		s.step = stateInStringEscU1
		return scanContinue
	}
	panic(s.newDocumentErrorAtCharacter("in \\u hexadecimal character escape"))
}

// after reading `"\u1` during a quoted string
func stateInStringEscU1(s *Scanner, c byte) state {
	if '0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F' {
		s.step = stateInStringEscU12
		return scanContinue
	}
	panic(s.newDocumentErrorAtCharacter("in \\u hexadecimal character escape"))
}

// after reading `"\u12` during a quoted string
func stateInStringEscU12(s *Scanner, c byte) state {
	if '0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F' {
		s.step = stateInStringEscU123
		return scanContinue
	}
	panic(s.newDocumentErrorAtCharacter("in \\u hexadecimal character escape"))
}

// after reading `"\u123` during a quoted string
func stateInStringEscU123(s *Scanner, c byte) state {
	if '0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F' {
		s.step = s.returnToStep.Pop() // = stateInAnnotationObjectKey for AnnotationObject
		return scanContinue
	}
	panic(s.newDocumentErrorAtCharacter("in \\u hexadecimal character escape"))
}

// after reading `-` during a number
func stateNeg(s *Scanner, c byte) state {
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
func state1(s *Scanner, c byte) state {
	if '0' <= c && c <= '9' {
		s.step = state1
		return scanContinue
	}
	return state0(s, c)
}

// after reading `0` during a number
func state0(s *Scanner, c byte) state {
	if c == '.' {
		s.unfinishedLiteral = true
		s.step = stateDot
		return scanContinue
	}
	if c == 'e' || c == 'E' {
		panic(s.newDocumentErrorAtCharacter(messageEIsNotAllowed))
	}
	return stateEndValue(s, c)
}

// after reading the integer and decimal point in a number, such as after reading `1.`
func stateDot(s *Scanner, c byte) state {
	if '0' <= c && c <= '9' {
		s.unfinishedLiteral = false
		s.step = stateDot0
		return scanContinue
	}
	panic(s.newDocumentErrorAtCharacter("after decimal point in numeric literal"))
}

// after reading the integer, decimal point, and subsequent
// digits of a number, such as after reading `3.14`
func stateDot0(s *Scanner, c byte) state {
	if '0' <= c && c <= '9' {
		return scanContinue
	}
	if c == 'e' || c == 'E' {
		panic(s.newDocumentErrorAtCharacter(messageEIsNotAllowed))
	}
	return stateEndValue(s, c)
}

// after reading `t`
func stateT(s *Scanner, c byte) state {
	if c == 'r' {
		s.step = stateTr
		return scanContinue
	}
	panic(s.newDocumentErrorAtCharacter("in literal true (expecting 'r')"))
}

// after reading `tr`
func stateTr(s *Scanner, c byte) state {
	if c == 'u' {
		s.step = stateTru
		return scanContinue
	}
	panic(s.newDocumentErrorAtCharacter("in literal true (expecting 'u')"))
}

// after reading `tru`
func stateTru(s *Scanner, c byte) state {
	if c == 'e' {
		s.step = stateEndValue
		s.unfinishedLiteral = false
		return scanContinue
	}
	panic(s.newDocumentErrorAtCharacter("in literal true (expecting 'e')"))
}

// after reading `f`
func stateF(s *Scanner, c byte) state {
	if c == 'a' {
		s.step = stateFa
		return scanContinue
	}
	panic(s.newDocumentErrorAtCharacter("in literal false (expecting 'a')"))
}

// after reading `fa`
func stateFa(s *Scanner, c byte) state {
	if c == 'l' {
		s.step = stateFal
		return scanContinue
	}
	panic(s.newDocumentErrorAtCharacter("in literal false (expecting 'l')"))
}

// after reading `fal`
func stateFal(s *Scanner, c byte) state {
	if c == 's' {
		s.step = stateFals
		return scanContinue
	}
	panic(s.newDocumentErrorAtCharacter("in literal false (expecting 's')"))
}

// after reading `fals`
func stateFals(s *Scanner, c byte) state {
	if c == 'e' {
		s.step = stateEndValue
		s.unfinishedLiteral = false
		return scanContinue
	}
	panic(s.newDocumentErrorAtCharacter("in literal false (expecting 'e')"))
}

// after reading `n`
func stateN(s *Scanner, c byte) state {
	if c == 'u' {
		s.step = stateNu
		return scanContinue
	}
	panic(s.newDocumentErrorAtCharacter("in literal null (expecting 'u')"))
}

// after reading `nu`
func stateNu(s *Scanner, c byte) state {
	if c == 'l' {
		s.step = stateNul
		return scanContinue
	}
	panic(s.newDocumentErrorAtCharacter("in literal null (expecting 'l')"))
}

// after reading `nul`
func stateNul(s *Scanner, c byte) state {
	if c == 'l' {
		s.step = stateEndValue
		s.unfinishedLiteral = false
		return scanContinue
	}
	panic(s.newDocumentErrorAtCharacter("in literal null (expecting 'l')"))
}

func stateTypesShortcutBeginOfSchemaName(s *Scanner, c byte) state {
	if bytes.IsValidUserTypeNameByte(c) {
		s.step = stateTypesShortcutSchemaName
		return scanContinue
	}
	panic(s.newDocumentErrorAtCharacter("in schema name"))
}

func stateTypesShortcutSchemaName(s *Scanner, c byte) state {
	if s.isAnnotationStart(c) {
		finishShortcut(s)
		s.switchToAnnotation()
		return scanContinue
	}

	if s.isCommentStart(c) {
		finishShortcut(s)
		s.switchToComment()
		return scanContinue
	}

	switch {
	case bytes.IsValidUserTypeNameByte(c):
		s.step = stateTypesShortcutSchemaName

	case c == ' ' || c == '\t':
		s.step = stateTypesShortcutBeforePipe

	case c == '|':
		s.step = stateTypesShortcutAfterPipe

	default:
		return stateEndValue(s, c)
	}
	return scanContinue
}

func stateTypesShortcutBeforePipe(s *Scanner, c byte) state {
	if s.isAnnotationStart(c) {
		finishShortcut(s)
		s.switchToAnnotation()
		return scanContinue
	}

	if s.isCommentStart(c) {
		finishShortcut(s)
		s.switchToComment()
		return scanContinue
	}

	switch {
	case c == ' ' || c == '\t':
		s.step = stateTypesShortcutBeforePipe

	case c == '|':
		s.step = stateTypesShortcutAfterPipe

	default:
		s.step = stateEndValue
		s.unfinishedLiteral = false
		return s.step(s, c)
	}
	return scanContinue
}

func stateTypesShortcutAfterPipe(s *Scanner, c byte) state {
	switch c {
	case ' ', '\t':
		s.step = stateTypesShortcutAfterPipe

	case '@':
		s.step = stateTypesShortcutBeginOfSchemaName

	default:
		panic(s.newDocumentErrorAtCharacter("expects ' ', '\\t', or '@'"))
	}
	return scanContinue
}

func stateAnyAnnotationStart(s *Scanner, c byte) state {
	switch c {
	case '/': // second slash - inline annotation
		s.annotation = annotationInline
		s.found(lexeme.InlineAnnotationBegin)
		s.step = stateInlineAnnotation
		return scanContinue
	case '*': // multi-line annotation
		s.annotation = annotationMultiLine
		s.found(lexeme.MultiLineAnnotationBegin)
		s.step = stateMultiLineAnnotation
		return scanContinue
	}
	panic(s.newDocumentErrorAtCharacter("after first slash"))
}

func stateInlineAnnotationStart(s *Scanner, c byte) state {
	if c == '/' { // second slash - inline annotation
		s.annotation = annotationInline
		s.found(lexeme.InlineAnnotationBegin)
		s.step = stateInlineAnnotation
		return scanContinue
	}
	panic(s.newDocumentErrorAtCharacter("after first slash on start inline annotation"))
}

func stateMultiLineAnnotation(s *Scanner, c byte) state {
	if s.isNewLine(c) {
		s.found(lexeme.NewLine)
		return scanSkipSpace
	}
	if bytes.IsSpace(c) {
		return scanSkipSpace
	}
	if c == '{' {
		return stateFoundRootValue(s, c)
	}
	s.found(lexeme.MultiLineAnnotationTextBegin)
	s.step = stateMultiLineAnnotationText
	return s.step(s, c)
}

func stateMultiLineAnnotationEndAfterObject(s *Scanner, c byte) state {
	if s.isNewLine(c) {
		s.found(lexeme.NewLine)
		return scanSkipSpace
	}
	if bytes.IsSpace(c) {
		return scanContinue
	}
	if c == '*' {
		s.step = stateMultiLineAnnotationEnd
		return scanContinue
	} else {
		// proceed annotation
		s.step = stateMultiLineAnnotation
		return scanContinue
	}
}

func stateMultiLineAnnotationEnd(s *Scanner, c byte) state {
	if c == '/' { // after *
		s.annotation = annotationNone
		s.found(lexeme.MultiLineAnnotationEnd)
		s.step = s.returnToStep.Pop()
		return scanContinue
	}
	panic(s.newDocumentErrorAtCharacter("in multi-line annotation after \"*\" character"))
}

func stateInlineAnnotation(s *Scanner, c byte) state {
	switch c {
	case ' ', '\t':
		return scanContinue

	case '{':
		return stateFoundRootValue(s, c)
	}

	s.found(lexeme.InlineAnnotationTextBegin)
	s.step = stateInlineAnnotationText
	return s.step(s, c)
}

func stateInlineAnnotationTextPrefix(s *Scanner, c byte) state {
	if c == ' ' || c == '\t' {
		return scanSkipSpace
	}
	if c == '\n' || c == '\r' {
		s.found(lexeme.InlineAnnotationEnd)
		s.found(lexeme.NewLine)
		s.step = s.returnToStep.Pop()
		if s.isInsideMultiLineAnnotation() {
			s.annotation = annotationMultiLine
		} else {
			s.annotation = annotationNone
		}
		return scanContinue
	}
	if s.isCommentStart(c) {
		s.switchToComment()
		return scanContinue
	}
	if c == '-' {
		s.step = stateInlineAnnotationTextPrefix2
		return scanContinue
	}
	panic(s.newDocumentErrorAtCharacter("after object in inline annotation"))
}

func stateInlineAnnotationTextPrefix2(s *Scanner, c byte) state {
	if c == ' ' || c == '\t' {
		return scanContinue
	}
	s.found(lexeme.InlineAnnotationTextBegin)
	s.step = stateInlineAnnotationText
	return s.step(s, c)
}

func stateInlineAnnotationText(s *Scanner, c byte) state {
	switch c {
	case '\n', '\r':
		s.found(lexeme.InlineAnnotationTextEnd)
		s.found(lexeme.InlineAnnotationEnd)
		s.found(lexeme.NewLine)
		fn := s.returnToStep.Pop()
		s.step = func(s *Scanner, c byte) state {
			if s.isAnnotationStart(c) {
				panic(s.newDocumentErrorAtCharacter("after inline annotation"))
			}
			return fn(s, c)
		}
		if s.isInsideMultiLineAnnotation() {
			s.annotation = annotationMultiLine
		} else {
			s.annotation = annotationNone
		}

	case '#':
		if !s.isInsideMultiLineAnnotation() {
			s.found(lexeme.InlineAnnotationTextEnd)
			s.found(lexeme.InlineAnnotationEnd)
			s.step = stateInlineAnnotationTextSkip
		}
	}
	return scanContinue
}

func stateInlineAnnotationTextSkip(s *Scanner, c byte) state {
	if c != '\n' && c != '\r' {
		return scanContinue
	}

	s.found(lexeme.NewLine)
	fn := s.returnToStep.Pop()
	s.step = func(s *Scanner, c byte) state {
		if s.isAnnotationStart(c) {
			panic(s.newDocumentErrorAtCharacter("after inline annotation"))
		}
		return fn(s, c)
	}
	if s.isInsideMultiLineAnnotation() {
		s.annotation = annotationMultiLine
	} else {
		s.annotation = annotationNone
	}
	return scanContinue
}

func stateBeginAnnotationObjectKeyOrEmpty(s *Scanner, c byte) state {
	if c == '}' {
		return stateFoundObjectEnd(s)
	}
	s.found(lexeme.ObjectKeyBegin)
	return stateBeginAnnotationObjectKey(s, c)
}

func stateBeginAnnotationObjectKey(s *Scanner, c byte) state {
	if c == '"' {
		s.boundary = '"'
		s.step = stateInString
		return scanBeginLiteral
	}

	s.boundary = 0 // default value
	s.step = stateInAnnotationObjectKeyFirstLetter
	return s.step(s, c)
}

func stateInAnnotationObjectKeyFirstLetter(s *Scanner, c byte) state {
	if (s.boundary == 0 && (c == ':' || c == '\n' || c == '\r' || c == '\\')) || c == s.boundary || c < 0x20 {
		panic(s.newDocumentError(errors.ErrInvalidCharacterInAnnotationObjectKey, c))
	}
	s.step = stateInAnnotationObjectKey
	return scanContinue
}

func stateInAnnotationObjectKey(s *Scanner, c byte) state {
	switch {
	case s.boundary == 0 && c == ':':
		return stateEndValue(s, c)

	case c == s.boundary:
		s.step = stateEndValue

	case c == ' ':
		s.step = stateInAnnotationObjectKeyAfter

	case c < 0x20 || (c == '"' || c == '\n' || c == '\r'):
		panic(s.newDocumentError(errors.ErrInvalidCharacterInAnnotationObjectKey, c))
	}
	return scanContinue
}

func stateInAnnotationObjectKeyAfter(s *Scanner, c byte) state {
	switch {
	case s.boundary == 0 && c == ':':
		return stateEndValue(s, c)

	case c == s.boundary:
		s.step = stateEndValue

	case c == ' ':
		return scanContinue
	}
	panic(s.newDocumentError(errors.ErrInvalidCharacterInAnnotationObjectKey, c))
}

func stateMultiLineAnnotationText(s *Scanner, c byte) state {
	if c == '*' && s.data[s.index] == '/' {
		s.found(lexeme.MultiLineAnnotationTextEnd)
		s.step = stateMultiLineAnnotationEnd
	}
	return scanContinue
}

func (s *Scanner) isCommentStart(c byte) bool {
	return (s.annotation == annotationNone || s.annotation == annotationInline) && c == '#'
}

func (s *Scanner) switchToComment() {
	if s.annotation != annotationNone && s.annotation != annotationInline {
		panic(s.newDocumentErrorAtCharacter("inside user inline comment"))
	}
	s.returnToStep.Push(s.step)
	s.step = stateAnyCommentStart
}

func stateAnyCommentStart(s *Scanner, c byte) state {
	if c != '#' {
		// any symbol inline user comment
		s.annotation = annotationNone
		s.step = stateInlineComment
		return scanContinue
	} else if s.data[s.index] == '#' { // third #
		s.annotation = annotationNone
		s.step = stateMultiLineComment
		return scanContinue
	}

	panic(s.newDocumentErrorAtCharacter("after first #"))
}

func stateInlineComment(s *Scanner, c byte) state {
	if c == '\n' || c == '\r' {
		s.step = s.returnToStep.Pop()
		s.found(lexeme.NewLine)
		s.index--
	}
	return scanContinue
}

func stateMultiLineComment(s *Scanner, c byte) state {
	if (s.index + 1) < s.dataSize {
		if c == '#' && s.data[s.index] == '#' && s.data[s.index+1] == '#' {
			s.index++ // skip second #
			s.index++ // skip third #
			s.step = s.returnToStep.Pop()
		}
	}
	return scanContinue
}

func stateKeyShortcut(s *Scanner, c byte) state {
	switch {
	case bytes.IsValidUserTypeNameByte(c):
		s.step = stateKeyShortcut
	default:
		return stateEndValue(s, c)
	}
	return scanContinue
}

const messageEIsNotAllowed = "isn't allowed 'cause not obvious it's a float or an integer"

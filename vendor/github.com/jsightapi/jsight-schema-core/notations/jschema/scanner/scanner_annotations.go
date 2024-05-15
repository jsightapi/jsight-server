package scanner

import (
	"github.com/jsightapi/jsight-schema-core/bytes"
	"github.com/jsightapi/jsight-schema-core/errs"
	"github.com/jsightapi/jsight-schema-core/kit"
	"github.com/jsightapi/jsight-schema-core/lexeme"
)

type annotation uint8

const (
	// annotationNone not inside the annotation.
	annotationNone annotation = iota

	// annotationInline inside the inline annotation.
	annotationInline

	// annotationMultiLine inside the multi-line annotation.
	annotationMultiLine
)

func (*Scanner) isAnnotationStart(c byte) bool {
	return c == '/'
}

func (s *Scanner) switchToAnnotation() {
	if !s.allowAnnotation {
		err := kit.NewJSchemaError(s.file, errs.ErrAnnotationNotAllowed.F())
		err.SetIndex(s.index - 1)
		panic(err)
	}

	s.returnToStep.Push(s.step)

	switch s.annotation {
	case annotationNone:
		s.step = stateAnyAnnotationStart

	case annotationMultiLine:
		s.step = stateInlineAnnotationStart

	default:
		panic(s.newJSchemaErrorAtCharacter("inside inline annotation"))
	}
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
	panic(s.newJSchemaErrorAtCharacter("after first slash"))
}

// ///////////////////////////
// Inline annotations states.

func stateInlineAnnotationStart(s *Scanner, c byte) state {
	// second slash - inline annotation
	if c != '/' {
		panic(s.newJSchemaErrorAtCharacter("after first slash on start inline annotation"))
	}
	s.annotation = annotationInline
	s.found(lexeme.InlineAnnotationBegin)
	s.step = stateInlineAnnotation
	return scanContinue
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
	switch {
	case bytes.IsSpace(c):

	case bytes.IsNewLine(c):
		s.found(lexeme.InlineAnnotationEnd)
		s.found(lexeme.NewLine)
		s.step = s.returnToStep.Pop()

		s.annotation = annotationNone
		if s.isInsideMultiLineAnnotation() {
			s.annotation = annotationMultiLine
		}

	case s.isCommentStart(c):
		s.switchToComment()

	case c == '-':
		s.step = stateInlineAnnotationTextPrefix2

	default:
		panic(s.newJSchemaErrorAtCharacter("after object in inline annotation"))
	}

	return scanContinue
}

func stateInlineAnnotationTextPrefix2(s *Scanner, c byte) state {
	if bytes.IsSpace(c) {
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
				panic(s.newJSchemaErrorAtCharacter("after inline annotation"))
			}
			return fn(s, c)
		}

		s.annotation = annotationNone
		if s.isInsideMultiLineAnnotation() {
			s.annotation = annotationMultiLine
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
	if !bytes.IsNewLine(c) {
		return scanContinue
	}

	s.found(lexeme.NewLine)
	fn := s.returnToStep.Pop()
	s.step = func(s *Scanner, c byte) state {
		if s.isAnnotationStart(c) {
			panic(s.newJSchemaErrorAtCharacter("after inline annotation"))
		}
		return fn(s, c)
	}

	s.annotation = annotationNone
	if s.isInsideMultiLineAnnotation() {
		s.annotation = annotationMultiLine
	}
	return scanContinue
}

// ///////////////////////////////
// Multi-line annotations states.

func stateMultiLineAnnotation(s *Scanner, c byte) state {
	if s.isNewLine(c) {
		s.found(lexeme.NewLine)
		return scanContinue
	}
	if bytes.IsBlank(c) {
		return scanContinue
	}
	if c == '{' {
		return stateFoundRootValue(s, c)
	}

	s.found(lexeme.MultiLineAnnotationTextBegin)
	s.step = stateMultiLineAnnotationText
	return s.step(s, c)
}

func stateMultiLineAnnotationTextPrefix(s *Scanner, c byte) state {
	switch {
	case bytes.IsNewLine(c):
		s.found(lexeme.NewLine)

	case bytes.IsSpace(c):

	case s.isCommentStart(c):
		s.switchToComment()

	case c == '*':
		s.step = stateMultiLineAnnotationEnd

	case c == '-':
		s.step = stateMultiLineAnnotationTextPrefix2

	default:
		panic(s.newJSchemaErrorAtCharacter("after object in multi-line annotation"))
	}
	return scanContinue
}

func stateMultiLineAnnotationTextPrefix2(s *Scanner, c byte) state {
	if bytes.IsSpace(c) {
		return scanContinue
	}
	s.found(lexeme.MultiLineAnnotationTextBegin)
	s.step = stateMultiLineAnnotationText
	return s.step(s, c)
}

func stateMultiLineAnnotationEnd(s *Scanner, c byte) state {
	if c != '/' {
		panic(s.newJSchemaErrorAtCharacter("in multi-line annotation after \"*\" character"))
	}
	// after *
	s.annotation = annotationNone
	s.found(lexeme.MultiLineAnnotationEnd)
	s.step = s.returnToStep.Pop()
	return scanContinue
}

func stateMultiLineAnnotationText(s *Scanner, c byte) state {
	if c == '*' && s.data.Byte(s.index) == '/' {
		s.found(lexeme.MultiLineAnnotationTextEnd)
		s.step = stateMultiLineAnnotationEnd
	}
	return scanContinue
}

// ///////////////////////////
// Common annotations states.

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
	if (s.boundary == 0 && (c == ':' || bytes.IsNewLine(c) || c == '\\')) || c == s.boundary || c < 0x20 {
		panic(s.newJSchemaError(errs.ErrInvalidCharacterInAnnotationObjectKey, c))
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

	case c < 0x20 || (c == '"' || bytes.IsNewLine(c)):
		panic(s.newJSchemaError(errs.ErrInvalidCharacterInAnnotationObjectKey, c))
	}
	return scanContinue
}

func stateInAnnotationObjectKeyAfter(s *Scanner, c byte) state {
	switch {
	case s.boundary == 0 && c == ':':
		return stateEndValue(s, c)

	case c == ' ':
		return scanContinue
	}
	panic(s.newJSchemaError(errs.ErrInvalidCharacterInAnnotationObjectKey, c))
}

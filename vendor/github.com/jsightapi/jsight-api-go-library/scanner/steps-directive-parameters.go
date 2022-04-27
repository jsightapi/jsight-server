package scanner

import (
	"github.com/jsightapi/jsight-api-go-library/jerr"
)

// stateParameterOrAnnotation scans the parameters of directives, and its annotation
func stateParameterOrAnnotation(s *Scanner, c byte) *jerr.JAPIError {
	switch c {
	case caseWhitespace(c):
		s.step = stateParameterOrAnnotationAfterFirstSpace
		return nil
	case CommentSign:
		return s.startComment()
	case caseNewLine(c), EOF:
		s.step = s.stepStack.Pop()
		return nil
	case AnnotationDelimiterPart:
		s.step = stateAnnotationSign2
		return nil
	default:
		return s.japiErrorUnexpectedChar("after directive keyword", "parameter")
	}
}

func stateParameterOrAnnotationAfterFirstSpace(s *Scanner, c byte) *jerr.JAPIError {
	switch c {
	case caseWhitespace(c):
		return nil
	case CommentSign:
		return s.startComment()
	case caseNewLine(c), EOF:
		s.step = s.stepStack.Pop()
		return nil
	case AnnotationDelimiterPart:
		s.step = stateAnnotationSign2
		return nil
	default:
		s.step = stateParameterStart
		return stateParameterStart(s, c)
	}
}

func stateParameterStart(s *Scanner, c byte) *jerr.JAPIError {
	s.found(ParameterBegin)

	switch c {
	case DoubleQuote:
		s.step = stateParameterInQuoted
	case caseNewLine(c):
		return s.japiErrorUnexpectedChar("in directive parameter", "")
	default:
		s.step = stateParameterWoQuoted
	}

	return nil
}

func stateParameterInQuoted(s *Scanner, c byte) *jerr.JAPIError {
	switch c {
	case caseNewLine(c), EOF:
		return s.japiErrorUnexpectedChar("in directive parameter", "closing quotation mark")
	case '\\':
		s.step = stateParameterInQuotedSlash
	case DoubleQuote:
		s.found(ParameterEnd)
		s.step = stateParameterOrAnnotation
	}
	return nil
}

func stateParameterInQuotedSlash(s *Scanner, c byte) *jerr.JAPIError {
	switch c {
	case '\\':
		s.step = stateParameterInQuoted
	case '"':
		s.step = stateParameterInQuoted
	default:
		return s.japiErrorUnexpectedChar("when escaping characters in parameters", "quotation marks or slash")
	}
	return nil
}

func stateParameterWoQuoted(s *Scanner, c byte) *jerr.JAPIError {
	switch c {
	case caseWhitespace(c), caseNewLine(c), CommentSign, EOF:
		s.foundAt(s.curIndex-1, ParameterEnd)
		s.step = stateParameterOrAnnotation
		return s.step(s, c)
	}
	return nil
}

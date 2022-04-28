package scanner

import (
	"github.com/jsightapi/jsight-api-go-library/jerr"
)

func statePa(s *Scanner, c byte) *jerr.JAPIError {
	switch c {
	case 't':
		s.step = statePat
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword Path", "t")
	}
}

func statePat(s *Scanner, c byte) *jerr.JAPIError {
	switch c {
	case 'h':
		s.found(KeywordEnd)
		s.stepStack.Push(statePathBody)
		s.step = stateParameterOrAnnotation
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword Path", "h")
	}
}

func statePathBody(s *Scanner, c byte) *jerr.JAPIError {
	switch c {
	case ContextOpenSign:
		s.found(ContextOpen)
		return nil
	case caseWhitespace(c), caseNewLine(c):
		return nil
	case CommentSign:
		return s.startComment()
	case ObjectOpen, LinkSymbol:
		return stateJSchema(s, c)
	default:
		return s.japiErrorUnexpectedChar("in Path body", "")
	}
}

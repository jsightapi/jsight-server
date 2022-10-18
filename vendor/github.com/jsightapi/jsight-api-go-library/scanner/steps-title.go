package scanner

import (
	"github.com/jsightapi/jsight-api-go-library/jerr"
)

func stateTi(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 't':
		s.step = stateTit
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword Title", "t")
	}
}

func stateTit(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'l':
		s.step = stateTitl
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword Title", "l")
	}
}

func stateTitl(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'e':
		s.found(KeywordEnd)
		s.stepStack.Push(stateExpectKeyword)
		s.step = stateParameterOrAnnotation
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword Title", "l")
	}
}

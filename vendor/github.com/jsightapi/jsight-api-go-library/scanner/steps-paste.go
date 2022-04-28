package scanner

import (
	"github.com/jsightapi/jsight-api-go-library/jerr"
)

func statePAS(s *Scanner, c byte) *jerr.JAPIError {
	switch c {
	case 'T':
		s.step = statePAST
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword MACRO", "u")
	}
}

func statePAST(s *Scanner, c byte) *jerr.JAPIError {
	switch c {
	case 'E':
		s.found(KeywordEnd)
		s.stepStack.Push(stateExpectKeyword)
		s.step = stateParameterOrAnnotation
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword MACRO", "y")
	}
}

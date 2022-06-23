package scanner

import (
	"github.com/jsightapi/jsight-api-go-library/jerr"
)

func stateMA(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'C':
		s.step = stateMAC
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword MACRO", "e")
	}
}

func stateMAC(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'R':
		s.step = stateMACR
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword MACRO", "r")
	}
}

func stateMACR(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'O':
		s.found(KeywordEnd)
		s.stepStack.Push(stateExpectKeyword)
		s.step = stateParameterOrAnnotation
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword MACRO", "y")
	}
}

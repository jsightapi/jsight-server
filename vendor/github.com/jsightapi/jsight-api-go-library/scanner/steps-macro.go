package scanner

import (
	"github.com/jsightapi/jsight-api-go-library/jerr"
)

func stateM(s *Scanner, c byte) *jerr.JAPIError {
	switch c {
	case 'A':
		s.step = stateMA
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword MACRO", "u")
	}
}

func stateMA(s *Scanner, c byte) *jerr.JAPIError {
	switch c {
	case 'C':
		s.step = stateMAC
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword MACRO", "e")
	}
}

func stateMAC(s *Scanner, c byte) *jerr.JAPIError {
	switch c {
	case 'R':
		s.step = stateMACR
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword MACRO", "r")
	}
}

func stateMACR(s *Scanner, c byte) *jerr.JAPIError {
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

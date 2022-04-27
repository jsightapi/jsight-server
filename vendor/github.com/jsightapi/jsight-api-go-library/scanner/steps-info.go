package scanner

import (
	"github.com/jsightapi/jsight-api-go-library/jerr"
)

func stateI(s *Scanner, c byte) *jerr.JAPIError {
	switch c {
	case 'N':
		s.step = stateIN
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword INFO", "N")
	}
}

func stateIN(s *Scanner, c byte) *jerr.JAPIError {
	switch c {
	case 'F':
		s.step = stateINF
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword INFO", "F")
	}
}

func stateINF(s *Scanner, c byte) *jerr.JAPIError {
	switch c {
	case 'O':
		s.found(KeywordEnd)
		s.stepStack.Push(stateExpectKeyword)
		s.step = stateParameterOrAnnotation
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword INFO", "O")
	}
}

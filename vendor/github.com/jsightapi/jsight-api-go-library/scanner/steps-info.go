package scanner

import (
	"github.com/jsightapi/jsight-api-go-library/jerr"
)

func stateI(s *Scanner, c byte) *jerr.JApiError {
	if c != 'N' {
		return stateInfoError(s, "N")
	}
	s.step = stateIN
	return nil
}

func stateIN(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'F':
		s.step = stateINF
	case 'C':
		s.step = stateINC
	default:
		return stateInfoError(s, "F")
	}
	return nil
}

func stateINF(s *Scanner, c byte) *jerr.JApiError {
	if c != 'O' {
		return stateInfoError(s, "O")
	}
	s.found(KeywordEnd)
	s.stepStack.Push(stateExpectKeyword)
	s.step = stateParameterOrAnnotation
	return nil
}

func stateInfoError(s *Scanner, expected string) *jerr.JApiError {
	return s.japiErrorUnexpectedChar("in keyword INFO", expected)
}

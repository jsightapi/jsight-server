package scanner

import (
	"github.com/jsightapi/jsight-api-go-library/jerr"
)

func stateINC(s *Scanner, c byte) *jerr.JApiError {
	if c != 'L' {
		return stateIncludeError(s, "L")
	}
	s.step = stateINCL
	return nil
}

func stateINCL(s *Scanner, c byte) *jerr.JApiError {
	if c != 'U' {
		return stateIncludeError(s, "U")
	}
	s.step = stateINCLU
	return nil
}

func stateINCLU(s *Scanner, c byte) *jerr.JApiError {
	if c != 'D' {
		return stateIncludeError(s, "D")
	}
	s.step = stateINCLUD
	return nil
}

func stateINCLUD(s *Scanner, c byte) *jerr.JApiError {
	if c != 'E' {
		return stateIncludeError(s, "E")
	}
	s.found(KeywordEnd)
	s.stepStack.Push(stateExpectKeyword)
	s.step = stateParameterOrAnnotation
	return nil
}

func stateIncludeError(s *Scanner, expected string) *jerr.JApiError {
	return s.japiErrorUnexpectedChar("in keyword INCLUDE", expected)
}

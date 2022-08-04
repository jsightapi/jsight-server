package scanner

import (
	"github.com/jsightapi/jsight-api-go-library/jerr"
)

func stateTA(s *Scanner, c byte) *jerr.JApiError {
	if c != 'G' {
		return s.japiErrorUnexpectedChar("in keyword TAG", "G")
	}
	s.found(KeywordEnd)
	s.stepStack.Push(stateExpectKeyword)
	s.step = stateParameterOrAnnotation
	return nil
}

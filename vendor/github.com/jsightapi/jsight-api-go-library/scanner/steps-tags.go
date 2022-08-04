package scanner

import (
	"github.com/jsightapi/jsight-api-go-library/jerr"
)

func stateTa(s *Scanner, c byte) *jerr.JApiError {
	if c != 'g' {
		return stateTagsError(s, "g")
	}
	s.step = stateTag
	return nil
}

func stateTag(s *Scanner, c byte) *jerr.JApiError {
	if c != 's' {
		return stateTagsError(s, "s")
	}
	s.found(KeywordEnd)
	s.stepStack.Push(stateExpectKeyword)
	s.step = stateParameterOrAnnotation
	return nil
}

func stateTagsError(s *Scanner, expected string) *jerr.JApiError {
	return s.japiErrorUnexpectedChar("in keyword \"Tags\"", expected)
}

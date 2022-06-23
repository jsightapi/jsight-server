package scanner

import (
	"fmt"

	"github.com/jsightapi/jsight-api-go-library/jerr"
)

func stateRegex(s *Scanner, c byte) *jerr.JApiError {
	if c != RegexDelimiter {
		return s.japiErrorUnexpectedChar("in the regular expression", fmt.Sprintf("%q character", RegexDelimiter))
	}

	s.found(TextBegin)
	s.step = stateRegexFirstChar
	return nil
}

func stateRegexFirstChar(s *Scanner, c byte) *jerr.JApiError {
	if c == RegexDelimiter {
		return s.japiErrorUnexpectedChar("empty regex", "")
	}

	s.step = stateRegexBody
	return s.step(s, c)
}

func stateRegexBody(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case RegexDelimiter:
		s.foundAt(s.curIndex, TextEnd)
		s.step = stateBodyEnded
	case EOF:
		return s.japiErrorUnexpectedChar("inside the regular expression", "")
	case '\\':
		s.step = stateRegexBodyAfterSlash
	}
	return nil
}

func stateRegexBodyAfterSlash(s *Scanner, _ byte) *jerr.JApiError {
	s.step = stateRegexBody
	return nil
}

package scanner //nolint:dupl

import (
	"github.com/jsightapi/jsight-api-go-library/jerr"
)

func stateBa(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 's':
		s.step = stateBas
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword BaseUrl", "s")
	}
}

func stateBas(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'e':
		s.step = stateBase
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword BaseUrl", "e")
	}
}

func stateBase(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'U':
		s.step = stateBaseU
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword BaseUrl", "U")
	}
}

func stateBaseU(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'r':
		s.step = stateBaseUr
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword BaseUrl", "r")
	}
}

func stateBaseUr(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'l':
		s.found(KeywordEnd)
		s.stepStack.Push(stateExpectKeyword)
		s.step = stateParameterOrAnnotation
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword BaseUrl", "l")
	}
}

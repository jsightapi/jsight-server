package scanner //nolint:dupl

import (
	"github.com/jsightapi/jsight-api-go-library/jerr"
)

func stateJ(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'S':
		s.step = stateJS
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword JSIGHT", "S")
	}
}

func stateJS(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'I':
		s.step = stateJSI
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword JSIGHT", "O")
	}
}

func stateJSI(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'G':
		s.step = stateJSIG
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword JSIGHT", "G")
	}
}

func stateJSIG(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'H':
		s.step = stateJSIGH
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword JSIGHT", "H")
	}
}

func stateJSIGH(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'T':
		s.found(KeywordEnd)
		s.stepStack.Push(stateExpectKeyword)
		s.step = stateParameterOrAnnotation
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword JSIGHT", "T")
	}
}

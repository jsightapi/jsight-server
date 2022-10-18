package scanner

import (
	"github.com/jsightapi/jsight-api-go-library/jerr"
)

func stateMe(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 't':
		s.step = stateMet
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword Method", "t")
	}
}

func stateMet(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'h':
		s.step = stateMeth
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword Method", "h")
	}
}

func stateMeth(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'o':
		s.step = stateMetho
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword Method", "o")
	}
}

func stateMetho(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'd':
		s.found(KeywordEnd)
		s.stepStack.Push(stateExpectKeyword)
		s.step = stateParameterOrAnnotation
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword Method", "d")
	}
}

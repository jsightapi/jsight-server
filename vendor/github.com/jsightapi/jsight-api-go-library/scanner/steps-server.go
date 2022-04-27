package scanner //nolint:dupl

import (
	"github.com/jsightapi/jsight-api-go-library/jerr"
)

func stateS(s *Scanner, c byte) *jerr.JAPIError {
	switch c {
	case 'E':
		s.step = stateSe
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword SERVER", "E")
	}
}

func stateSe(s *Scanner, c byte) *jerr.JAPIError {
	switch c {
	case 'R':
		s.step = stateSer
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword SERVER", "R")
	}
}

func stateSer(s *Scanner, c byte) *jerr.JAPIError {
	switch c {
	case 'V':
		s.step = stateServ
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword SERVER", "V")
	}
}

func stateServ(s *Scanner, c byte) *jerr.JAPIError {
	switch c {
	case 'E':
		s.step = stateServe
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword SERVER", "E")
	}
}

func stateServe(s *Scanner, c byte) *jerr.JAPIError {
	switch c {
	case 'R':
		s.found(KeywordEnd)
		s.stepStack.Push(stateExpectKeyword)
		s.step = stateParameterOrAnnotation
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword SERVER", "R")
	}
}

package scanner //nolint:dupl

import (
	"github.com/jsightapi/jsight-api-go-library/jerr"
)

func statePr(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'o':
		s.step = statePro
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword Protocol", "o")
	}
}

func statePro(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 't':
		s.step = stateProt
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword Protocol", "t")
	}
}

func stateProt(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'o':
		s.step = stateProto
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword Protocol", "o")
	}
}

func stateProto(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'c':
		s.step = stateProtoc
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword Protocol", "c")
	}
}

func stateProtoc(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'o':
		s.step = stateProtoco
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword Protocol", "o")
	}
}

func stateProtoco(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'l':
		s.found(KeywordEnd)
		s.stepStack.Push(stateExpectKeyword)
		s.step = stateParameterOrAnnotation
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword Protocol", "l")
	}
}

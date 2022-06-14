package scanner

import (
	"github.com/jsightapi/jsight-api-go-library/jerr"
)

func stateV(s *Scanner, c byte) *jerr.JAPIError {
	switch c {
	case 'e':
		s.step = stateVe
		return nil
	default:
		return s.japiErrorUnexpectedChar("in directive Version", "e")
	}
}

func stateVe(s *Scanner, c byte) *jerr.JAPIError {
	switch c {
	case 'r':
		s.step = stateVer
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword Version", "r")
	}
}

func stateVer(s *Scanner, c byte) *jerr.JAPIError {
	switch c {
	case 's':
		s.step = stateVers
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword Version", "s")
	}
}

func stateVers(s *Scanner, c byte) *jerr.JAPIError {
	switch c {
	case 'i':
		s.step = stateVersi
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword Version", "i")
	}
}

func stateVersi(s *Scanner, c byte) *jerr.JAPIError {
	switch c {
	case 'o':
		s.step = stateVersio
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword Version", "o")
	}
}

func stateVersio(s *Scanner, c byte) *jerr.JAPIError {
	switch c {
	case 'n':
		s.found(KeywordEnd)
		s.stepStack.Push(stateExpectKeyword)
		s.step = stateParameterOrAnnotation
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword Version", "n")
	}
}

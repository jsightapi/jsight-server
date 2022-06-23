package scanner

import (
	"github.com/jsightapi/jsight-api-go-library/jerr"
)

func stateG(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'E':
		s.step = stateGE
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword GET", "E")
	}
}

func stateGE(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'T':
		s.found(KeywordEnd)
		s.stepStack.Push(stateExpectKeyword)
		s.step = stateParameterOrAnnotation
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword GET", "T")
	}
}

func statePO(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'S':
		s.step = statePOS
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword POST", "S")
	}
}

func statePOS(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'T':
		s.found(KeywordEnd)
		s.stepStack.Push(stateExpectKeyword)
		s.step = stateParameterOrAnnotation
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword POST", "T")
	}
}

func statePU(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'T':
		s.found(KeywordEnd)
		s.stepStack.Push(stateExpectKeyword)
		s.step = stateParameterOrAnnotation
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword PUT", "T")
	}
}

func statePAT(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'C':
		s.step = statePATC
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword PATCH", "C")
	}
}

func statePATC(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'H':
		s.found(KeywordEnd)
		s.stepStack.Push(stateExpectKeyword)
		s.step = stateParameterOrAnnotation
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword PATCH", "H")
	}
}

func stateDE(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'L':
		s.step = stateDEL
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword DELETE", "L")
	}
}

func stateDEL(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'E':
		s.step = stateDELE
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword DELETE", "E")
	}
}

func stateDELE(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'T':
		s.step = stateDELET
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword DELETE", "T")
	}
}

func stateDELET(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'E':
		s.found(KeywordEnd)
		s.stepStack.Push(stateExpectKeyword)
		s.step = stateParameterOrAnnotation
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword DELETE", "E")
	}
}

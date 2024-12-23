package scanner

import (
	"github.com/jsightapi/jsight-api-core/jerr"
)

func stateO(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'p':
		s.step = stateOp
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword OperationId", "p")
	}
}

func stateOp(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'e':
		s.step = stateOpe
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword OperationId", "e")
	}
}

func stateOpe(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'r':
		s.step = stateOper
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword OperationId", "r")
	}
}

func stateOper(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'a':
		s.step = stateOpera
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword OperationId", "a")
	}
}

func stateOpera(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 't':
		s.step = stateOperat
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword OperationId", "t")
	}
}

func stateOperat(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'i':
		s.step = stateOperati
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword OperationId", "i")
	}
}

func stateOperati(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'o':
		s.step = stateOperatio
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword OperationId", "o")
	}
}

func stateOperatio(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'n':
		s.step = stateOperation
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword OperationId", "n")
	}
}

func stateOperation(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'I':
		s.step = stateOperationI
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword OperationId", "I")
	}
}

func stateOperationI(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'd':
		s.found(KeywordEnd)
		s.stepStack.Push(stateExpectKeyword)
		s.step = stateParameterOrAnnotation
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword OperationId", "d")
	}
}

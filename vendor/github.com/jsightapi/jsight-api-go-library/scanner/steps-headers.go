package scanner

import (
	"github.com/jsightapi/jsight-api-go-library/jerr"
)

func stateH(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'e':
		s.step = stateHe
		return nil
	default:
		return s.japiErrorUnexpectedChar("in directive Headers", "e")
	}
}

func stateHe(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'a':
		s.step = stateHea
		return nil
	default:
		return s.japiErrorUnexpectedChar("in directive Headers", "a")
	}
}

func stateHea(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'd':
		s.step = stateHead
		return nil
	default:
		return s.japiErrorUnexpectedChar("in directive Headers", "d")
	}
}

func stateHead(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'e':
		s.step = stateHeade
		return nil
	default:
		return s.japiErrorUnexpectedChar("in directive Headers", "e")
	}
}

func stateHeade(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'r':
		s.step = stateHeader
		return nil
	default:
		return s.japiErrorUnexpectedChar("in directive Headers", "r")
	}
}

func stateHeader(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 's':
		s.found(KeywordEnd)
		s.stepStack.Push(stateHeaderBody)
		s.step = stateParameterOrAnnotation
		return nil
	default:
		return s.japiErrorUnexpectedChar("in directive Headers", "s")
	}
}

func stateHeaderBody(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case ContextOpenSign:
		s.found(ContextOpen)
		return nil
	case caseWhitespace(c), caseNewLine(c):
		return nil
	case CommentSign:
		return s.startComment()
	case ObjectOpen, LinkSymbol:
		return stateJSchema(s, c)
	default:
		return s.japiErrorUnexpectedChar("in Headers body", "")
	}
}

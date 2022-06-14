package scanner

import (
	"github.com/jsightapi/jsight-api-go-library/jerr"
)

func stateQ(s *Scanner, c byte) *jerr.JAPIError {
	switch c {
	case 'u':
		s.step = stateQu
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword Query", "u")
	}
}

func stateQu(s *Scanner, c byte) *jerr.JAPIError {
	switch c {
	case 'e':
		s.step = stateQue
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword Query", "e")
	}
}

func stateQue(s *Scanner, c byte) *jerr.JAPIError {
	switch c {
	case 'r':
		s.step = stateQuer
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword Query", "r")
	}
}

func stateQuer(s *Scanner, c byte) *jerr.JAPIError {
	switch c {
	case 'y':
		s.found(KeywordEnd)
		s.stepStack.Push(stateQueryBodyOrKeyword)
		s.step = stateParameterOrAnnotation
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword Query", "y")
	}
}

func stateQueryBodyOrKeyword(s *Scanner, c byte) *jerr.JAPIError {
	switch c {
	case ContextOpenSign:
		s.found(ContextOpen)
		return nil
	case caseWhitespace(c), caseNewLine(c):
		return nil
	case CommentSign:
		return s.startComment()
	default:
		return stateJSchema(s, c)
	}
}

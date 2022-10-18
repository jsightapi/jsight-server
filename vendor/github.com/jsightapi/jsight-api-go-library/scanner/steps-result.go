package scanner //nolint:dupl

import (
	"github.com/jsightapi/jsight-api-go-library/jerr"
)

func stateRes(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'u':
		s.step = stateResu
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword Result", "u")
	}
}

func stateResu(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'l':
		s.step = stateResul
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword Result", "l")
	}
}

func stateResul(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 't':
		s.found(KeywordEnd)
		s.stepStack.Push(stateResultBody)
		s.step = stateParameterOrAnnotation
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword Result", "t")
	}
}

func stateResultBody(s *Scanner, c byte) *jerr.JApiError {
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

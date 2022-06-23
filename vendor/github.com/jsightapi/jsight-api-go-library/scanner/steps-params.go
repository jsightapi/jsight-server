package scanner //nolint:dupl

import (
	"github.com/jsightapi/jsight-api-go-library/jerr"
)

func statePar(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'a':
		s.step = statePara
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword Params", "a")
	}
}

func statePara(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'm':
		s.step = stateParam
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword Params", "m")
	}
}

func stateParam(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 's':
		s.found(KeywordEnd)
		s.stepStack.Push(stateParamsBody)
		s.step = stateParameterOrAnnotation
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword Params", "s")
	}
}

func stateParamsBody(s *Scanner, c byte) *jerr.JApiError {
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

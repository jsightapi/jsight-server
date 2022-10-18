package scanner

import (
	"github.com/jsightapi/jsight-api-go-library/jerr"
)

func stateTy(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'P':
		s.step = stateTyp
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword TYPE", "P")
	}
}

func stateTyp(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'E':
		s.found(KeywordEnd)
		s.stepStack.Push(stateTypeBodyOrKeyword)
		s.step = stateParameterOrAnnotation
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword TYPE", "E")
	}
}

func stateTypeBodyOrKeyword(s *Scanner, c byte) *jerr.JApiError {
	if s.isDirectiveParameterHasAnyOrEmpty() {
		if s.isDirectiveParameterHasRegexNotation() {
			s.stepStack.Push(stateRegex)
		} else {
			s.stepStack.Push(stateJSchema)
		}
		s.step = stateTypeBody
	} else {
		s.step = stateExpectKeyword
	}
	return s.step(s, c)
}

func stateTypeBody(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case caseWhitespace(c), caseNewLine(c):
		return nil
	case ContextOpenSign:
		s.found(ContextOpen)
		return nil
	default:
		s.step = s.stepStack.Pop()
		return s.step(s, c)
	}
}

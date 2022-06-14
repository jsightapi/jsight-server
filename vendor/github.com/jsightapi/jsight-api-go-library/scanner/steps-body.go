package scanner

import (
	"github.com/jsightapi/jsight-api-go-library/jerr"
)

func stateBo(s *Scanner, c byte) *jerr.JAPIError {
	switch c {
	case 'd':
		s.step = stateBod
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword Body", "d")
	}
}

func stateBod(s *Scanner, c byte) *jerr.JAPIError {
	switch c {
	case 'y':
		s.found(KeywordEnd)
		s.stepStack.Push(stateBodyBodyOrKeyword)
		s.step = stateParameterOrAnnotation
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword Body", "y")
	}
}

func stateBodyBodyOrKeyword(s *Scanner, c byte) *jerr.JAPIError {
	if !s.isDirectiveParameterHasTypeOrAnyOrEmpty() {
		if s.isDirectiveParameterHasRegexNotation() {
			s.stepStack.Push(stateRegex)
		} else {
			s.stepStack.Push(stateJSchema)
		}
		s.step = stateBodyBody
	} else {
		s.step = stateExpectKeyword
	}
	return s.step(s, c)
}

func stateBodyBody(s *Scanner, c byte) *jerr.JAPIError {
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

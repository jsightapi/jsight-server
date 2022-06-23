package scanner

import (
	"github.com/jsightapi/jsight-api-go-library/jerr"
)

func stateReq(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'u':
		s.step = stateRequ
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword Request", "u")
	}
}

func stateRequ(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'e':
		s.step = stateReque
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword Request", "e")
	}
}

func stateReque(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 's':
		s.step = stateReques
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword Request", "s")
	}
}

func stateReques(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 't':
		s.found(KeywordEnd)
		s.stepStack.Push(stateRequestBodyOrKeyword)
		s.step = stateParameterOrAnnotation
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword Request", "t")
	}
}

func stateRequestBodyOrKeyword(s *Scanner, c byte) *jerr.JApiError {
	if !s.isDirectiveParameterHasTypeOrAnyOrEmpty() {
		if s.isDirectiveParameterHasRegexNotation() {
			s.stepStack.Push(stateRegex)
		} else {
			s.stepStack.Push(stateJSchema)
		}
		s.step = stateRequestBody
	} else {
		s.step = stateExpectKeyword
	}
	return s.step(s, c)
}

func stateRequestBody(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case ContextOpenSign:
		s.found(ContextOpen)
		return nil
	case caseWhitespace(c), caseNewLine(c):
		return nil
	case CommentSign:
		return s.startComment()
	case 'B', 'H', 'P', 'I': // directives: Body, Header, PASTE, INCLUDE
		return stateExpectKeyword(s, c)
	default:
		s.step = s.stepStack.Pop()
		return s.step(s, c)
	}
}

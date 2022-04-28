package scanner

import (
	"github.com/jsightapi/jsight-api-go-library/jerr"
)

func stateResponseKeywordStarted(s *Scanner, c byte) *jerr.JAPIError {
	switch c {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		s.step = stateResponseKeywordSecond
		return nil
	default:
		return s.japiErrorUnexpectedChar("at response directive", "digit")
	}
}

func stateResponseKeywordSecond(s *Scanner, c byte) *jerr.JAPIError {
	switch c {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		s.found(KeywordEnd)
		s.stepStack.Push(stateResponseBodyOrKeyword)
		s.step = stateParameterOrAnnotation
		return nil
	default:
		return s.japiErrorUnexpectedChar("at response directive", "digit")
	}
}

func stateResponseBodyOrKeyword(s *Scanner, c byte) *jerr.JAPIError {
	if !s.isDirectiveParameterHasTypeOrAnyOrEmpty() {
		if s.isDirectiveParameterHasRegexNotation() {
			s.stepStack.Push(stateRegex)
		} else {
			s.stepStack.Push(stateJSchema)
		}
		s.step = stateResponseBody
	} else {
		s.step = stateExpectKeyword
	}
	return s.step(s, c)
}

func stateResponseBody(s *Scanner, c byte) *jerr.JAPIError {
	switch c {
	case ContextOpenSign:
		s.found(ContextOpen)
		return nil
	case caseWhitespace(c), caseNewLine(c):
		return nil
	case CommentSign:
		return s.startComment()
	case 'B', 'H', 'P': // directives: Body, Header, PASTE
		return stateExpectKeyword(s, c)
	default:
		s.step = s.stepStack.Pop()
		return s.step(s, c)
	}
}

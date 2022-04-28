package scanner

import (
	"github.com/jsightapi/jsight-api-go-library/directive"
	"github.com/jsightapi/jsight-api-go-library/jerr"
)

func stateDe(s *Scanner, c byte) *jerr.JAPIError {
	switch c {
	case 's':
		s.step = stateDes
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword Description", "s")
	}
}

func stateDes(s *Scanner, c byte) *jerr.JAPIError {
	switch c {
	case 'c':
		s.step = stateDesc
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword Description", "c")
	}
}

func stateDesc(s *Scanner, c byte) *jerr.JAPIError {
	switch c {
	case 'r':
		s.step = stateDescr
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword Description", "r")
	}
}

func stateDescr(s *Scanner, c byte) *jerr.JAPIError {
	switch c {
	case 'i':
		s.step = stateDescri
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword Description", "i")
	}
}

func stateDescri(s *Scanner, c byte) *jerr.JAPIError {
	switch c {
	case 'p':
		s.step = stateDescrip
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword Description", "p")
	}
}

func stateDescrip(s *Scanner, c byte) *jerr.JAPIError {
	switch c {
	case 't':
		s.step = stateDescript
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword Description", "t")
	}
}

func stateDescript(s *Scanner, c byte) *jerr.JAPIError {
	switch c {
	case 'i':
		s.step = stateDescripti
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword Description", "i")
	}
}

func stateDescripti(s *Scanner, c byte) *jerr.JAPIError {
	switch c {
	case 'o':
		s.step = stateDescriptio
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword Description", "o")
	}
}

func stateDescriptio(s *Scanner, c byte) *jerr.JAPIError {
	switch c {
	case 'n':
		s.found(KeywordEnd)
		s.stepStack.Push(stateDescriptionTextBeginStarter)
		s.step = stateParameterOrAnnotation
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword Description", "n")
	}
}

func stateDescriptionTextBeginStarter(s *Scanner, c byte) *jerr.JAPIError {
	s.found(TextBegin)
	s.step = stateDescriptionTextBegin
	return stateDescriptionTextBegin(s, c)
}

func stateDescriptionTextBegin(s *Scanner, c byte) *jerr.JAPIError {
	switch c {
	case caseNewLine(c), caseWhitespace(c):
		return nil
	case EOF:
		s.foundAt(s.curIndex-1, TextEnd)
		return nil
	case '(':
		s.step = stateDescriptionTextBracketsInner
		return nil
	default:
		s.step = stateDescriptionTextNewline
		return stateDescriptionTextNewline(s, c)
	}
}

func stateDescriptionTextBracketsInner(s *Scanner, c byte) *jerr.JAPIError {
	if isNewLine(c) {
		s.step = stateDescriptionTextBracketsInnerNewLine
	}
	return nil
}

func stateDescriptionTextBracketsInnerNewLine(s *Scanner, c byte) *jerr.JAPIError {
	switch c {
	case caseWhitespace(c), caseNewLine(c):
		return nil
	case ContextCloseSign:
		s.found(TextEnd)
		s.step = stateExpectKeyword
		return nil
	default:
		s.step = stateDescriptionTextBracketsInner
		return nil
	}
}

func stateDescriptionText(s *Scanner, c byte) *jerr.JAPIError {
	switch c {
	case caseNewLine(c):
		s.step = stateDescriptionTextNewline
		return nil
	case EOF:
		s.foundAt(s.curIndex-1, TextEnd)
		return nil
	default:
		return nil
	}
}

func stateDescriptionTextNewline(s *Scanner, c byte) *jerr.JAPIError {
	switch c {
	case caseWhitespace(c), caseNewLine(c):
		return nil
	case EOF:
		s.foundAt(s.curIndex-1, TextEnd)
		return nil
	default:
		switch {
		case s.isDirective():
			s.foundAt(s.curIndex-1, TextEnd)
			s.step = stateExpectKeyword
			s.curIndex--
			return nil
		case c == ContextCloseSign:
			s.foundAt(s.curIndex-1, TextEnd)
			s.found(ContextClose)
			s.step = stateExpectKeyword
			return nil
		default:
			s.step = stateDescriptionText
			return nil
		}
	}
}

func (s *Scanner) isDirective() bool {
	b, err := s.data.LineFrom(s.curIndex)
	if err != nil {
		return false
	}

	return directive.IsStartWithDirective(b)
}

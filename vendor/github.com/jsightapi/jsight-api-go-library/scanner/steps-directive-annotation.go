package scanner

import (
	"github.com/jsightapi/jsight-api-go-library/jerr"
)

func stateAnnotationSign2(s *Scanner, c byte) *jerr.JAPIError {
	switch c {
	case AnnotationDelimiterPart:
		s.step = stateAnnotationTextStart
	case '*':
		s.step = stateMultilineAnnotationTextStart
	default:
		s.curIndex -= 2
		s.step = stateParameterStart
	}
	return nil
}

func stateAnnotationTextStart(s *Scanner, c byte) *jerr.JAPIError {
	s.found(AnnotationBegin)
	s.step = stateAnnotation
	return stateAnnotation(s, c)
}

func stateAnnotation(s *Scanner, c byte) *jerr.JAPIError {
	switch c {
	case CommentSign:
		s.foundAt(s.curIndex-1, AnnotationEnd)
		s.step = stateSingleComment
		return nil
	case caseNewLine(c), EOF:
		s.foundAt(s.curIndex-1, AnnotationEnd)
		s.step = s.stepStack.Pop()
		return s.step(s, c)
	}
	return nil
}

func stateMultilineAnnotationTextStart(s *Scanner, c byte) *jerr.JAPIError {
	s.foundAt(s.curIndex, AnnotationBegin)
	s.step = stateMultilineAnnotation
	return stateMultilineAnnotation(s, c)
}

func stateMultilineAnnotation(s *Scanner, c byte) *jerr.JAPIError {
	if c == AnnotationDelimiterPart && s.data[s.curIndex-1] == '*' {
		s.foundAt(s.curIndex-2, AnnotationEnd)
		s.step = s.stepStack.Pop()
	} else if c == EOF {
		return s.japiErrorUnexpectedChar("multiline annotation", "*/")
	}
	return nil
}

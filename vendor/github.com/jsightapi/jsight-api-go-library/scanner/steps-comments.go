package scanner

import (
	"github.com/jsightapi/jsight-api-go-library/jerr"
)

// scanner goes into comment mode
func (s *Scanner) startComment() *jerr.JApiError {
	s.stepStack.Push(s.step)
	s.step = stateCommentStarted
	return nil
}

// byte that "ended" comment (eof, break) should be scanned as part of previous step
func (s *Scanner) endCommentLine(c byte) *jerr.JApiError {
	s.step = s.stepStack.Pop()
	return s.step(s, c)
}

func stateCommentStarted(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case CommentSign:
		s.step = stateCommentDouble
		return nil
	default:
		return stateSingleComment(s, c)
	}
}

func stateCommentDouble(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case CommentSign:
		s.step = stateCommentBlock
		return nil
	default:
		return stateSingleComment(s, c)
	}
}

func stateSingleComment(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case caseNewLine(c), EOF:
		return s.endCommentLine(c)
	default:
		return nil // anything allowed
	}
}

func stateCommentBlock(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case EOF:
		return s.japiErrorUnexpectedChar("not found boundary end symbols", "###")
	case CommentSign:
		s.step = stateCommentOnceClosed
		return nil
	default:
		return nil // anything allowed
	}
}

func stateCommentOnceClosed(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case CommentSign:
		s.step = stateCommentTwiceClosed
		return nil
	default:
		s.step = stateCommentBlock
		return s.step(s, c)
	}
}

func stateCommentTwiceClosed(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case CommentSign:
		// end comment block
		s.step = s.stepStack.Pop()
		return nil
	default:
		s.step = stateCommentBlock
		return s.step(s, c)
	}
}

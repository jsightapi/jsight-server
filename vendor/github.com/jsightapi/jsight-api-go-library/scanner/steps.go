package scanner

import (
	"github.com/jsightapi/jsight-api-go-library/jerr"
)

func stateRoot(s *Scanner, c byte) *jerr.JApiError {
	if c == CommentSign {
		return s.startComment()
	}
	return stateExpectKeyword(s, c)
}

func stateExpectKeyword(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case caseNewLine(c), caseWhitespace(c), EOF:
		return nil
	case CommentSign:
		return s.startComment()
	case ContextOpenSign:
		s.found(ContextOpen)
		s.step = stateContextOpenedOnNewline
		return nil
	case ContextCloseSign:
		s.found(ContextClose)
		s.step = stateContextClosed
		return nil
	case 'U': // URL
		s.found(KeywordBegin)
		s.step = stateU
		return nil
	case 'R': // Request, Result
		s.found(KeywordBegin)
		s.step = stateR
		return nil
	case 'G': // GET
		s.found(KeywordBegin)
		s.step = stateG
		return nil
	case 'P': // POST, PUT, PATCH, Path, PASTE, Protocol, Params
		s.found(KeywordBegin)
		s.step = stateP
		return nil
	case 'D': // DELETE, Description
		s.found(KeywordBegin)
		s.step = stateD
		return nil
	case 'S': // SERVER
		s.found(KeywordBegin)
		s.step = stateS
		return nil
	case 'I': // INFO, INCLUDE
		s.found(KeywordBegin)
		s.step = stateI
		return nil
	case 'T': // Title, TYPE
		s.found(KeywordBegin)
		s.step = stateT
		return nil
	case 'V': // Version
		s.found(KeywordBegin)
		s.step = stateV
		return nil
	case 'J': // JSIGHT
		s.found(KeywordBegin)
		s.step = stateJ
		return nil
	case 'E': // ENUM
		s.found(KeywordBegin)
		s.step = stateE
		return nil
	case 'Q': // Query
		s.found(KeywordBegin)
		s.step = stateQ
		return nil
	case 'M': // MACRO, Method
		s.found(KeywordBegin)
		s.step = stateM
		return nil
	case 'B': // Body, BaseUrl
		s.found(KeywordBegin)
		s.step = stateB
		return nil
	case 'H': // Headers
		s.found(KeywordBegin)
		s.step = stateH
		return nil
	case '1', '2', '3', '4', '5': // HTTP responses
		s.found(KeywordBegin)
		s.step = stateResponseKeywordStarted
		return nil
	}
	return s.japiErrorUnexpectedChar("at directive beginning", "")
}

func stateB(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'o': // Body
		s.step = stateBo
		return nil
	case 'a': // BaseUrl
		s.step = stateBa
		return nil
	default:
		return s.japiErrorUnexpectedChar("in directive name", "")
	}
}

func stateD(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'E': // DELETE
		s.step = stateDE
		return nil
	case 'e': // Description
		s.step = stateDe
		return nil
	default:
		return s.japiErrorUnexpectedChar("in directive name", "")
	}
}

func stateM(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'A': // MACRO
		s.step = stateMA
		return nil
	case 'e': // Method
		s.step = stateMe
		return nil
	default:
		return s.japiErrorUnexpectedChar("in directive name", "")
	}
}

func stateP(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'O': // POST
		s.step = statePO
		return nil
	case 'U': // PUT
		s.step = statePU
		return nil
	case 'A': // PATCH, PASTE
		s.step = statePA
		return nil
	case 'a': // Path, Params
		s.step = statePa
		return nil
	case 'r': // Protocol
		s.step = statePr
		return nil
	default:
		return s.japiErrorUnexpectedChar("in directive name", "")
	}
}

func statePA(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'S': // PASTE
		s.step = statePAS
		return nil
	case 'T': // PATCH
		s.step = statePAT
		return nil
	default:
		return s.japiErrorUnexpectedChar("in directive name", "")
	}
}

func statePa(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 't': // Path
		s.step = statePat
		return nil
	case 'r': // Params
		s.step = statePar
		return nil
	default:
		return s.japiErrorUnexpectedChar("in directive name", "")
	}
}

func stateR(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'e': // Request, Result
		s.step = stateRe
		return nil
	default:
		return s.japiErrorUnexpectedChar("in directive name", "")
	}
}

func stateRe(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'q': // Request
		s.step = stateReq
		return nil
	case 's': // Result
		s.step = stateRes
		return nil
	default:
		return s.japiErrorUnexpectedChar("in directive name", "")
	}
}

func stateT(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'i': // Title
		s.step = stateTi
		return nil
	case 'Y': // TYPE
		s.step = stateTy
		return nil
	default:
		return s.japiErrorUnexpectedChar("in directive name", "")
	}
}

// this is only good until schema takes up whole body.
// if we later need to continue body after schema, we can make processSchema() take 'nextStep' func as a parameter.
func stateSchemaClosed(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case caseWhitespace(c):
		s.foundAt(s.curIndex-1, SchemaEnd)
		s.step = stateBodyEnded
		return nil
	case caseNewLine(c), EOF:
		s.foundAt(s.curIndex-1, SchemaEnd)
		s.step = stateExpectKeyword
		return nil
	default:
		return s.japiErrorUnexpectedChar("after schema", "")
	}
}

// any directive's body, not the "Body" directive
// this state allows comments, because body was properly ended at least with whitespace
func stateBodyEnded(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case caseWhitespace(c):
		return nil
	case caseNewLine(c), EOF:
		s.step = stateExpectKeyword
		return nil
	case CommentSign:
		return s.startComment()
	default:
		return s.japiErrorUnexpectedChar("after body", "")
	}
}

func stateContextClosed(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case caseWhitespace(c), EOF:
		return nil
	case caseNewLine(c):
		s.step = stateExpectKeyword
		return nil
	case CommentSign:
		return s.startComment()
	default:
		return s.japiErrorUnexpectedChar("after explicit context close", "")
	}
}

func stateContextOpenedOnNewline(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case caseWhitespace(c):
		return nil
	case caseNewLine(c):
		s.step = stateExpectKeyword
		return nil
	case CommentSign:
		return s.startComment()
	default:
		return s.japiErrorBasic(jerr.ApartFromTheOpeningParenthesis)
	}
}

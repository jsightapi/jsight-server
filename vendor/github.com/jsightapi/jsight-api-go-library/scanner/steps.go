package scanner

import (
	"github.com/jsightapi/jsight-api-go-library/jerr"
)

func stateRoot(s *Scanner, c byte) *jerr.JAPIError {
	if c == CommentSign {
		return s.startComment()
	}
	return stateExpectKeyword(s, c)
}

func stateExpectKeyword(s *Scanner, c byte) *jerr.JAPIError {
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
	case 'R': // URL
		s.found(KeywordBegin)
		s.step = stateR
		return nil
	case 'G': // GET
		s.found(KeywordBegin)
		s.step = stateG
		return nil
	case 'P': // POST, PUT, PATCH, Path, PASTE
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
	case 'I': // INFO
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
	case 'M': // MACRO
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

func stateB(s *Scanner, c byte) *jerr.JAPIError {
	switch c {
	case 'o': // Body
		s.step = stateBo
		return nil
	case 'a': // BaseUrl
		s.step = stateBa
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword Body/BaseUrl", "o")
	}
}

func stateD(s *Scanner, c byte) *jerr.JAPIError {
	switch c {
	case 'E': // DELETE
		s.step = stateDE
		return nil
	case 'e': // Description
		s.step = stateDe
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword DELETE", "E")
	}
}

func stateP(s *Scanner, c byte) *jerr.JAPIError {
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
	case 'a': // Path
		s.step = statePa
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword ", "'O', 'U', 'A', 'a')")
	}
}

func statePA(s *Scanner, c byte) *jerr.JAPIError {
	switch c {
	case 'S': // PASTE
		s.step = statePAS
		return nil
	case 'T': // PATCH
		s.step = statePAT
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword PATCH", "T")
	}
}

func stateT(s *Scanner, c byte) *jerr.JAPIError {
	switch c {
	case 'i': // Title
		s.step = stateTi
		return nil
	case 'Y': // TYPE
		s.step = stateTy
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword Title", "'Y' or 'i'")
	}
}

// this is only good until schema takes up whole body.
// if we later need to continue body after schema, we can make processSchema() take 'nextStep' func as a parameter.
func stateSchemaClosed(s *Scanner, c byte) *jerr.JAPIError {
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

func stateJsonArrayClosed(s *Scanner, c byte) *jerr.JAPIError {
	switch c {
	case caseWhitespace(c):
		s.found(JsonArrayEnd)
		s.step = stateBodyEnded
		return nil
	case caseNewLine(c), EOF:
		s.found(JsonArrayEnd)
		s.step = stateExpectKeyword
		return nil
	default:
		return s.japiErrorUnexpectedChar("after json array", "")
	}
}

// any directive's body, not the "Body" directive
// this state allows comments, because body was properly ended at least with whitespace
func stateBodyEnded(s *Scanner, c byte) *jerr.JAPIError {
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

func stateContextClosed(s *Scanner, c byte) *jerr.JAPIError {
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

func stateContextOpenedOnNewline(s *Scanner, c byte) *jerr.JAPIError {
	switch c {
	case caseWhitespace(c):
		return nil
	case caseNewLine(c):
		s.step = stateExpectKeyword
		return nil
	case CommentSign:
		return s.startComment()
	default:
		return s.japiErrorUnexpectedChar("after explicit context open", "")
	}
}

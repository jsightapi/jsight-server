package scanner

import (
	"github.com/jsightapi/jsight-schema-go-library/bytes"
	"github.com/jsightapi/jsight-schema-go-library/fs"
	"github.com/jsightapi/jsight-schema-go-library/kit"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema/rules"

	"github.com/jsightapi/jsight-api-go-library/jerr"
)

func stateE(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'N':
		s.step = stateEN
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword ENUM", "N")
	}
}

func stateEN(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'U':
		s.step = stateENU
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword ENUM", "U")
	}
}

func stateENU(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case 'M':
		s.found(KeywordEnd)
		s.stepStack.Push(stateEnumBody)
		s.step = stateParameterOrAnnotation
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword ENUM", "M")
	}
}

func stateEnumBody(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case ContextOpenSign:
		s.found(ContextOpen)
		return nil
	case caseWhitespace(c), caseNewLine(c):
		return nil
	case CommentSign:
		return s.startComment()
	case ArrayOpen:
		return s.scanEnumBody(c)
	default:
		return s.japiErrorUnexpectedChar("after Enum directive", "")
	}
}

func (s *Scanner) scanEnumBody(_ byte) *jerr.JApiError {
	s.found(EnumBegin)
	enumLength, je := s.readEnumWithJsc()
	if je != nil {
		return je
	}
	if enumLength > 0 {
		s.curIndex += bytes.Index(enumLength - 1)
	}
	s.step = stateEnumBodyClose
	return nil
}

func stateEnumBodyClose(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case caseWhitespace(c):
		s.foundAt(s.curIndex-1, EnumEnd)
		s.step = stateEnumBodyEnded
		return nil
	case caseNewLine(c), EOF:
		s.foundAt(s.curIndex-1, EnumEnd)
		s.step = stateExpectKeyword
		return nil
	default:
		return s.japiErrorUnexpectedChar("after enum", "")
	}
}

// stateEnumBodyEnded any directive's body, not the "Body" directive
// this state allows comments, because body was properly ended at least with whitespace
func stateEnumBodyEnded(s *Scanner, c byte) *jerr.JApiError {
	switch c {
	case caseWhitespace(c):
		return nil
	case caseNewLine(c), EOF:
		s.step = stateExpectKeyword
		return nil
	case CommentSign:
		return s.startComment()
	default:
		return s.japiErrorUnexpectedChar("after enum body", "")
	}
}

func (s *Scanner) readEnumWithJsc() (uint, *jerr.JApiError) {
	fc := s.file.Content()
	file := fs.NewFile("", fc.Slice(s.curIndex, bytes.Index(fc.Len()-1)))

	l, err := rules.EnumFromFile(file).Len()
	if err != nil {
		err := kit.ConvertError(file, err)
		return 0, s.japiError(err.Message(), s.curIndex+bytes.Index(err.Position()))
	}
	return l, nil
}

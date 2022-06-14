package scanner

import (
	"encoding/json"

	"github.com/jsightapi/jsight-schema-go-library/bytes"
	"github.com/jsightapi/jsight-schema-go-library/fs"
	"github.com/jsightapi/jsight-schema-go-library/kit"

	"github.com/jsightapi/jsight-api-go-library/jerr"
)

func stateE(s *Scanner, c byte) *jerr.JAPIError {
	switch c {
	case 'N':
		s.step = stateEN
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword ENUM", "N")
	}
}

func stateEN(s *Scanner, c byte) *jerr.JAPIError {
	switch c {
	case 'U':
		s.step = stateENU
		return nil
	default:
		return s.japiErrorUnexpectedChar("in keyword ENUM", "U")
	}
}

func stateENU(s *Scanner, c byte) *jerr.JAPIError {
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

func stateEnumBody(s *Scanner, c byte) *jerr.JAPIError {
	switch c {
	case caseWhitespace(c), caseNewLine(c):
		return nil
	case CommentSign:
		return s.startComment()
	case ArrayOpen:
		return s.scanJsonArray(c)
	default:
		return s.japiErrorUnexpectedChar("after Enum directive", "")
	}
}

// pass rest of the file to jsc scanner to find out where array ends
func (s *Scanner) scanJsonArray(_ byte) *jerr.JAPIError {
	s.found(JsonArrayBegin)
	arrLength, je := s.readArrayWithJsc()
	if je != nil {
		return je
	}
	s.curIndex += bytes.Index(arrLength - 1)
	s.step = stateJsonArrayClosed
	return nil
}

func (s *Scanner) readArrayWithJsc() (uint, *jerr.JAPIError) {
	b := s.file.Content()
	bb := b.Slice(s.curIndex, bytes.Index(b.Len()-1))
	f := fs.NewFile("", bb)
	jsonLength, err := kit.LengthOfJson(f)
	if err != nil {
		return 0, s.japiError(err.Message(), s.curIndex+bytes.Index(err.Position()))
	}

	// validate that it is indeed array, not json (though that will never happen in current logic)
	arr := bb[:jsonLength]
	model := make([]interface{}, 0)
	marshalErr := json.Unmarshal(arr, &model)
	if marshalErr != nil {
		return 0, s.japiError("ENUM body is not a json array", s.curIndex)
	}

	return jsonLength, nil
}

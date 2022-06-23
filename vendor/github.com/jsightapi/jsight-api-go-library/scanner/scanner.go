package scanner

import (
	"github.com/jsightapi/jsight-schema-go-library/bytes"
	"github.com/jsightapi/jsight-schema-go-library/fs"

	"github.com/jsightapi/jsight-api-go-library/jerr"
)

type stepFunc func(*Scanner, byte) *jerr.JApiError

type Scanner struct {
	data bytes.Bytes
	file *fs.File

	// step a function, that will be evaluated for next byte.
	step stepFunc

	// stepStack keeping previous step when going into comments to return to when
	// comment is done.
	stepStack stepFuncStack

	// finds gather beginnings or ends of lexemes during steps.
	finds []LexemeEvent

	// stack to keep the beginning of a lexeme until ending is found.
	stack                   eventStack
	lastDirectiveParameters []*Lexeme
	curIndex                bytes.Index
	dataSize                bytes.Index
}

func NewJApiScanner(file *fs.File) *Scanner {
	s := Scanner{
		step:                    stateRoot,
		file:                    file,
		data:                    file.Content(),
		finds:                   make([]LexemeEvent, 0, 5),
		stepStack:               make(stepFuncStack, 0, 5),
		lastDirectiveParameters: make([]*Lexeme, 0, 5),
		stack:                   make(eventStack, 0, 5),
	}
	s.dataSize = bytes.Index(len(s.data))
	return &s
}

func (s *Scanner) File() *fs.File {
	return s.file
}

// Next reads japi file by bytes, detects lexemes beginnings and ends and returns them as soon as they found
// returns false for the end of file
func (s *Scanner) Next() (*Lexeme, *jerr.JApiError) {
	if len(s.finds) != 0 { // found beginning or end of lexeme
		lex, je := s.processLexemeEvent(s.shiftFound())
		if je != nil {
			return nil, je
		}
		if lex != nil {
			return lex, nil
		}
	}

	for s.curIndex <= s.dataSize {
		var c byte
		if s.curIndex == s.dataSize { // file ended
			// s.closeAllOpenedLexemeEvent()
			c = EOF
		} else {
			c = s.data[s.curIndex]
			if c == EOF { // we use it to imitate end of file but strictly under our control
				return nil, s.japiErrorBasic("File cannot contain byte zero")
			}
		}

		je := s.step(s, c) // evaluate byte
		if je != nil {
			return nil, je
		}
		s.curIndex++

		for range s.finds { // sometimes one char means two finds
			lex, je := s.processLexemeEvent(s.shiftFound())
			if je != nil {
				return nil, je
			}
			if lex != nil {
				switch lex.Type() {
				case Parameter:
					s.lastDirectiveParameters = append(s.lastDirectiveParameters, lex)
				case Keyword:
					s.lastDirectiveParameters = s.lastDirectiveParameters[:0]
				default:
					// none
				}
				return lex, nil
			}
		}
	}

	return nil, nil
}

func (s *Scanner) CurrentIndex() bytes.Index {
	return s.curIndex
}

// SetCurrentIndex to continue scanning from certain position
func (s *Scanner) SetCurrentIndex(i bytes.Index) {
	s.curIndex = i
}

// it is important to rely here on Event index only, not current scanner index
func (s *Scanner) processLexemeEvent(lexEvent LexemeEvent) (*Lexeme, *jerr.JApiError) {
	eventType := lexEvent.type_
	switch {
	case eventType.IsBeginning():
		s.stack.Push(lexEvent)
		return nil, nil
	case eventType.IsEnding():
		startEvent := s.stack.Pop()
		startType := startEvent.type_
		switch {
		case startType == KeywordBegin && eventType == KeywordEnd,
			startType == AnnotationBegin && eventType == AnnotationEnd,
			startType == SchemaBegin && eventType == SchemaEnd,
			startType == TextBegin && eventType == TextEnd,
			startType == ParameterBegin && eventType == ParameterEnd,
			startType == EnumBegin && eventType == EnumEnd:

			lex := NewLexeme(eventType.ToLexemeType(), startEvent.position, lexEvent.position, s.file)

			return lex, nil
		default:
			return nil, s.japiErrorBasic("Ending lexeme event does not match beginning event")
		}
	case eventType.IsSingle():
		lex := NewLexeme(eventType.ToLexemeType(), lexEvent.position, lexEvent.position, s.file)
		return lex, nil
	default:
		return nil, s.japiErrorBasic("Unsupported lexeme event type")
	}
}

func (s *Scanner) foundAt(i bytes.Index, t LexemeEventType) {
	s.finds = append(s.finds, LexemeEvent{t, i})
}

func (s *Scanner) found(t LexemeEventType) {
	s.foundAt(s.curIndex, t)
}

func (s *Scanner) shiftFound() LexemeEvent {
	length := len(s.finds)
	if length == 0 {
		panic("Empty set of found lexemes")
	}
	lexEvent := s.finds[0]
	copy(s.finds[0:], s.finds[1:])
	s.finds = s.finds[:length-1]
	return lexEvent
}

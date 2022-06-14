package scanner

import (
	"github.com/jsightapi/jsight-schema-go-library/bytes"
)

type LexemeEvent struct {
	type_    LexemeEventType
	position bytes.Index
}

type LexemeEventType uint8

const (
	KeywordBegin LexemeEventType = iota
	KeywordEnd
	ParameterBegin
	ParameterEnd
	AnnotationBegin
	AnnotationEnd
	SchemaBegin
	SchemaEnd
	JsonArrayBegin
	JsonArrayEnd
	TextBegin
	TextEnd
	ContextOpen
	ContextClose
)

func (e LexemeEventType) IsBeginning() bool {
	switch e {
	case KeywordBegin, AnnotationBegin, SchemaBegin, TextBegin, JsonArrayBegin, ParameterBegin:
		return true
	default:
		return false
	}
}

func (e LexemeEventType) IsEnding() bool {
	switch e {
	case KeywordEnd, AnnotationEnd, SchemaEnd, TextEnd, JsonArrayEnd, ParameterEnd:
		return true
	default:
		return false
	}
}

func (e LexemeEventType) IsSingle() bool {
	switch e {
	case ContextOpen, ContextClose:
		return true
	default:
		return false
	}
}

func (e LexemeEventType) String() string {
	switch e {
	case KeywordBegin:
		return "keyword-begin"
	case KeywordEnd:
		return "keyword-end"
	case AnnotationBegin:
		return "annotation-begin"
	case AnnotationEnd:
		return "annotation-end"
	case SchemaBegin:
		return "schema-begin"
	case SchemaEnd:
		return "schema-end"
	case JsonArrayBegin:
		return "array-begin"
	case JsonArrayEnd:
		return "array-end"
	case TextBegin:
		return "text-begin"
	case TextEnd:
		return "text-end"
	case ContextOpen:
		return "context-open"
	case ContextCloseSign:
		return "context-close"
	case ParameterBegin:
		return "property-begin"
	case ParameterEnd:
		return "property-end"
	default:
		return "Unknown-lexeme-event-type"
	}
}

func (e LexemeEventType) ToLexemeType() LexemeType {
	switch e {
	case KeywordBegin, KeywordEnd:
		return Keyword
	case AnnotationBegin, AnnotationEnd:
		return Annotation
	case SchemaBegin, SchemaEnd:
		return Schema
	case TextBegin, TextEnd:
		return Text
	case JsonArrayBegin, JsonArrayEnd:
		return Array
	case ContextOpen:
		return ContextExplicitOpening
	case ContextClose:
		return ContextExplicitClosing
	case ParameterBegin, ParameterEnd:
		return Parameter
	default:
		panic("Unknown lexeme event type")
	}
}

// func (e LexemeEventType) ClosingPair() LexemeEventType {
// 	switch e {
// 	case KeywordBegin:
// 		return KeywordEnd
// 	case ValueBegin:
// 		return ValueEnd
// 	case AnnotationBegin:
// 		return AnnotationEnd
// 	case SchemaBegin:
// 		return SchemaEnd
// 	case StringBodyBegin:
// 		return StringBodyEnd
// 	case JsonArrayBegin:
// 		return JsonArrayEnd
// 	case TextBegin:
// 		return TextEnd
// 	case ContextOpen:
// 		return ContextClose
// 	case ParameterBegin, ParameterEnd:
// 		return ParameterEnd
// 	default:
// 		panic("Unknown beginning lexeme")
// 	}
// }

package lexeme

type LexEventType uint8

const (
	LiteralBegin LexEventType = iota
	LiteralEnd
	ObjectBegin
	ObjectEnd
	ObjectKeyBegin
	ObjectKeyEnd
	ObjectValueBegin
	ObjectValueEnd
	ArrayBegin
	ArrayEnd
	ArrayItemBegin
	ArrayItemEnd
	InlineAnnotationBegin
	InlineAnnotationEnd
	InlineAnnotationTextBegin
	InlineAnnotationTextEnd
	MultiLineAnnotationBegin
	MultiLineAnnotationEnd
	MultiLineAnnotationTextBegin
	MultiLineAnnotationTextEnd
	NewLine

	// TypesShortcutBegin indicates that "type" or "or" shortcut was began.
	TypesShortcutBegin

	// TypesShortcutEnd indicates that "type" or "or" shortcut was ended.
	TypesShortcutEnd

	KeyShortcutBegin
	KeyShortcutEnd

	// MixedValueBegin indicates that here can be anything: scalar, array, or object.
	MixedValueBegin
	MixedValueEnd

	EndTop // character after the last closing JSON or SCHEMA lexeme event
)

// IsOneOf returns true if given lexeme type is equal to at least one of specified.
func (e LexEventType) IsOneOf(ll ...LexEventType) bool {
	for _, l := range ll {
		if l == e {
			return true
		}
	}
	return false
}

func (e LexEventType) IsOpening() bool {
	switch e { //nolint:exhaustive // It's okay.
	case LiteralBegin,
		ObjectBegin,
		ObjectKeyBegin,
		ObjectValueBegin,
		ArrayBegin,
		ArrayItemBegin,
		MultiLineAnnotationBegin,
		InlineAnnotationBegin,
		InlineAnnotationTextBegin,
		MultiLineAnnotationTextBegin,
		TypesShortcutBegin,
		KeyShortcutBegin,
		MixedValueBegin:
		return true
	}
	return false
}

func (e LexEventType) String() string { //nolint:gocyclo // todo try to make this more readable
	switch e {
	case LiteralBegin:
		return "literal-begin"
	case LiteralEnd:
		return "literal-end"
	case ObjectBegin:
		return "object-begin"
	case ObjectEnd:
		return "object-end"
	case ObjectKeyBegin:
		return "key-begin"
	case ObjectKeyEnd:
		return "key-end"
	case ObjectValueBegin:
		return "value-begin"
	case ObjectValueEnd:
		return "value-end"
	case ArrayBegin:
		return "array-begin"
	case ArrayEnd:
		return "array-end"
	case ArrayItemBegin:
		return "item-begin"
	case ArrayItemEnd:
		return "item-end"
	case InlineAnnotationBegin:
		return "inline-annotation-begin"
	case InlineAnnotationEnd:
		return "inline-annotation-end"
	case InlineAnnotationTextBegin:
		return "inline-annotation-text-begin"
	case InlineAnnotationTextEnd:
		return "inline-annotation-text-end"
	case MultiLineAnnotationBegin:
		return "multi-line-annotation-begin"
	case MultiLineAnnotationEnd:
		return "multi-line-annotation-end"
	case MultiLineAnnotationTextBegin:
		return "multi-line-annotation-text-begin"
	case MultiLineAnnotationTextEnd:
		return "multi-line-annotation-text-end"
	case NewLine:
		return "new-line"
	case TypesShortcutBegin:
		return "types-shortcut-begin"
	case TypesShortcutEnd:
		return "types-shortcut-end"
	case KeyShortcutBegin:
		return "key-shortcut-begin"
	case KeyShortcutEnd:
		return "key-shortcut-end"
	case MixedValueBegin:
		return "mixed-value-begin"
	case MixedValueEnd:
		return "mixed-value-end"
	case EndTop:
		return "end-top"
	default:
		panic("Unknown lexical event type")
	}
}

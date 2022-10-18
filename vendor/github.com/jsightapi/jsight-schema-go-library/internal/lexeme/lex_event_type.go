package lexeme

// LexEventType available lexeme event types.
// gen:Stringer e Unknown lexical event type
type LexEventType uint8

const (
	LiteralBegin                 LexEventType = iota // literal-begin
	LiteralEnd                                       // literal-end
	ObjectBegin                                      // object-begin
	ObjectEnd                                        // object-end
	ObjectKeyBegin                                   // key-begin
	ObjectKeyEnd                                     // key-end
	ObjectValueBegin                                 // value-begin
	ObjectValueEnd                                   // value-end
	ArrayBegin                                       // array-begin
	ArrayEnd                                         // array-end
	ArrayItemBegin                                   // item-begin
	ArrayItemEnd                                     // item-end
	InlineAnnotationBegin                            // inline-annotation-begin
	InlineAnnotationEnd                              // inline-annotation-end
	InlineAnnotationTextBegin                        // inline-annotation-text-begin
	InlineAnnotationTextEnd                          // inline-annotation-text-end
	MultiLineAnnotationBegin                         // multi-line-annotation-begin
	MultiLineAnnotationEnd                           // multi-line-annotation-end
	MultiLineAnnotationTextBegin                     // multi-line-annotation-text-begin
	MultiLineAnnotationTextEnd                       // multi-line-annotation-text-end
	NewLine                                          // new-line

	// TypesShortcutBegin indicates that "type" or "or" shortcut was began.
	TypesShortcutBegin // types-shortcut-begin

	// TypesShortcutEnd indicates that "type" or "or" shortcut was ended.
	TypesShortcutEnd // types-shortcut-end

	KeyShortcutBegin // key-shortcut-begin
	KeyShortcutEnd   // key-shortcut-end

	// MixedValueBegin indicates that here can be anything: scalar, array, or object.
	MixedValueBegin // mixed-value-begin
	MixedValueEnd   // mixed-value-end

	// EndTop character after the last closing JSON or SCHEMA lexeme event.
	EndTop // end-top
)

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

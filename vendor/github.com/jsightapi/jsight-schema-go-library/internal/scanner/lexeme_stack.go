package scanner

import (
	"github.com/jsightapi/jsight-schema-go-library/internal/lexeme"
)

// LexemesStack a stack of lexemes.
type LexemesStack []lexeme.LexEvent

// NewLexemesStack creates a stack of lexemes.
func NewLexemesStack() LexemesStack {
	return make(LexemesStack, 0, 10)
}

// Len returns length of the stack.
func (stack *LexemesStack) Len() int {
	return len(*stack)
}

// Push pushes a lexeme onto the stack.
func (stack *LexemesStack) Push(lex lexeme.LexEvent) {
	*stack = append(*stack, lex)
}

// Pop pops a lexeme from the stack.
func (stack *LexemesStack) Pop() lexeme.LexEvent {
	lex := stack.Peek()
	count := len(*stack)
	*stack = (*stack)[:count-1]
	return lex
}

// Peek returns lexeme of the stack from the end, without removing.
func (stack *LexemesStack) Peek() lexeme.LexEvent {
	count := len(*stack)
	if count == 0 {
		panic("Reading from empty stack")
	}
	return (*stack)[count-1]
}

// Get returns lexeme of the stack, without removing.
func (stack *LexemesStack) Get(i int) lexeme.LexEvent {
	count := len(*stack)
	if i > count-1 {
		panic("Reading a nonexistent element of the stack")
	}
	return (*stack)[i]
}

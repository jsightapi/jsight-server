package scanner

type eventStack []LexemeEvent

func (stack eventStack) Len() int {
	return len(stack)
}

func (stack *eventStack) Push(lex LexemeEvent) {
	*stack = append(*stack, lex)
}

func (stack *eventStack) Pop() LexemeEvent {
	lex := stack.Peek()
	count := len(*stack)
	*stack = (*stack)[:count-1]
	return lex
}

func (stack *eventStack) Peek() LexemeEvent {
	count := len(*stack)
	if count == 0 {
		panic("Reading from empty stack")
	}
	return (*stack)[count-1]
}

// Get return lexeme of the stack, without removing
func (stack eventStack) Get(i int) LexemeEvent {
	count := len(stack)
	if i > count-1 {
		panic("Reading a nonexistent lexeme of the stack")
	}
	return stack[i]
}

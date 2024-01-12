package scanner

type eventStack []LexemeEvent

func (stack *eventStack) Push(lex LexemeEvent) {
	*stack = append(*stack, lex)
}

func (stack *eventStack) Pop() LexemeEvent {
	lex := stack.peek()
	count := len(*stack)
	*stack = (*stack)[:count-1]
	return lex
}

func (stack *eventStack) peek() LexemeEvent {
	count := len(*stack)
	if count == 0 {
		panic("Reading from empty stack")
	}
	return (*stack)[count-1]
}

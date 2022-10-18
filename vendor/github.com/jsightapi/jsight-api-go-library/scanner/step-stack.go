package scanner

type stepFuncStack []stepFunc // tod taken

func (stack *stepFuncStack) Push(val stepFunc) {
	*stack = append(*stack, val)
}

func (stack *stepFuncStack) Pop() stepFunc {
	f := stack.peek()
	count := len(*stack)
	*stack = (*stack)[:count-1]
	return f
}

func (stack *stepFuncStack) peek() stepFunc {
	count := len(*stack)
	if count == 0 {
		panic("Reading from empty stack")
	}
	return (*stack)[count-1]
}

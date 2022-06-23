package scanner

// Stack of step functions, to preserve the sequence of steps (and return to them) in some cases.
type stepFuncStack []stepFunc

func (stack *stepFuncStack) Push(val stepFunc) {
	*stack = append(*stack, val)
}

func (stack *stepFuncStack) Pop() stepFunc {
	f := stack.Peek()
	count := len(*stack)
	*stack = (*stack)[:count-1]
	return f
}

func (stack *stepFuncStack) Peek() stepFunc {
	count := len(*stack)
	if count == 0 {
		panic("Reading from empty stack")
	}
	return (*stack)[count-1]
}

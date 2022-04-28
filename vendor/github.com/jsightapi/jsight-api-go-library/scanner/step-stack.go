package scanner

type stepFuncStack []stepFunc // tod taken

func (stack *stepFuncStack) Len() int {
	return len(*stack)
}

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

// Get return element of the stack, without removing
func (stack *stepFuncStack) Get(i int) stepFunc {
	count := len(*stack)
	if i > count-1 {
		panic("Reading a nonexistent stepFunc of the stack")
	}
	return (*stack)[i]
}

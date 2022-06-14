package scanner

// Stack of step functions, to preserve the sequence of steps (and return to them) in some cases.
type contextStack []context

func (s *contextStack) Push(v context) {
	*s = append(*s, v)
}

func (s *contextStack) Pop() context {
	v := s.Peek()
	*s = (*s)[:len(*s)-1]
	return v
}

func (s *contextStack) Peek() context {
	l := len(*s)
	if l == 0 {
		panic("Reading from empty stack")
	}

	return (*s)[l-1]
}

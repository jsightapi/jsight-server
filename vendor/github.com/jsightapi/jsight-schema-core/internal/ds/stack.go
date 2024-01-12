package ds

import "github.com/jsightapi/jsight-schema-core/errs"

// Stack represent generic stack.
// Not a thread safe!
type Stack[T any] struct {
	vals []T
}

// Len returns length of the stack.
func (s *Stack[T]) Len() int {
	if s == nil {
		return 0
	}
	return len(s.vals)
}

// Push pushes a value onto the stack.
func (s *Stack[T]) Push(v T) {
	s.vals = append(s.vals, v)
}

// Pop pops a value from the stack.
func (s *Stack[T]) Pop() T {
	lex := s.Peek()
	s.vals = s.vals[:s.Len()-1]
	return lex
}

// Peek returns a value of the stack from the end, without removing.
func (s *Stack[T]) Peek() T {
	l := s.Len()
	if l == 0 {
		panic(errs.ErrRuntimeFailure.F())
	}
	return s.vals[l-1]
}

// Get returns a value of the stack, without removing.
func (s *Stack[T]) Get(i int) T {
	if i < 0 || i > s.Len()-1 {
		panic(errs.ErrRuntimeFailure.F())
	}
	return s.vals[i]
}

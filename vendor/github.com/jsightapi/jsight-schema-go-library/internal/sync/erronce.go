package sync

import "sync"

// ErrOnce we same as sync.Once but require a function which can return an error.
// This error will be hold inside this type and return every time when someone
// call `Do` method.
type ErrOnce struct {
	err  error
	once sync.Once
}

// Do doing the stuff.
func (e *ErrOnce) Do(fn func() error) error {
	e.once.Do(func() {
		e.err = fn()
	})
	return e.err
}

// ErrOnceWithValue we same as ErrOnce but holds the value as well.
type ErrOnceWithValue[T any] struct {
	value T
	err   error
	once  sync.Once
}

// Do doing the stuff.
func (e *ErrOnceWithValue[T]) Do(fn func() (T, error)) (T, error) {
	e.once.Do(func() {
		e.value, e.err = fn()
	})
	return e.value, e.err
}

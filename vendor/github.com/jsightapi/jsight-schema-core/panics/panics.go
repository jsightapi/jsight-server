package panics

// Handle handles panics properly.
func Handle(r interface{}, originErr error) error {
	if originErr != nil {
		return originErr
	}

	if r == nil {
		return nil
	}

	rErr, ok := r.(error)
	if !ok {
		panic(r)
	}
	return rErr
}

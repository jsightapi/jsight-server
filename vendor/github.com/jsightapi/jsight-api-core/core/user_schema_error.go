package core

type userTypeError struct {
	userTypeName string
	err          error
}

func (e userTypeError) Error() string {
	return e.err.Error()
}

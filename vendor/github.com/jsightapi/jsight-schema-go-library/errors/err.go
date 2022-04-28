package errors

type Err interface {
	error

	Code() ErrorCode
}

type Error interface {
	Filename() string
	Position() uint
	Message() string
	ErrCode() int
	IncorrectUserType() string
}

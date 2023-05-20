package errs

type CodeKeeper interface {
	error
	Code() Code
}

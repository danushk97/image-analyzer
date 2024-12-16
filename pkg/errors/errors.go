package errors

type IError interface {
	error
	Cause() error
	Wrap(err error) IError
}

// AppError is the error type with a string type and a parent error class.
type AppError struct {
	message string
	cause   error
}

// NewAppError creates a new error with the given message and class.
func NewAppError(message string) IError {
	return AppError{
		message: message,
	}
}

// AppError returns the error message.
func (e AppError) Error() string {
	return e.message
}

func (e AppError) Wrap(err error) IError {
	e.cause = err

	return e
}

func (e AppError) Cause() error {
	return e.cause
}

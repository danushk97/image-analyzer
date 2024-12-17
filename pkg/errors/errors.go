package errors

type IError interface {
	error
	Cause() error
	Wrap(err error) IError
	IsOfType(ErrorType) bool
}

// AppError is the error type with a string type and a parent error class.
type AppError struct {
	message string
	cause   error
	Type    ErrorType
}

// NewAppError creates a new error with the given message and class.
func NewBadRequestError(message string) IError {
	return AppError{
		message: message,
		Type:    BAD_REQUEST_ERROR,
	}
}

// NewAppError creates a new error with the given message and class.
func NewServerError(message string) IError {
	return AppError{
		message: message,
		Type:    INTERNAL_SERVER_ERROR,
	}
}

// NewAppError creates a new error with the given message and class.
func NewAuthorizationError(message string) IError {
	return AppError{
		message: message,
		Type:    AUTHORIZATION_ERROR,
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

func (e AppError) IsOfType(eType ErrorType) bool {
	return e.Type == eType
}

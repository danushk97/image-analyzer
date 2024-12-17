package errors

type ErrorType string

const (
	BAD_REQUEST_ERROR     ErrorType = "BAD_REQUEST_ERROR"
	INTERNAL_SERVER_ERROR ErrorType = "INTERNAL_SERVER_ERROR"
	AUTHORIZATION_ERROR   ErrorType = "AUTHORIZATION_ERROR"
)

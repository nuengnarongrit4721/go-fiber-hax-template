package errors

import "errors"

var (
	ErrNotFound       = errors.New("resource not found")
	ErrConflict       = errors.New("resource already exists")
	ErrInvalidInput   = errors.New("invalid input data")
	ErrUnauthorized   = errors.New("unauthorized access")
	ErrForbidden      = errors.New("forbidden resource coverage")
	ErrInternalServer = errors.New("internal server error")
	ErrDatabase       = errors.New("database operation failed")
)

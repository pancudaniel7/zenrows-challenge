package apperr

import "fmt"

type appError struct {
	code  string
	msg   string
	cause error
}

func newAppError(code, msg string, cause error) appError {
	return appError{code: code, msg: msg, cause: cause}
}

func (e appError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.code, e.msg, e.cause)
	}
	return fmt.Sprintf("[%s] %s", e.code, e.msg)
}

func (e appError) Code() string      { return e.code }
func (e appError) Message() string   { return e.msg }
func (e appError) CauseError() error { return e.cause }
func (e appError) Unwrap() error     { return e.cause }

type InvalidArgErr struct{ appError }

// NewInvalidArgErr builds an INVALID_ARGUMENT Error for bad client input.
func NewInvalidArgErr(msg string, cause error) *InvalidArgErr {
	return &InvalidArgErr{appError: newAppError("INVALID_ARGUMENT", msg, cause)}
}

// Error renders the InvalidArgErr as a string.
func (e *InvalidArgErr) Error() string { return e.appError.Error() }

type NotFoundErr struct{ appError }

// NewNotFoundErr builds a NOT_FOUND Error when a resource is missing.
func NewNotFoundErr(msg string, cause error) *NotFoundErr {
	return &NotFoundErr{appError: newAppError("NOT_FOUND", msg, cause)}
}

// Error renders the NotFoundErr as a string.
func (e *NotFoundErr) Error() string { return e.appError.Error() }

type AlreadyExistsErr struct{ appError }

// NewAlreadyExistsErr builds an ALREADY_EXISTS Error for duplicate resources.
func NewAlreadyExistsErr(msg string, cause error) *AlreadyExistsErr {
	return &AlreadyExistsErr{appError: newAppError("ALREADY_EXISTS", msg, cause)}
}

// Error renders the AlreadyExistsErr as a string.
func (e *AlreadyExistsErr) Error() string { return e.appError.Error() }

type NotAuthorizedErr struct{ appError }

// NewNotAuthorizedErr builds a NOT_AUTHORIZED Error for failed authz/authn.
func NewNotAuthorizedErr(msg string, cause error) *NotAuthorizedErr {
	return &NotAuthorizedErr{appError: newAppError("NOT_AUTHORIZED", msg, cause)}
}

// Error renders the NotAuthorizedErr as a string.
func (e *NotAuthorizedErr) Error() string { return e.appError.Error() }

type InternalErr struct{ appError }

// NewInternalErr builds an INTERNAL_ERROR for unexpected failures.
func NewInternalErr(msg string, cause error) *InternalErr {
	return &InternalErr{appError: newAppError("INTERNAL_ERROR", msg, cause)}
}

// Error renders the InternalErr as a string.
func (e *InternalErr) Error() string { return e.appError.Error() }

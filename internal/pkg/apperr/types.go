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

func (e appError) error() string {
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

func NewInvalidArgErr(msg string, cause error) *InvalidArgErr {
	return &InvalidArgErr{appError: newAppError("INVALID_ARGUMENT", msg, cause)}
}

func (e *InvalidArgErr) Error() string { return e.appError.error() }

type NotFoundErr struct{ appError }

func NewNotFoundErr(msg string, cause error) *NotFoundErr {
	return &NotFoundErr{appError: newAppError("NOT_FOUND", msg, cause)}
}

func (e *NotFoundErr) Error() string { return e.appError.error() }

type AlreadyExistsErr struct{ appError }

func NewAlreadyExistsErr(msg string, cause error) *AlreadyExistsErr {
	return &AlreadyExistsErr{appError: newAppError("ALREADY_EXISTS", msg, cause)}
}

func (e *AlreadyExistsErr) Error() string { return e.appError.error() }

type NotAuthorizedErr struct{ appError }

func NewNotAuthorizedErr(msg string, cause error) *NotAuthorizedErr {
	return &NotAuthorizedErr{appError: newAppError("NOT_AUTHORIZED", msg, cause)}
}

func (e *NotAuthorizedErr) Error() string { return e.appError.error() }

type InternalErr struct{ appError }

func NewInternalErr(msg string, cause error) *InternalErr {
	return &InternalErr{appError: newAppError("INTERNAL_ERROR", msg, cause)}
}

func (e *InternalErr) Error() string { return e.appError.error() }

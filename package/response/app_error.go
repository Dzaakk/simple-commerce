package response

import (
	"fmt"
)

// AppError is a typed error that carries an HTTP status code and safe client message.
// Use this in service layer so handlers can simply ctx.Error(err).
type AppError struct {
	Code    int
	Message string
	Details interface{}
	Err     error
}

func (e *AppError) Error() string {
	if e == nil {
		return ""
	}
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

func NewAppError(code int, message string, details ...interface{}) *AppError {
	var detail interface{}
	if len(details) > 0 {
		detail = details[0]
	}
	return &AppError{
		Code:    code,
		Message: message,
		Details: detail,
	}
}

func WrapAppError(code int, message string, err error, details ...interface{}) *AppError {
	appErr := NewAppError(code, message, details...)
	appErr.Err = err
	return appErr
}

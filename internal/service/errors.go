package service

import "fmt"

type ErrorCode string

const (
	ErrorCodeValidationFailed       ErrorCode = "validation_failed"
	ErrorCodeEmailAlreadyRegistered ErrorCode = "email_already_registered"
	ErrorCodeInvalidCredentials     ErrorCode = "invalid_credentials"
	ErrorCodeInvalidRefreshToken    ErrorCode = "invalid_refresh_token"
	ErrorCodeUnauthorized           ErrorCode = "unauthorized"
	ErrorCodeForbidden              ErrorCode = "forbidden"
	ErrorCodeNotFound               ErrorCode = "not_found"
)

type Error struct {
	Code    ErrorCode
	Message string
	Details map[string]any
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func NewError(code ErrorCode, message string, details map[string]any) *Error {
	return &Error{Code: code, Message: message, Details: details}
}

package errors

import (
	"errors"
	"fmt"
)

// Error codes
type Code string

const (
	CodeUnknown       Code = "UNKNOWN"
	CodeNotFound      Code = "NOT_FOUND"
	CodeAlreadyExists Code = "ALREADY_EXISTS"
	CodeInvalidInput  Code = "INVALID_INPUT"
	CodeDBError       Code = "DB_ERROR"
	CodeClaudeError   Code = "CLAUDE_ERROR"
	CodeTimeout       Code = "TIMEOUT"
	CodeCancelled     Code = "CANCELLED"
	CodePermission    Code = "PERMISSION"
	CodeConfig        Code = "CONFIG"
)

// Error is a structured error with code
type Error struct {
	Code    Code
	Message string
	Err     error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *Error) Unwrap() error {
	return e.Err
}

// New creates a new error with code
func New(code Code, message string) *Error {
	return &Error{Code: code, Message: message}
}

// Wrap wraps an existing error with code
func Wrap(code Code, message string, err error) *Error {
	return &Error{Code: code, Message: message, Err: err}
}

// Convenience constructors

func NotFound(resource string) *Error {
	return New(CodeNotFound, fmt.Sprintf("%s not found", resource))
}

func AlreadyExists(resource string) *Error {
	return New(CodeAlreadyExists, fmt.Sprintf("%s already exists", resource))
}

func InvalidInput(message string) *Error {
	return New(CodeInvalidInput, message)
}

func DBError(err error) *Error {
	return Wrap(CodeDBError, "database error", err)
}

func ClaudeError(err error) *Error {
	return Wrap(CodeClaudeError, "claude execution error", err)
}

func Timeout(message string) *Error {
	return New(CodeTimeout, message)
}

func Cancelled(message string) *Error {
	return New(CodeCancelled, message)
}

func ConfigError(message string) *Error {
	return New(CodeConfig, message)
}

// IsCode checks if error has specific code
func IsCode(err error, code Code) bool {
	var e *Error
	if errors.As(err, &e) {
		return e.Code == code
	}
	return false
}

// GetCode returns error code or CodeUnknown
func GetCode(err error) Code {
	var e *Error
	if errors.As(err, &e) {
		return e.Code
	}
	return CodeUnknown
}

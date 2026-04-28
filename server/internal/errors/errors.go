package errors

import (
	"fmt"
	"net/http"
)

// ErrorType defines the type of error for categorization
type ErrorType string

const (
	ValidationError ErrorType = "VALIDATION_ERROR"
	NotFoundError   ErrorType = "NOT_FOUND_ERROR"
	ConflictError   ErrorType = "CONFLICT_ERROR"
	AuthError       ErrorType = "AUTH_ERROR"
	ServerError     ErrorType = "SERVER_ERROR"
	ForbiddenError  ErrorType = "FORBIDDEN_ERROR"
)

// AppError is the standard error type for the application
type AppError struct {
	Type       ErrorType
	Message    string
	StatusCode int
	Err        error
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Type, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewValidationError creates a new validation error
func NewValidationError(message string) *AppError {
	return &AppError{
		Type:       ValidationError,
		Message:    message,
		StatusCode: http.StatusBadRequest,
	}
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(resource string) *AppError {
	return &AppError{
		Type:       NotFoundError,
		Message:    fmt.Sprintf("%s not found", resource),
		StatusCode: http.StatusNotFound,
	}
}

// NewConflictError creates a new conflict error
func NewConflictError(message string) *AppError {
	return &AppError{
		Type:       ConflictError,
		Message:    message,
		StatusCode: http.StatusConflict,
	}
}

// NewAuthError creates a new authentication error
func NewAuthError(message string) *AppError {
	return &AppError{
		Type:       AuthError,
		Message:    message,
		StatusCode: http.StatusUnauthorized,
	}
}

// NewForbiddenError creates a new forbidden error
func NewForbiddenError(message string) *AppError {
	return &AppError{
		Type:       ForbiddenError,
		Message:    message,
		StatusCode: http.StatusForbidden,
	}
}

// NewServerError creates a new server error
func NewServerError(message string, err error) *AppError {
	return &AppError{
		Type:       ServerError,
		Message:    message,
		StatusCode: http.StatusInternalServerError,
		Err:        err,
	}
}

// Wrap wraps an error with context
func Wrap(err error, message string) *AppError {
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}
	return NewServerError(message, err)
}

// ErrorResponse is the JSON response for errors
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	TraceID string `json:"trace_id,omitempty"`
}

package errors

import (
	"fmt"
	"time"
)

// AppError represents an application error
type AppError struct {
	Code      string         `json:"code"`
	Message   string         `json:"message"`
	Details   map[string]any `json:"details,omitempty"`
	Cause     error          `json:"-"`
	Timestamp time.Time      `json:"timestamp"`
	RequestID string         `json:"request_id,omitempty"`
	TraceID   string         `json:"trace_id,omitempty"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Cause
}

// NewValidationError creates a new validation error
func NewValidationError(code, message string, details map[string]any) *AppError {
	return &AppError{
		Code:      code,
		Message:   message,
		Details:   details,
		Timestamp: time.Now(),
	}
}

// NewBusinessError creates a new business error
func NewBusinessError(code, message string, details map[string]any) *AppError {
	return &AppError{
		Code:      code,
		Message:   message,
		Details:   details,
		Timestamp: time.Now(),
	}
}

// NewSystemError creates a new system error
func NewSystemError(code, message string, details map[string]any) *AppError {
	return &AppError{
		Code:      code,
		Message:   message,
		Details:   details,
		Timestamp: time.Now(),
	}
}

// WrapError wraps an existing error with additional context
func WrapError(err error, code, message string) *AppError {
	return &AppError{
		Code:      code,
		Message:   message,
		Cause:     err,
		Timestamp: time.Now(),
	}
}

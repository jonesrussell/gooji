package video

import (
	"fmt"
	"net/http"
)

// ErrorType represents the type of error
type ErrorType string

const (
	// ErrorTypeValidation represents validation errors
	ErrorTypeValidation ErrorType = "validation"
	// ErrorTypeNotFound represents not found errors
	ErrorTypeNotFound ErrorType = "not_found"
	// ErrorTypeInternal represents internal server errors
	ErrorTypeInternal ErrorType = "internal"
	// ErrorTypeSecurity represents security-related errors
	ErrorTypeSecurity ErrorType = "security"
	// ErrorTypeUpload represents upload-related errors
	ErrorTypeUpload ErrorType = "upload"
)

// VideoError represents a structured error with context
type VideoError struct {
	Type    ErrorType `json:"type"`
	Message string    `json:"message"`
	Code    int       `json:"code"`
	Err     error     `json:"-"`
}

// Error implements the error interface
func (e *VideoError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the underlying error
func (e *VideoError) Unwrap() error {
	return e.Err
}

// NewValidationError creates a new validation error
func NewValidationError(message string, err error) *VideoError {
	return &VideoError{
		Type:    ErrorTypeValidation,
		Message: message,
		Code:    http.StatusBadRequest,
		Err:     err,
	}
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(message string, err error) *VideoError {
	return &VideoError{
		Type:    ErrorTypeNotFound,
		Message: message,
		Code:    http.StatusNotFound,
		Err:     err,
	}
}

// NewInternalError creates a new internal server error
func NewInternalError(message string, err error) *VideoError {
	return &VideoError{
		Type:    ErrorTypeInternal,
		Message: message,
		Code:    http.StatusInternalServerError,
		Err:     err,
	}
}

// NewSecurityError creates a new security error
func NewSecurityError(message string, err error) *VideoError {
	return &VideoError{
		Type:    ErrorTypeSecurity,
		Message: message,
		Code:    http.StatusForbidden,
		Err:     err,
	}
}

// NewUploadError creates a new upload error
func NewUploadError(message string, err error) *VideoError {
	return &VideoError{
		Type:    ErrorTypeUpload,
		Message: message,
		Code:    http.StatusBadRequest,
		Err:     err,
	}
}

// IsValidationError checks if an error is a validation error
func IsValidationError(err error) bool {
	if err != nil && err.Error() != "" {
		if videoErr, ok := err.(*VideoError); ok {
			return videoErr.Type == ErrorTypeValidation
		}
	}
	return false
}

// IsNotFoundError checks if an error is a not found error
func IsNotFoundError(err error) bool {
	if err != nil && err.Error() != "" {
		if videoErr, ok := err.(*VideoError); ok {
			return videoErr.Type == ErrorTypeNotFound
		}
	}
	return false
}

// IsSecurityError checks if an error is a security error
func IsSecurityError(err error) bool {
	if err != nil && err.Error() != "" {
		if videoErr, ok := err.(*VideoError); ok {
			return videoErr.Type == ErrorTypeSecurity
		}
	}
	return false
}

// GetHTTPStatusCode returns the appropriate HTTP status code for an error
func GetHTTPStatusCode(err error) int {
	if err != nil && err.Error() != "" {
		if videoErr, ok := err.(*VideoError); ok {
			return videoErr.Code
		}
	}
	return http.StatusInternalServerError
}

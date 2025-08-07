package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Error represents an application error
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

// New creates a new error
func New(code int, message string, err error) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Error implements the error interface
func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// WriteError writes an error response to the HTTP response writer
func WriteError(w http.ResponseWriter, err error) {
	var appErr *Error
	if e, ok := err.(*Error); ok {
		appErr = e
	} else {
		appErr = New(http.StatusInternalServerError, "Internal server error", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.Code)
	if err := json.NewEncoder(w).Encode(appErr); err != nil {
		// If we can't encode the error, write a simple error message
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// Common error codes
const (
	ErrCodeBadRequest         = 400
	ErrCodeUnauthorized       = 401
	ErrCodeForbidden          = 403
	ErrCodeNotFound           = 404
	ErrCodeMethodNotAllowed   = 405
	ErrCodeInternalServer     = 500
	ErrCodeServiceUnavailable = 503
)

// Common errors
var (
	ErrBadRequest         = New(ErrCodeBadRequest, "Bad request", nil)
	ErrUnauthorized       = New(ErrCodeUnauthorized, "Unauthorized", nil)
	ErrForbidden          = New(ErrCodeForbidden, "Forbidden", nil)
	ErrNotFound           = New(ErrCodeNotFound, "Not found", nil)
	ErrMethodNotAllowed   = New(ErrCodeMethodNotAllowed, "Method not allowed", nil)
	ErrInternalServer     = New(ErrCodeInternalServer, "Internal server error", nil)
	ErrServiceUnavailable = New(ErrCodeServiceUnavailable, "Service unavailable", nil)
)

package domain

import "errors"

var (
	ErrVideoNotFound    = errors.New("video not found")
	ErrDeviceNotFound   = errors.New("device not found")
	ErrInvalidAPIKey    = errors.New("invalid API key")
	ErrInvalidStatus    = errors.New("invalid status transition")
	ErrValidationFailed = errors.New("validation failed")
)

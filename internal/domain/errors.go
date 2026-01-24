package domain

import "errors"

var (
	ErrVideoNotFound    = errors.New("video not found")
	ErrInvalidVideo     = errors.New("invalid video")
	ErrInvalidInput     = errors.New("invalid input")
	ErrStorageFailed    = errors.New("storage operation failed")
	ErrProcessingFailed = errors.New("video processing failed")
)

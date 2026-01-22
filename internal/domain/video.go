package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type VideoStatus string

const (
	StatusPending   VideoStatus = "pending"
	StatusPublished VideoStatus = "published"
	StatusRejected  VideoStatus = "rejected"
)

type Video struct {
	ID              string
	Title           string
	Description     string
	Language        string
	DurationSeconds int
	FilePath        string
	ThumbnailPath   string
	FileSizeBytes   int64
	MimeType        string
	Status          VideoStatus
	RejectionReason string
	DeviceID        *string
	CreatedAt       time.Time
	PublishedAt     *time.Time
	ModeratedAt     *time.Time
}

func NewVideo(title, filePath, mimeType string) *Video {
	now := time.Now()
	return &Video{
		ID:        uuid.New().String(),
		Title:     title,
		FilePath:  filePath,
		MimeType:  mimeType,
		Language:  "oji",
		Status:    StatusPending,
		CreatedAt: now,
	}
}

func (v *Video) Publish() error {
	if v.Status == StatusRejected {
		return errors.New("cannot publish rejected video")
	}
	now := time.Now()
	v.Status = StatusPublished
	v.PublishedAt = &now
	v.ModeratedAt = &now
	return nil
}

func (v *Video) Reject(reason string) error {
	if v.Status == StatusPublished {
		return errors.New("cannot reject published video")
	}
	now := time.Now()
	v.Status = StatusRejected
	v.RejectionReason = reason
	v.ModeratedAt = &now
	return nil
}

func (v *Video) IsPending() bool {
	return v.Status == StatusPending
}

func (v *Video) IsPublished() bool {
	return v.Status == StatusPublished
}

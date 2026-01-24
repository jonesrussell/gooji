package domain

import "time"

type Video struct {
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	Language        string    `json:"language"`
	DurationSeconds float64   `json:"duration_seconds"`
	FilePath        string    `json:"-"`
	ThumbnailPath   string    `json:"-"`
	FileSizeBytes   int64     `json:"file_size_bytes"`
	MimeType        string    `json:"mime_type"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type Tag struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type VideoListOptions struct {
	Limit  int
	Offset int
}

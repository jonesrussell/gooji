package video

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"gooji/internal/logger"
	"gooji/pkg/ffmpeg"
)

// Service defines the interface for video processing operations
type Service interface {
	ProcessUpload(ctx context.Context, file multipart.File, header *multipart.FileHeader, metadata *UploadMetadata) (*VideoMetadata, error)
	GetVideo(ctx context.Context, id string) (*VideoMetadata, error)
	ListVideos(ctx context.Context) ([]VideoMetadata, error)
	DeleteVideo(ctx context.Context, id string) error
	GenerateThumbnail(ctx context.Context, videoPath string) error
}

// Repository defines the interface for data persistence operations
type Repository interface {
	SaveVideo(ctx context.Context, file multipart.File, filename string) (string, error)
	SaveMetadata(ctx context.Context, metadata *VideoMetadata) error
	GetMetadata(ctx context.Context, id string) (*VideoMetadata, error)
	ListMetadata(ctx context.Context) ([]VideoMetadata, error)
	DeleteVideo(ctx context.Context, id string) error
	VideoExists(ctx context.Context, id string) bool
}

// Processor defines the interface for video processing operations
type Processor interface {
	GetVideoInfo(inputPath string) (*ffmpeg.VideoInfo, error)
	GenerateThumbnail(inputPath, outputPath string, timestamp float64) error
	ValidateVideo(inputPath string) error
}

// VideoMetadata represents metadata for a recorded video
type VideoMetadata struct {
	ID          string    `json:"id"`
	Filename    string    `json:"filename"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Duration    float64   `json:"duration"`
	CreatedAt   time.Time `json:"created_at"`
	Tags        []string  `json:"tags"`
}

// UploadMetadata represents metadata for video uploads
type UploadMetadata struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

// VideoInfo contains metadata about a video file
type VideoInfo ffmpeg.VideoInfo

// service implements the Service interface
type service struct {
	repo      Repository
	processor Processor
	logger    *logger.Logger
}

// NewService creates a new video service
func NewService(repo Repository, processor Processor, logger *logger.Logger) Service {
	return &service{
		repo:      repo,
		processor: processor,
		logger:    logger,
	}
}

// ProcessUpload handles the complete video upload process
func (s *service) ProcessUpload(ctx context.Context, file multipart.File, header *multipart.FileHeader, metadata *UploadMetadata) (*VideoMetadata, error) {
	// Validate upload
	if err := s.validateUpload(file, header); err != nil {
		return nil, fmt.Errorf("upload validation failed: %w", err)
	}

	// Generate secure filename
	filename := s.generateSecureFilename(header.Filename)

	// Save video file
	videoPath, err := s.repo.SaveVideo(ctx, file, filename)
	if err != nil {
		return nil, fmt.Errorf("failed to save video: %w", err)
	}

	// Get video information
	info, err := s.processor.GetVideoInfo(videoPath)
	if err != nil {
		// Clean up saved file on error
		s.repo.DeleteVideo(ctx, filename)
		return nil, fmt.Errorf("failed to get video info: %w", err)
	}

	// Create video metadata
	videoMetadata := &VideoMetadata{
		ID:          filename,
		Filename:    filename,
		Title:       s.sanitizeInput(metadata.Title),
		Description: s.sanitizeInput(metadata.Description),
		Duration:    info.Duration,
		CreatedAt:   time.Now(),
		Tags:        s.sanitizeTags(metadata.Tags),
	}

	// Save metadata
	if err := s.repo.SaveMetadata(ctx, videoMetadata); err != nil {
		// Clean up saved file on error
		s.repo.DeleteVideo(ctx, filename)
		return nil, fmt.Errorf("failed to save metadata: %w", err)
	}

	// Generate thumbnail asynchronously
	go func() {
		if err := s.GenerateThumbnail(context.Background(), videoPath); err != nil {
			s.logger.Error("Failed to generate thumbnail for %s: %v", videoPath, err)
		} else {
			s.logger.Info("Successfully generated thumbnail for: %s", videoPath)
		}
	}()

	s.logger.Info("Successfully processed video upload: %s", filename)
	return videoMetadata, nil
}

// GetVideo retrieves video metadata by ID
func (s *service) GetVideo(ctx context.Context, id string) (*VideoMetadata, error) {
	if id == "" {
		return nil, fmt.Errorf("video ID is required")
	}

	metadata, err := s.repo.GetMetadata(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get video metadata: %w", err)
	}

	return metadata, nil
}

// ListVideos retrieves all video metadata
func (s *service) ListVideos(ctx context.Context) ([]VideoMetadata, error) {
	videos, err := s.repo.ListMetadata(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list videos: %w", err)
	}

	return videos, nil
}

// DeleteVideo removes a video and its metadata
func (s *service) DeleteVideo(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("video ID is required")
	}

	if err := s.repo.DeleteVideo(ctx, id); err != nil {
		return fmt.Errorf("failed to delete video: %w", err)
	}

	s.logger.Info("Successfully deleted video: %s", id)
	return nil
}

// GenerateThumbnail creates a thumbnail for a video
func (s *service) GenerateThumbnail(ctx context.Context, videoPath string) error {
	if videoPath == "" {
		return fmt.Errorf("video path is required")
	}

	// Extract filename and create thumbnail path
	filename := filepath.Base(videoPath)
	thumbnailPath := filepath.Join("storage/thumbnails", strings.TrimSuffix(filename, filepath.Ext(filename))+".jpg")

	// Generate thumbnail at 1 second mark
	if err := s.processor.GenerateThumbnail(videoPath, thumbnailPath, 1.0); err != nil {
		return fmt.Errorf("failed to generate thumbnail: %w", err)
	}

	return nil
}

// validateUpload validates the uploaded file
func (s *service) validateUpload(file multipart.File, header *multipart.FileHeader) error {
	// Check file size (max 100MB)
	const maxSize = 100 * 1024 * 1024
	if header.Size > maxSize {
		return fmt.Errorf("file size %d exceeds maximum allowed size %d", header.Size, maxSize)
	}

	// Validate MIME type
	allowedTypes := []string{"video/mp4", "video/webm", "video/avi", "video/mov"}
	contentType := header.Header.Get("Content-Type")

	isAllowed := false
	for _, allowedType := range allowedTypes {
		if contentType == allowedType {
			isAllowed = true
			break
		}
	}

	if !isAllowed {
		return fmt.Errorf("content type %s is not allowed", contentType)
	}

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowedExts := []string{".mp4", ".webm", ".avi", ".mov"}

	isValidExt := false
	for _, allowedExt := range allowedExts {
		if ext == allowedExt {
			isValidExt = true
			break
		}
	}

	if !isValidExt {
		return fmt.Errorf("file extension %s is not allowed", ext)
	}

	// Validate file header magic bytes to detect actual file type
	if err := s.validateFileMagicBytes(file); err != nil {
		return fmt.Errorf("file type validation failed: %w", err)
	}

	return nil
}

// validateFileMagicBytes validates the file's magic bytes to ensure it's actually a video file
func (s *service) validateFileMagicBytes(file multipart.File) error {
	// Read first 12 bytes to check magic numbers
	header := make([]byte, 12)
	_, err := file.Read(header)
	if err != nil {
		return fmt.Errorf("failed to read file header: %w", err)
	}

	// Reset file position for later use
	if _, err := file.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to reset file position: %w", err)
	}

	// Check for common video file magic bytes
	// MP4: ftyp
	if len(header) >= 8 && string(header[4:8]) == "ftyp" {
		return nil
	}
	// WebM: EBML
	if len(header) >= 4 && string(header[0:4]) == "\x1a\x45\xdf\xa3" {
		return nil
	}
	// AVI: RIFF
	if len(header) >= 4 && string(header[0:4]) == "RIFF" {
		return nil
	}
	// MOV: ftyp (same as MP4)
	if len(header) >= 8 && string(header[4:8]) == "ftyp" {
		return nil
	}

	return fmt.Errorf("file does not appear to be a valid video file (invalid magic bytes)")
}

// generateSecureFilename creates a secure filename using UUID
func (s *service) generateSecureFilename(originalName string) string {
	ext := filepath.Ext(originalName)
	timestamp := time.Now().Unix()
	return fmt.Sprintf("%d_%s%s", timestamp, generateUUID(), ext)
}

// sanitizeInput sanitizes user input to prevent XSS
func (s *service) sanitizeInput(input string) string {
	// Trim whitespace and limit length
	input = strings.TrimSpace(input)
	if len(input) > 200 {
		input = input[:200]
	}

	// Basic HTML escaping (additional protection beyond template escaping)
	input = strings.ReplaceAll(input, "<", "&lt;")
	input = strings.ReplaceAll(input, ">", "&gt;")
	input = strings.ReplaceAll(input, "\"", "&quot;")
	input = strings.ReplaceAll(input, "'", "&#39;")

	return input
}

// sanitizeTags sanitizes and validates tags
func (s *service) sanitizeTags(tags []string) []string {
	if tags == nil {
		return []string{"ojibwe", "language", "culture"}
	}

	sanitized := make([]string, 0, len(tags))
	for _, tag := range tags {
		tag = s.sanitizeInput(tag)
		if tag != "" && len(tag) <= 50 {
			sanitized = append(sanitized, tag)
		}
	}

	// Ensure we have at least default tags
	if len(sanitized) == 0 {
		sanitized = []string{"ojibwe", "language", "culture"}
	}

	return sanitized
}

// generateUUID generates a simple UUID-like string
// In production, use github.com/google/uuid
func generateUUID() string {
	return fmt.Sprintf("%x", time.Now().UnixNano())
}

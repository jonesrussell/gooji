package video

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"gooji/internal/config"
	"gooji/internal/logger"
)

// repository implements the Repository interface
type repository struct {
	storage config.Storage
	logger  *logger.Logger
}

// NewRepository creates a new video repository
func NewRepository(storage config.Storage, logger *logger.Logger) Repository {
	return &repository{
		storage: storage,
		logger:  logger,
	}
}

// SaveVideo saves a video file to storage
func (r *repository) SaveVideo(ctx context.Context, file multipart.File, filename string) (string, error) {
	// Ensure uploads directory exists
	if err := os.MkdirAll(r.storage.Uploads, 0750); err != nil {
		return "", fmt.Errorf("failed to create uploads directory: %w", err)
	}

	// Create secure file path
	videoPath := filepath.Join(r.storage.Uploads, filename)

	// Validate path is within allowed directory
	if err := r.validatePath(videoPath, r.storage.Uploads); err != nil {
		return "", fmt.Errorf("invalid file path: %w", err)
	}

	// Create destination file
	dst, err := os.Create(videoPath)
	if err != nil {
		return "", fmt.Errorf("failed to create video file: %w", err)
	}
	defer dst.Close()

	// Copy file content with buffer for large files
	buffer := make([]byte, 32*1024) // 32KB buffer
	for {
		n, err := file.Read(buffer)
		if n > 0 {
			if _, writeErr := dst.Write(buffer[:n]); writeErr != nil {
				return "", fmt.Errorf("failed to write video file: %w", writeErr)
			}
		}
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return "", fmt.Errorf("failed to read video file: %w", err)
		}
	}

	r.logger.Debug("Successfully saved video file: %s", videoPath)
	return videoPath, nil
}

// SaveMetadata saves video metadata to storage
func (r *repository) SaveMetadata(ctx context.Context, metadata *VideoMetadata) error {
	// Ensure metadata directory exists
	if err := os.MkdirAll(r.storage.Metadata, 0750); err != nil {
		return fmt.Errorf("failed to create metadata directory: %w", err)
	}

	// Create metadata file path
	metadataPath := filepath.Join(r.storage.Metadata, metadata.ID+".json")

	// Validate path is within allowed directory
	if err := r.validatePath(metadataPath, r.storage.Metadata); err != nil {
		return fmt.Errorf("invalid metadata path: %w", err)
	}

	// Create metadata file
	file, err := os.Create(metadataPath)
	if err != nil {
		return fmt.Errorf("failed to create metadata file: %w", err)
	}
	defer file.Close()

	// Encode metadata as JSON
	if err := json.NewEncoder(file).Encode(metadata); err != nil {
		return fmt.Errorf("failed to encode metadata: %w", err)
	}

	r.logger.Debug("Successfully saved metadata: %s", metadataPath)
	return nil
}

// GetMetadata retrieves video metadata by ID
func (r *repository) GetMetadata(ctx context.Context, id string) (*VideoMetadata, error) {
	if id == "" {
		return nil, fmt.Errorf("metadata ID is required")
	}

	// Create metadata file path
	metadataPath := filepath.Join(r.storage.Metadata, id+".json")

	// Validate path is within allowed directory
	if err := r.validatePath(metadataPath, r.storage.Metadata); err != nil {
		return nil, fmt.Errorf("invalid metadata path: %w", err)
	}

	// Check if file exists
	if _, err := os.Stat(metadataPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("metadata not found: %s", id)
	}

	// Open metadata file
	file, err := os.Open(metadataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open metadata file: %w", err)
	}
	defer file.Close()

	// Decode metadata
	var metadata VideoMetadata
	if err := json.NewDecoder(file).Decode(&metadata); err != nil {
		return nil, fmt.Errorf("failed to decode metadata: %w", err)
	}

	return &metadata, nil
}

// ListMetadata retrieves all video metadata
func (r *repository) ListMetadata(ctx context.Context) ([]VideoMetadata, error) {
	// Ensure metadata directory exists
	if err := os.MkdirAll(r.storage.Metadata, 0750); err != nil {
		return nil, fmt.Errorf("failed to create metadata directory: %w", err)
	}

	// Read metadata directory
	files, err := os.ReadDir(r.storage.Metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata directory: %w", err)
	}

	var videos []VideoMetadata
	for _, file := range files {
		// Only process JSON files
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			metadataPath := filepath.Join(r.storage.Metadata, file.Name())

			// Validate path is within allowed directory
			if err := r.validatePath(metadataPath, r.storage.Metadata); err != nil {
				r.logger.Error("Invalid metadata path found: %s", metadataPath)
				continue
			}

			// Open and decode metadata file
			metadataFile, err := os.Open(metadataPath)
			if err != nil {
				r.logger.Error("Failed to open metadata file %s: %v", metadataPath, err)
				continue
			}

			var metadata VideoMetadata
			if err := json.NewDecoder(metadataFile).Decode(&metadata); err != nil {
				metadataFile.Close()
				r.logger.Error("Failed to decode metadata file %s: %v", metadataPath, err)
				continue
			}
			metadataFile.Close()

			videos = append(videos, metadata)
		}
	}

	return videos, nil
}

// DeleteVideo removes a video file and its metadata
func (r *repository) DeleteVideo(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("video ID is required")
	}

	// Delete video file
	videoPath := filepath.Join(r.storage.Uploads, id)
	if err := r.validatePath(videoPath, r.storage.Uploads); err == nil {
		if err := os.Remove(videoPath); err != nil && !os.IsNotExist(err) {
			r.logger.Error("Failed to delete video file %s: %v", videoPath, err)
		} else {
			r.logger.Debug("Deleted video file: %s", videoPath)
		}
	}

	// Delete metadata file
	metadataPath := filepath.Join(r.storage.Metadata, id+".json")
	if err := r.validatePath(metadataPath, r.storage.Metadata); err == nil {
		if err := os.Remove(metadataPath); err != nil && !os.IsNotExist(err) {
			r.logger.Error("Failed to delete metadata file %s: %v", metadataPath, err)
		} else {
			r.logger.Debug("Deleted metadata file: %s", metadataPath)
		}
	}

	// Delete thumbnail file
	thumbnailPath := filepath.Join(r.storage.Thumbnails, strings.TrimSuffix(id, filepath.Ext(id))+".jpg")
	if err := r.validatePath(thumbnailPath, r.storage.Thumbnails); err == nil {
		if err := os.Remove(thumbnailPath); err != nil && !os.IsNotExist(err) {
			r.logger.Error("Failed to delete thumbnail file %s: %v", thumbnailPath, err)
		} else {
			r.logger.Debug("Deleted thumbnail file: %s", thumbnailPath)
		}
	}

	return nil
}

// VideoExists checks if a video file exists
func (r *repository) VideoExists(ctx context.Context, id string) bool {
	if id == "" {
		return false
	}

	videoPath := filepath.Join(r.storage.Uploads, id)
	if err := r.validatePath(videoPath, r.storage.Uploads); err != nil {
		return false
	}

	_, err := os.Stat(videoPath)
	return err == nil
}

// validatePath ensures a file path is within the allowed directory
func (r *repository) validatePath(filePath, allowedDir string) error {
	// Resolve absolute paths
	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	absAllowedDir, err := filepath.Abs(allowedDir)
	if err != nil {
		return fmt.Errorf("failed to resolve allowed directory: %w", err)
	}

	// Check if file path is within allowed directory
	if !strings.HasPrefix(absFilePath, absAllowedDir) {
		return fmt.Errorf("file path %s is outside allowed directory %s", absFilePath, absAllowedDir)
	}

	return nil
}

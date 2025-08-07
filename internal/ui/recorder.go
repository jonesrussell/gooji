package ui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Recorder handles video recording UI operations
type Recorder struct {
	apiURL string
}

// NewRecorder creates a new video recorder UI component
func NewRecorder(apiURL string) *Recorder {
	return &Recorder{
		apiURL: apiURL,
	}
}

// validatePath ensures a file path is secure and within allowed directory
func (r *Recorder) validatePath(filePath, allowedDir string) error {
	if filePath == "" {
		return fmt.Errorf("file path cannot be empty")
	}

	// Check for path traversal attempts
	if strings.Contains(filePath, "..") {
		return fmt.Errorf("path traversal not allowed: %s", filePath)
	}

	// Check for absolute paths outside allowed directory
	if filepath.IsAbs(filePath) {
		if allowedDir != "" {
			absAllowedDir, err := filepath.Abs(allowedDir)
			if err != nil {
				return fmt.Errorf("failed to resolve allowed directory: %w", err)
			}

			absFilePath, err := filepath.Abs(filePath)
			if err != nil {
				return fmt.Errorf("failed to resolve file path: %w", err)
			}

			if !strings.HasPrefix(absFilePath, absAllowedDir) {
				return fmt.Errorf("file path %s is outside allowed directory %s", absFilePath, absAllowedDir)
			}
		}
	}

	// Check for dangerous characters
	dangerousChars := []string{"|", "&", ";", "`", "$", "(", ")", "{", "}", "[", "]", "*", "?", "\\"}
	for _, char := range dangerousChars {
		if strings.Contains(filePath, char) {
			return fmt.Errorf("dangerous character '%s' not allowed in path: %s", char, filePath)
		}
	}

	return nil
}

// UploadVideo uploads a recorded video to the server
func (r *Recorder) UploadVideo(videoPath, title, description string) error {
	// Validate video path is within allowed directory (current working directory)
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	if err := r.validatePath(videoPath, currentDir); err != nil {
		return fmt.Errorf("invalid video path: %w", err)
	}

	// Create multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add video file
	// Additional security check: ensure path is absolute and no path traversal
	if !filepath.IsAbs(videoPath) {
		return fmt.Errorf("video path must be absolute: %s", videoPath)
	}
	if strings.Contains(videoPath, "..") {
		return fmt.Errorf("path traversal not allowed in video path: %s", videoPath)
	}

	// Final security check: validate path is safe before opening
	if err := r.validatePath(videoPath, currentDir); err != nil {
		return fmt.Errorf("video path validation failed: %w", err)
	}

	file, err := os.Open(videoPath) //nolint:gosec // Path validated above
	if err != nil {
		return fmt.Errorf("failed to open video file: %w", err)
	}
	defer file.Close()

	part, err := writer.CreateFormFile("video", filepath.Base(videoPath))
	if err != nil {
		return fmt.Errorf("failed to create form file: %w", err)
	}

	if _, err := io.Copy(part, file); err != nil {
		return fmt.Errorf("failed to copy video file: %w", err)
	}

	// Add metadata
	if err := writer.WriteField("title", title); err != nil {
		return fmt.Errorf("failed to write title field: %w", err)
	}
	if err := writer.WriteField("description", description); err != nil {
		return fmt.Errorf("failed to write description field: %w", err)
	}

	// Close writer
	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	// Create request
	req, err := http.NewRequest("POST", r.apiURL+"/api/videos/upload", body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check response
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("upload failed with status %d and failed to read response body: %w", resp.StatusCode, err)
		}
		return fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// GetVideos retrieves the list of available videos
func (r *Recorder) GetVideos() ([]map[string]interface{}, error) {
	resp, err := http.Get(r.apiURL + "/api/videos")
	if err != nil {
		return nil, fmt.Errorf("failed to get videos: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to get videos with status %d and failed to read response body: %w", resp.StatusCode, err)
		}
		return nil, fmt.Errorf("failed to get videos with status %d: %s", resp.StatusCode, string(body))
	}

	var videos []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&videos); err != nil {
		return nil, fmt.Errorf("failed to decode videos: %w", err)
	}

	return videos, nil
}

// GetVideoURL returns the URL for a video
func (r *Recorder) GetVideoURL(id string) string {
	return fmt.Sprintf("%s/api/videos/%s", r.apiURL, id)
}

// GetThumbnailURL returns the URL for a video thumbnail
func (r *Recorder) GetThumbnailURL(id string) string {
	return fmt.Sprintf("%s/api/videos/thumbnail/%s", r.apiURL, id)
}

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
)

// Recorder handles the video recording UI
type Recorder struct {
	apiURL string
}

// NewRecorder creates a new video recorder UI component
func NewRecorder(apiURL string) *Recorder {
	return &Recorder{
		apiURL: apiURL,
	}
}

// UploadVideo uploads a recorded video to the server
func (r *Recorder) UploadVideo(videoPath, title, description string) error {
	// Create multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add video file
	file, err := os.Open(videoPath)
	if err != nil {
		return fmt.Errorf("failed to open video file: %v", err)
	}
	defer file.Close()

	part, err := writer.CreateFormFile("video", filepath.Base(videoPath))
	if err != nil {
		return fmt.Errorf("failed to create form file: %v", err)
	}

	if _, err := io.Copy(part, file); err != nil {
		return fmt.Errorf("failed to copy video file: %v", err)
	}

	// Add metadata
	if err := writer.WriteField("title", title); err != nil {
		return fmt.Errorf("failed to write title field: %v", err)
	}
	if err := writer.WriteField("description", description); err != nil {
		return fmt.Errorf("failed to write description field: %v", err)
	}

	// Close writer
	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %v", err)
	}

	// Create request
	req, err := http.NewRequest("POST", r.apiURL+"/api/videos/upload", body)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check response
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// GetVideos retrieves the list of available videos
func (r *Recorder) GetVideos() ([]map[string]interface{}, error) {
	resp, err := http.Get(r.apiURL + "/api/videos")
	if err != nil {
		return nil, fmt.Errorf("failed to get videos: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get videos with status %d: %s", resp.StatusCode, string(body))
	}

	var videos []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&videos); err != nil {
		return nil, fmt.Errorf("failed to decode videos: %v", err)
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

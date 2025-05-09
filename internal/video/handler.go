package video

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"gooji/pkg/ffmpeg"
)

// Handler manages video recording and processing
type Handler struct {
	processor *ffmpeg.Processor
	videoDir  string
}

// VideoMetadata represents metadata for a recorded video
type VideoMetadata struct {
	ID          string    `json:"id"`
	Filename    string    `json:"filename"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Duration    int       `json:"duration"`
	CreatedAt   time.Time `json:"created_at"`
	Tags        []string  `json:"tags"`
}

// NewHandler creates a new video handler
func NewHandler(processor *ffmpeg.Processor, videoDir string) (*Handler, error) {
	if err := os.MkdirAll(videoDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create video directory: %v", err)
	}

	return &Handler{
		processor: processor,
		videoDir:  videoDir,
	}, nil
}

// HandleUpload processes an uploaded video file
func (h *Handler) HandleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get video file
	file, header, err := r.FormFile("video")
	if err != nil {
		http.Error(w, "Failed to get video file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create unique filename
	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), header.Filename)
	filepath := filepath.Join(h.videoDir, filename)

	// Save file
	dst, err := os.Create(filepath)
	if err != nil {
		http.Error(w, "Failed to save video", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Failed to save video", http.StatusInternalServerError)
		return
	}

	// Get video metadata
	_, err = h.processor.GetVideoInfo(filepath)
	if err != nil {
		http.Error(w, "Failed to get video info", http.StatusInternalServerError)
		return
	}

	// Create metadata
	metadata := VideoMetadata{
		ID:          filename,
		Filename:    filename,
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		CreatedAt:   time.Now(),
		Tags:        []string{"ojibwe", "language", "culture"},
	}

	// Save metadata
	metadataPath := filepath.Join(h.videoDir, filename+".json")
	metadataFile, err := os.Create(metadataPath)
	if err != nil {
		http.Error(w, "Failed to save metadata", http.StatusInternalServerError)
		return
	}
	defer metadataFile.Close()

	if err := json.NewEncoder(metadataFile).Encode(metadata); err != nil {
		http.Error(w, "Failed to save metadata", http.StatusInternalServerError)
		return
	}

	// Generate thumbnail
	thumbnailPath := filepath.Join(h.videoDir, filename+".jpg")
	if err := h.processor.GenerateThumbnail(filepath, thumbnailPath, 1); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Failed to generate thumbnail: %v\n", err)
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"id":       metadata.ID,
		"filename": metadata.Filename,
	})
}

// GetVideo returns a video file
func (h *Handler) GetVideo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing video ID", http.StatusBadRequest)
		return
	}

	filepath := filepath.Join(h.videoDir, id)
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		http.Error(w, "Video not found", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, filepath)
}

// GetThumbnail returns a video thumbnail
func (h *Handler) GetThumbnail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing video ID", http.StatusBadRequest)
		return
	}

	thumbnailPath := filepath.Join(h.videoDir, id+".jpg")
	if _, err := os.Stat(thumbnailPath); os.IsNotExist(err) {
		http.Error(w, "Thumbnail not found", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, thumbnailPath)
}

// ListVideos returns a list of available videos
func (h *Handler) ListVideos(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	files, err := os.ReadDir(h.videoDir)
	if err != nil {
		http.Error(w, "Failed to list videos", http.StatusInternalServerError)
		return
	}

	var videos []VideoMetadata
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			metadataPath := filepath.Join(h.videoDir, file.Name())
			metadataFile, err := os.Open(metadataPath)
			if err != nil {
				continue
			}

			var metadata VideoMetadata
			if err := json.NewDecoder(metadataFile).Decode(&metadata); err != nil {
				metadataFile.Close()
				continue
			}
			metadataFile.Close()

			videos = append(videos, metadata)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(videos)
}

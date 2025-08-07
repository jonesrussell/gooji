package video

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"gooji/internal/config"
	"gooji/pkg/ffmpeg"
)

// Handler manages video recording and processing
type Handler struct {
	processor *ffmpeg.Processor
	storage   config.Storage
	templates *template.Template
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

// NewHandler creates a new video handler
func NewHandler(processor *ffmpeg.Processor, storage config.Storage) (*Handler, error) {
	// Create storage directories
	dirs := []string{storage.Uploads, storage.Temp, storage.Logs, storage.Thumbnails, storage.Metadata}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory %s: %v", dir, err)
		}
	}

	// Parse templates
	templates, err := template.ParseFiles(
		"web/templates/base.html",
		"web/templates/record.html",
		"web/templates/gallery.html",
		"web/templates/camera-test.html",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %v", err)
	}

	return &Handler{
		processor: processor,
		storage:   storage,
		templates: templates,
	}, nil
}

// HandleHome serves the home page
func (h *Handler) HandleHome(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := h.templates.ExecuteTemplate(w, "base.html", map[string]interface{}{
		"Page":         "home",
		"IsRecordPage": false,
	}); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// HandleRecord serves the recording page
func (h *Handler) HandleRecord(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := h.templates.ExecuteTemplate(w, "base.html", map[string]interface{}{
		"Page":         "record",
		"IsRecordPage": true,
	}); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// HandleEdit serves the video editing page
func (h *Handler) HandleEdit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := h.templates.ExecuteTemplate(w, "base.html", map[string]interface{}{
		"Page":         "edit",
		"IsRecordPage": false,
	}); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// HandleGallery serves the video gallery page
func (h *Handler) HandleGallery(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := h.templates.ExecuteTemplate(w, "base.html", map[string]interface{}{
		"Page":         "gallery",
		"IsRecordPage": false,
	}); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// HandleCameraTest serves the camera test page
func (h *Handler) HandleCameraTest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := h.templates.ExecuteTemplate(w, "base.html", map[string]interface{}{
		"Page":         "camera-test",
		"IsRecordPage": false,
	}); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// HandleHealth provides system health information
func (h *Handler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Test FFmpeg availability
	ffmpegWorking := false
	if _, err := exec.LookPath("ffmpeg"); err == nil {
		ffmpegWorking = true
	}

	// Check video directory
	videoDirAccessible := false
	if _, err := os.Stat(h.storage.Uploads); err == nil {
		videoDirAccessible = true
	}

	health := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now(),
		"checks": map[string]bool{
			"ffmpeg":    ffmpegWorking,
			"video_dir": videoDirAccessible,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(health)
}

// HandleVideos handles video-related API endpoints
func (h *Handler) HandleVideos(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.ListVideos(w, r)
	case http.MethodPost:
		h.HandleUpload(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HandleVideo handles individual video API endpoints
func (h *Handler) HandleVideo(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetVideo(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HandleUpload processes an uploaded video file
func (h *Handler) HandleUpload(w http.ResponseWriter, r *http.Request) {
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
	videoPath := filepath.Join(h.storage.Uploads, filename)

	// Save file
	dst, err := os.Create(videoPath)
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
	info, err := h.processor.GetVideoInfo(videoPath)
	if err != nil {
		// Log the specific error for debugging
		fmt.Printf("FFmpeg error for file %s: %v\n", videoPath, err)
		http.Error(w, fmt.Sprintf("Failed to process video: %v", err), http.StatusInternalServerError)
		return
	}

	// Create metadata
	metadata := VideoMetadata{
		ID:          filename,
		Filename:    filename,
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		Duration:    info.Duration,
		CreatedAt:   time.Now(),
		Tags:        []string{"ojibwe", "language", "culture"},
	}

	// Save metadata
	metadataPath := filepath.Join(h.storage.Metadata, filename+".json")
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
	thumbnailPath := filepath.Join(h.storage.Thumbnails, filename+".jpg")
	if err := h.processor.GenerateThumbnail(videoPath, thumbnailPath, 1); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Failed to generate thumbnail for %s: %v\n", videoPath, err)
	} else {
		fmt.Printf("Successfully generated thumbnail: %s\n", thumbnailPath)
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
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing video ID", http.StatusBadRequest)
		return
	}

	videoPath := filepath.Join(h.storage.Uploads, id)
	if _, err := os.Stat(videoPath); os.IsNotExist(err) {
		http.Error(w, "Video not found", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, videoPath)
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

	thumbnailPath := filepath.Join(h.storage.Thumbnails, id+".jpg")
	if _, err := os.Stat(thumbnailPath); os.IsNotExist(err) {
		http.Error(w, "Thumbnail not found", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, thumbnailPath)
}

// ListVideos returns a list of available videos
func (h *Handler) ListVideos(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir(h.storage.Metadata)
	if err != nil {
		http.Error(w, "Failed to list videos", http.StatusInternalServerError)
		return
	}

	var videos []VideoMetadata
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			metadataPath := filepath.Join(h.storage.Metadata, file.Name())
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

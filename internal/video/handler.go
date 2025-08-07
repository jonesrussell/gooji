package video

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"gooji/internal/config"
	"gooji/internal/logger"
	"gooji/pkg/ffmpeg"
)

// Handler manages video recording and processing
type Handler struct {
	service   Service
	templates map[string]*template.Template
	logger    *logger.Logger
	storage   config.Storage
}

// NewHandler creates a new video handler
func NewHandler(processor *ffmpeg.Processor, storage config.Storage, log *logger.Logger) (*Handler, error) {
	// Create storage directories
	if err := createStorageDirectories(storage); err != nil {
		return nil, fmt.Errorf("failed to create storage directories: %w", err)
	}

	// Create secure processor with allowed directory restriction
	secureProcessor := ffmpeg.NewProcessorWithSecurity(processor.FFmpegPath(), storage.Uploads)

	// Create repository and service
	repo := NewRepository(storage, log)
	service := NewService(repo, secureProcessor, log)

	// Parse templates
	templates, err := parseTemplates()
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}

	// Log template parsing success
	log.Info("Templates parsed successfully")
	for name := range templates {
		log.Debug("Available template: %s", name)
	}

	return &Handler{
		service:   service,
		templates: templates,
		logger:    log,
		storage:   storage,
	}, nil
}

// HandleHome serves the home page
func (h *Handler) HandleHome(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.handleMethodNotAllowed(w, r)
		return
	}

	h.logger.Debug("Serving home page")
	h.logger.Debug("Available templates: %v", len(h.templates))

	if err := h.templates["home"].ExecuteTemplate(w, "base.html", map[string]interface{}{
		"Page":         "home",
		"IsRecordPage": false,
	}); err != nil {
		h.logger.Error("Template execution error: %v", err)
		h.handleInternalError(w, r, err)
		return
	}
}

// HandleRecord serves the recording page
func (h *Handler) HandleRecord(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.handleMethodNotAllowed(w, r)
		return
	}

	h.logger.Debug("Serving record page")

	if err := h.templates["record"].ExecuteTemplate(w, "base.html", map[string]interface{}{
		"Page":         "record",
		"IsRecordPage": true,
	}); err != nil {
		h.handleInternalError(w, r, err)
		return
	}
}

// HandleEdit serves the video editing page
func (h *Handler) HandleEdit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.handleMethodNotAllowed(w, r)
		return
	}

	if err := h.templates["editor"].ExecuteTemplate(w, "editor.html", map[string]interface{}{
		"Page":         "edit",
		"IsRecordPage": false,
	}); err != nil {
		h.handleInternalError(w, r, err)
		return
	}
}

// HandleGallery serves the video gallery page
func (h *Handler) HandleGallery(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.handleMethodNotAllowed(w, r)
		return
	}

	if err := h.templates["gallery"].ExecuteTemplate(w, "gallery.html", map[string]interface{}{
		"Page":         "gallery",
		"IsRecordPage": false,
	}); err != nil {
		h.handleInternalError(w, r, err)
		return
	}
}

// HandleCameraTest serves the camera test page
func (h *Handler) HandleCameraTest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.handleMethodNotAllowed(w, r)
		return
	}

	h.logger.Debug("Serving camera test page")

	if err := h.templates["camera-test"].ExecuteTemplate(w, "base.html", map[string]interface{}{
		"Page":         "camera-test",
		"IsRecordPage": false,
	}); err != nil {
		h.handleInternalError(w, r, err)
		return
	}
}

// HandleHealth provides system health information
func (h *Handler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.handleMethodNotAllowed(w, r)
		return
	}

	health := h.getHealthStatus()
	h.writeJSONResponse(w, health)
}

// HandleVideos handles video-related API endpoints
func (h *Handler) HandleVideos(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.ListVideos(w, r)
	case http.MethodPost:
		h.HandleUpload(w, r)
	default:
		h.handleMethodNotAllowed(w, r)
	}
}

// HandleVideo handles individual video API endpoints
func (h *Handler) HandleVideo(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetVideo(w, r)
	default:
		h.handleMethodNotAllowed(w, r)
	}
}

// HandleUpload processes an uploaded video file
func (h *Handler) HandleUpload(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		h.handleValidationError(w, r, "Failed to parse form", err)
		return
	}

	// Get video file
	file, header, err := r.FormFile("video")
	if err != nil {
		h.handleValidationError(w, r, "Failed to get video file", err)
		return
	}
	defer file.Close()

	// Create upload metadata
	metadata := &UploadMetadata{
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		Tags:        []string{"ojibwe", "language", "culture"}, // Default tags
	}

	// Process upload through service
	videoMetadata, err := h.service.ProcessUpload(r.Context(), file, header, metadata)
	if err != nil {
		h.handleServiceError(w, r, err)
		return
	}

	// Return success response
	response := map[string]string{
		"id":       videoMetadata.ID,
		"filename": videoMetadata.Filename,
	}
	h.writeJSONResponse(w, response)
}

// GetVideo returns a video file
func (h *Handler) GetVideo(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		h.handleValidationError(w, r, "Missing video ID", nil)
		return
	}

	videoPath := filepath.Join(h.storage.Uploads, id)
	if _, err := os.Stat(videoPath); os.IsNotExist(err) {
		h.handleNotFoundError(w, r, "Video not found", err)
		return
	}

	http.ServeFile(w, r, videoPath)
}

// GetThumbnail returns a video thumbnail
func (h *Handler) GetThumbnail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.handleMethodNotAllowed(w, r)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		h.handleValidationError(w, r, "Missing video ID", nil)
		return
	}

	thumbnailPath := filepath.Join(h.storage.Thumbnails, id+".jpg")
	if _, err := os.Stat(thumbnailPath); os.IsNotExist(err) {
		h.handleNotFoundError(w, r, "Thumbnail not found", err)
		return
	}

	http.ServeFile(w, r, thumbnailPath)
}

// ListVideos returns a list of available videos
func (h *Handler) ListVideos(w http.ResponseWriter, r *http.Request) {
	videos, err := h.service.ListVideos(r.Context())
	if err != nil {
		h.handleServiceError(w, r, err)
		return
	}

	h.writeJSONResponse(w, videos)
}

// Helper functions

// createStorageDirectories creates all required storage directories
func createStorageDirectories(storage config.Storage) error {
	dirs := []string{storage.Uploads, storage.Temp, storage.Logs, storage.Thumbnails, storage.Metadata}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0o750); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}
	return nil
}

// parseTemplates parses all HTML templates
func parseTemplates() (map[string]*template.Template, error) {
	// Parse base template first
	baseTemplate, err := template.ParseFiles("web/templates/base.html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse base template: %w", err)
	}

	// Parse individual page templates that use base template
	homeTemplate, err := template.Must(baseTemplate.Clone()).ParseFiles("web/templates/home.html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse home template: %w", err)
	}

	recordTemplate, err := template.Must(baseTemplate.Clone()).ParseFiles("web/templates/record.html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse record template: %w", err)
	}

	cameraTestTemplate, err := template.Must(baseTemplate.Clone()).ParseFiles("web/templates/camera-test.html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse camera test template: %w", err)
	}

	// Parse standalone templates
	galleryTemplate, err := template.ParseFiles("web/templates/gallery.html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse gallery template: %w", err)
	}

	editorTemplate, err := template.ParseFiles("web/templates/editor.html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse editor template: %w", err)
	}

	// Create a template map for easy access
	templates := map[string]*template.Template{
		"home":        homeTemplate,
		"record":      recordTemplate,
		"gallery":     galleryTemplate,
		"editor":      editorTemplate,
		"camera-test": cameraTestTemplate,
	}

	return templates, nil
}

// getHealthStatus returns the current health status
func (h *Handler) getHealthStatus() map[string]interface{} {
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

	return map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now(),
		"checks": map[string]bool{
			"ffmpeg":    ffmpegWorking,
			"video_dir": videoDirAccessible,
		},
	}
}

// Error handling helpers

// handleMethodNotAllowed handles HTTP method not allowed errors
func (h *Handler) handleMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

// handleValidationError handles validation errors
func (h *Handler) handleValidationError(w http.ResponseWriter, r *http.Request, message string, err error) {
	h.logger.Error("Validation error: %s - %v (method: %s, path: %s, remote: %s)",
		message, err, r.Method, r.URL.Path, r.RemoteAddr)
	http.Error(w, message, http.StatusBadRequest)
}

// handleNotFoundError handles not found errors
func (h *Handler) handleNotFoundError(w http.ResponseWriter, r *http.Request, message string, err error) {
	h.logger.Error("Not found error: %s - %v (method: %s, path: %s, remote: %s)",
		message, err, r.Method, r.URL.Path, r.RemoteAddr)
	http.Error(w, message, http.StatusNotFound)
}

// handleInternalError handles internal server errors
func (h *Handler) handleInternalError(w http.ResponseWriter, r *http.Request, err error) {
	h.logger.Error("Internal server error: %v (method: %s, path: %s, remote: %s)",
		err, r.Method, r.URL.Path, r.RemoteAddr)
	http.Error(w, "Internal server error", http.StatusInternalServerError)
}

// handleServiceError handles service layer errors
func (h *Handler) handleServiceError(w http.ResponseWriter, r *http.Request, err error) {
	h.logger.Error("Service error: %v (method: %s, path: %s, remote: %s)",
		err, r.Method, r.URL.Path, r.RemoteAddr)

	// Use structured error handling if available
	statusCode := GetHTTPStatusCode(err)
	http.Error(w, "Service error", statusCode)
}

// writeJSONResponse writes a JSON response
func (h *Handler) writeJSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("Failed to encode JSON response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

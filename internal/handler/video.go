package handler

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/romsar/gonertia"
	"go.uber.org/zap"

	"gooji/internal/domain"
	"gooji/internal/service"
)

type VideoHandler struct {
	inertia *gonertia.Inertia
	service *service.VideoService
	logger  *zap.Logger
	storage string
}

func NewVideoHandler(i *gonertia.Inertia, s *service.VideoService, l *zap.Logger, storagePath string) *VideoHandler {
	return &VideoHandler{
		inertia: i,
		service: s,
		logger:  l,
		storage: storagePath,
	}
}

func (h *VideoHandler) Upload(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(100 << 20); err != nil {
		h.jsonError(w, "failed to parse form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("video")
	if err != nil {
		h.jsonError(w, "missing video file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	req := service.UploadRequest{
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
	}

	video, err := h.service.Upload(r.Context(), file, header, req)
	if err != nil {
		h.logger.Error("upload failed", zap.Error(err))
		h.jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.jsonResponse(w, video, http.StatusCreated)
}

func (h *VideoHandler) List(w http.ResponseWriter, r *http.Request) {
	videos, err := h.service.List(r.Context(), domain.VideoListOptions{Limit: 50})
	if err != nil {
		h.logger.Error("list failed", zap.Error(err))
		h.jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.jsonResponse(w, videos, http.StatusOK)
}

func (h *VideoHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := h.extractID(r.URL.Path, "/api/videos/")
	if id == "" {
		h.jsonError(w, "missing video ID", http.StatusBadRequest)
		return
	}

	video, err := h.service.Get(r.Context(), id)
	if err != nil {
		if err == domain.ErrVideoNotFound {
			h.jsonError(w, "video not found", http.StatusNotFound)
			return
		}
		h.jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.ServeFile(w, r, video.FilePath)
}

func (h *VideoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := h.extractID(r.URL.Path, "/api/videos/")
	if id == "" {
		h.jsonError(w, "missing video ID", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		if err == domain.ErrVideoNotFound {
			h.jsonError(w, "video not found", http.StatusNotFound)
			return
		}
		h.jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.jsonResponse(w, map[string]string{"status": "deleted"}, http.StatusOK)
}

func (h *VideoHandler) Thumbnail(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	// Remove extension if present and add .jpg
	base := strings.TrimSuffix(id, filepath.Ext(id))
	thumbnailPath := filepath.Join(h.storage, "thumbnails", base+".jpg")

	http.ServeFile(w, r, thumbnailPath)
}

func (h *VideoHandler) extractID(path, prefix string) string {
	return strings.TrimPrefix(path, prefix)
}

func (h *VideoHandler) jsonResponse(w http.ResponseWriter, data any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *VideoHandler) jsonError(w http.ResponseWriter, msg string, status int) {
	h.jsonResponse(w, map[string]string{"error": msg}, status)
}

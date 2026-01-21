package handler

import (
	"net/http"
)

type VideoDTO struct {
	ID              string `json:"id"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	ThumbnailPath   string `json:"thumbnailPath"`
	DurationSeconds int    `json:"durationSeconds"`
}

func (h *Handler) Gallery(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	videos, err := h.videoService.ListPublished(ctx)
	if err != nil {
		http.Error(w, "Failed to load videos", http.StatusInternalServerError)
		return
	}

	// Convert to DTOs
	dtos := make([]VideoDTO, 0, len(videos))
	for _, v := range videos {
		dtos = append(dtos, VideoDTO{
			ID:              v.ID,
			Title:           v.Title,
			Description:     v.Description,
			ThumbnailPath:   v.ThumbnailPath,
			DurationSeconds: v.DurationSeconds,
		})
	}

	err = h.inertia.Render(w, r, "Gallery", map[string]any{
		"videos": dtos,
	})
	if err != nil {
		http.Error(w, "Failed to render", http.StatusInternalServerError)
	}
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	err := h.inertia.Render(w, r, "Home", map[string]any{})
	if err != nil {
		http.Error(w, "Failed to render", http.StatusInternalServerError)
	}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

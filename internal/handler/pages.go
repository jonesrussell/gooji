package handler

import (
	"net/http"

	"github.com/romsar/gonertia"
	"go.uber.org/zap"

	"gooji/internal/domain"
	"gooji/internal/service"
)

type PageHandler struct {
	inertia *gonertia.Inertia
	service *service.VideoService
	logger  *zap.Logger
}

func NewPageHandler(i *gonertia.Inertia, s *service.VideoService, l *zap.Logger) *PageHandler {
	return &PageHandler{
		inertia: i,
		service: s,
		logger:  l,
	}
}

func (h *PageHandler) Home(w http.ResponseWriter, r *http.Request) {
	h.inertia.Render(w, r, "Home", gonertia.Props{
		"message": "Welcome to Gooji",
	})
}

func (h *PageHandler) Gallery(w http.ResponseWriter, r *http.Request) {
	videos, err := h.service.List(r.Context(), domain.VideoListOptions{Limit: 50})
	if err != nil {
		h.logger.Error("failed to list videos", zap.Error(err))
		videos = []domain.Video{}
	}

	h.inertia.Render(w, r, "Gallery", gonertia.Props{
		"videos": videos,
	})
}

func (h *PageHandler) Upload(w http.ResponseWriter, r *http.Request) {
	h.inertia.Render(w, r, "Upload", gonertia.Props{})
}

func (h *PageHandler) Record(w http.ResponseWriter, r *http.Request) {
	h.inertia.Render(w, r, "Record", gonertia.Props{})
}

func (h *PageHandler) VideoShow(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/video/"):]
	video, err := h.service.Get(r.Context(), id)
	if err != nil {
		h.inertia.Render(w, r, "Error", gonertia.Props{
			"message": "Video not found",
		})
		return
	}

	h.inertia.Render(w, r, "Video/Show", gonertia.Props{
		"video": video,
	})
}

func (h *PageHandler) Edit(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/edit/"):]
	video, err := h.service.Get(r.Context(), id)
	if err != nil {
		h.inertia.Render(w, r, "Error", gonertia.Props{
			"message": "Video not found",
		})
		return
	}

	h.inertia.Render(w, r, "Edit", gonertia.Props{
		"video": video,
	})
}

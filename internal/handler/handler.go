package handler

import (
	"net/http"

	inertia "github.com/petaki/inertia-go"

	"gooji/internal/service"
)

type Handler struct {
	inertia      *inertia.Inertia
	videoService *service.VideoService
}

func New(videoService *service.VideoService, version string) *Handler {
	i := inertia.New("http://localhost:5173", "app.html", version)

	return &Handler{
		inertia:      i,
		videoService: videoService,
	}
}

func (h *Handler) ShareFunc(r *http.Request) map[string]any {
	return map[string]any{
		"appName": "Gooji",
	}
}

func (h *Handler) Inertia() *inertia.Inertia {
	return h.inertia
}

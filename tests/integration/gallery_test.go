package integration_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"gooji/internal/domain"
	"gooji/internal/handler"
	"gooji/internal/repository/sqlite"
	"gooji/internal/service"
)

func setupTestServer(t *testing.T) (*httptest.Server, *service.VideoService) {
	db, err := sqlite.NewDB(":memory:")
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	videoRepo := sqlite.NewVideoRepository(db)
	videoService := service.NewVideoService(videoRepo)
	h := handler.New(videoService, "")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", h.Health)

	return httptest.NewServer(mux), videoService
}

func TestHealthEndpoint(t *testing.T) {
	server, _ := setupTestServer(t)
	defer server.Close()

	resp, err := http.Get(server.URL + "/health")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var result map[string]string
	json.Unmarshal(body, &result)

	if result["status"] != "ok" {
		t.Errorf("expected status ok, got %s", result["status"])
	}
}

func TestVideoServiceIntegration(t *testing.T) {
	_, videoService := setupTestServer(t)
	ctx := context.Background()

	// Create a video
	video := domain.NewVideo("Integration Test", "/test.mp4", "video/mp4")
	err := videoService.Create(ctx, video)
	if err != nil {
		t.Fatalf("failed to create: %v", err)
	}

	// List pending (should have 1)
	pending, _ := videoService.ListPending(ctx)
	if len(pending) != 1 {
		t.Errorf("expected 1 pending, got %d", len(pending))
	}

	// Approve it
	err = videoService.Approve(ctx, video.ID)
	if err != nil {
		t.Fatalf("failed to approve: %v", err)
	}

	// List published (should have 1)
	published, _ := videoService.ListPublished(ctx)
	if len(published) != 1 {
		t.Errorf("expected 1 published, got %d", len(published))
	}

	// List pending (should have 0)
	pending, _ = videoService.ListPending(ctx)
	if len(pending) != 0 {
		t.Errorf("expected 0 pending, got %d", len(pending))
	}
}

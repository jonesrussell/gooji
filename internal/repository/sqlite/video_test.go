package sqlite_test

import (
	"context"
	"testing"

	"gooji/internal/domain"
	"gooji/internal/repository/sqlite"
)

func setupTestDB(t *testing.T) *sqlite.DB {
	db, err := sqlite.NewDB(":memory:")
	if err != nil {
		t.Fatalf("failed to create test db: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func TestVideoRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := sqlite.NewVideoRepository(db)
	ctx := context.Background()

	video := domain.NewVideo("Test Video", "/videos/test.mp4", "video/mp4")
	video.Description = "A test video"

	err := repo.Create(ctx, video)
	if err != nil {
		t.Fatalf("failed to create video: %v", err)
	}

	// Verify it was saved
	saved, err := repo.GetByID(ctx, video.ID)
	if err != nil {
		t.Fatalf("failed to get video: %v", err)
	}
	if saved.Title != "Test Video" {
		t.Errorf("expected title 'Test Video', got '%s'", saved.Title)
	}
	if saved.Description != "A test video" {
		t.Errorf("expected description 'A test video', got '%s'", saved.Description)
	}
}

func TestVideoRepository_List(t *testing.T) {
	db := setupTestDB(t)
	repo := sqlite.NewVideoRepository(db)
	ctx := context.Background()

	// Create videos with different statuses
	v1 := domain.NewVideo("Video 1", "/v1.mp4", "video/mp4")
	v2 := domain.NewVideo("Video 2", "/v2.mp4", "video/mp4")
	v2.Publish()
	v3 := domain.NewVideo("Video 3", "/v3.mp4", "video/mp4")

	repo.Create(ctx, v1)
	repo.Create(ctx, v2)
	repo.Create(ctx, v3)

	// List all
	all, err := repo.List(ctx, nil)
	if err != nil {
		t.Fatalf("failed to list all: %v", err)
	}
	if len(all) != 3 {
		t.Errorf("expected 3 videos, got %d", len(all))
	}

	// List pending only
	pending := domain.StatusPending
	pendingVideos, err := repo.List(ctx, &pending)
	if err != nil {
		t.Fatalf("failed to list pending: %v", err)
	}
	if len(pendingVideos) != 2 {
		t.Errorf("expected 2 pending videos, got %d", len(pendingVideos))
	}
}

func TestVideoRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := sqlite.NewVideoRepository(db)
	ctx := context.Background()

	video := domain.NewVideo("Original", "/test.mp4", "video/mp4")
	repo.Create(ctx, video)

	video.Title = "Updated"
	video.Publish()
	err := repo.Update(ctx, video)
	if err != nil {
		t.Fatalf("failed to update: %v", err)
	}

	updated, _ := repo.GetByID(ctx, video.ID)
	if updated.Title != "Updated" {
		t.Errorf("expected 'Updated', got '%s'", updated.Title)
	}
	if updated.Status != domain.StatusPublished {
		t.Errorf("expected Published, got %s", updated.Status)
	}
}

func TestVideoRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := sqlite.NewVideoRepository(db)
	ctx := context.Background()

	video := domain.NewVideo("To Delete", "/test.mp4", "video/mp4")
	repo.Create(ctx, video)

	err := repo.Delete(ctx, video.ID)
	if err != nil {
		t.Fatalf("failed to delete: %v", err)
	}

	_, err = repo.GetByID(ctx, video.ID)
	if err != domain.ErrVideoNotFound {
		t.Errorf("expected ErrVideoNotFound, got %v", err)
	}
}

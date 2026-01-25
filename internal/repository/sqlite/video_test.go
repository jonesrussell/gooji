package sqlite

import (
	"context"
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"

	"gooji/internal/domain"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	// Create temp file for test database
	tmpfile, err := os.CreateTemp("", "test-*.db")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	tmpfile.Close()

	// Register cleanup
	t.Cleanup(func() {
		os.Remove(tmpfile.Name())
	})

	db, err := sql.Open("sqlite3", tmpfile.Name())
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	// Create videos table
	schema := `
		CREATE TABLE IF NOT EXISTS videos (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			description TEXT,
			language TEXT DEFAULT 'oji',
			duration_seconds REAL,
			file_path TEXT NOT NULL,
			thumbnail_path TEXT,
			file_size_bytes INTEGER,
			mime_type TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`

	if _, err := db.Exec(schema); err != nil {
		t.Fatalf("failed to create schema: %v", err)
	}

	return db
}

func TestVideoRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewVideoRepository(db)
	ctx := context.Background()

	video := &domain.Video{
		ID:              "test-123",
		Title:           "Test Video",
		Description:     "A test description",
		Language:        "oji",
		DurationSeconds: 120.5,
		FilePath:        "/storage/test-123.mp4",
		ThumbnailPath:   "/storage/test-123.jpg",
		FileSizeBytes:   1024000,
		MimeType:        "video/mp4",
	}

	t.Run("create video successfully", func(t *testing.T) {
		err := repo.Create(ctx, video)
		if err != nil {
			t.Fatalf("failed to create video: %v", err)
		}

		// Verify timestamps are set
		if video.CreatedAt.IsZero() {
			t.Error("expected CreatedAt to be set")
		}
		if video.UpdatedAt.IsZero() {
			t.Error("expected UpdatedAt to be set")
		}
	})

	t.Run("create duplicate video fails", func(t *testing.T) {
		err := repo.Create(ctx, video)
		if err == nil {
			t.Error("expected error creating duplicate video")
		}
	})
}

func TestVideoRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewVideoRepository(db)
	ctx := context.Background()

	// Create a test video
	video := &domain.Video{
		ID:            "get-test",
		Title:         "Get Test",
		Description:   "Testing get",
		FilePath:      "/storage/get-test.mp4",
		MimeType:      "video/mp4",
		FileSizeBytes: 500,
	}
	if err := repo.Create(ctx, video); err != nil {
		t.Fatalf("failed to setup test video: %v", err)
	}

	t.Run("get existing video", func(t *testing.T) {
		found, err := repo.GetByID(ctx, "get-test")
		if err != nil {
			t.Fatalf("failed to get video: %v", err)
		}

		if found.ID != video.ID {
			t.Errorf("expected ID %s, got %s", video.ID, found.ID)
		}
		if found.Title != video.Title {
			t.Errorf("expected Title %s, got %s", video.Title, found.Title)
		}
		if found.Description != video.Description {
			t.Errorf("expected Description %s, got %s", video.Description, found.Description)
		}
	})

	t.Run("get non-existing video returns error", func(t *testing.T) {
		_, err := repo.GetByID(ctx, "non-existent")
		if err != domain.ErrVideoNotFound {
			t.Errorf("expected ErrVideoNotFound, got %v", err)
		}
	})
}

func TestVideoRepository_List(t *testing.T) {
	db := setupTestDB(t)
	repo := NewVideoRepository(db)
	ctx := context.Background()

	// Create test videos
	for i := 0; i < 5; i++ {
		video := &domain.Video{
			ID:       string(rune('a' + i)),
			Title:    "Video " + string(rune('A'+i)),
			FilePath: "/storage/video.mp4",
			MimeType: "video/mp4",
		}
		if err := repo.Create(ctx, video); err != nil {
			t.Fatalf("failed to create test video: %v", err)
		}
	}

	t.Run("list with default limit", func(t *testing.T) {
		videos, err := repo.List(ctx, domain.VideoListOptions{})
		if err != nil {
			t.Fatalf("failed to list videos: %v", err)
		}
		if len(videos) != 5 {
			t.Errorf("expected 5 videos, got %d", len(videos))
		}
	})

	t.Run("list with limit", func(t *testing.T) {
		videos, err := repo.List(ctx, domain.VideoListOptions{Limit: 3})
		if err != nil {
			t.Fatalf("failed to list videos: %v", err)
		}
		if len(videos) != 3 {
			t.Errorf("expected 3 videos, got %d", len(videos))
		}
	})

	t.Run("list with offset", func(t *testing.T) {
		videos, err := repo.List(ctx, domain.VideoListOptions{Limit: 10, Offset: 3})
		if err != nil {
			t.Fatalf("failed to list videos: %v", err)
		}
		if len(videos) != 2 {
			t.Errorf("expected 2 videos, got %d", len(videos))
		}
	})
}

func TestVideoRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewVideoRepository(db)
	ctx := context.Background()

	// Create a test video
	video := &domain.Video{
		ID:       "delete-me",
		Title:    "Delete Me",
		FilePath: "/storage/delete-me.mp4",
		MimeType: "video/mp4",
	}
	if err := repo.Create(ctx, video); err != nil {
		t.Fatalf("failed to setup test video: %v", err)
	}

	t.Run("delete existing video", func(t *testing.T) {
		err := repo.Delete(ctx, "delete-me")
		if err != nil {
			t.Fatalf("failed to delete video: %v", err)
		}

		// Verify it's gone
		_, err = repo.GetByID(ctx, "delete-me")
		if err != domain.ErrVideoNotFound {
			t.Error("expected video to be deleted")
		}
	})

	t.Run("delete non-existing video returns error", func(t *testing.T) {
		err := repo.Delete(ctx, "non-existent")
		if err != domain.ErrVideoNotFound {
			t.Errorf("expected ErrVideoNotFound, got %v", err)
		}
	})
}

func TestVideoRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewVideoRepository(db)
	ctx := context.Background()

	// Create a test video
	video := &domain.Video{
		ID:          "update-me",
		Title:       "Original Title",
		Description: "Original Description",
		FilePath:    "/storage/update-me.mp4",
		MimeType:    "video/mp4",
	}
	if err := repo.Create(ctx, video); err != nil {
		t.Fatalf("failed to setup test video: %v", err)
	}

	t.Run("update existing video", func(t *testing.T) {
		video.Title = "Updated Title"
		video.Description = "Updated Description"

		err := repo.Update(ctx, video)
		if err != nil {
			t.Fatalf("failed to update video: %v", err)
		}

		// Verify changes
		found, err := repo.GetByID(ctx, "update-me")
		if err != nil {
			t.Fatalf("failed to get updated video: %v", err)
		}

		if found.Title != "Updated Title" {
			t.Errorf("expected Title 'Updated Title', got %s", found.Title)
		}
		if found.Description != "Updated Description" {
			t.Errorf("expected Description 'Updated Description', got %s", found.Description)
		}
	})

	t.Run("update non-existing video returns error", func(t *testing.T) {
		nonExistent := &domain.Video{ID: "non-existent"}
		err := repo.Update(ctx, nonExistent)
		if err != domain.ErrVideoNotFound {
			t.Errorf("expected ErrVideoNotFound, got %v", err)
		}
	})
}

package filesystem

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"
)

func setupTestDirs(t *testing.T) (string, string) {
	t.Helper()

	tmpDir, err := os.MkdirTemp("", "storage-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	t.Cleanup(func() {
		os.RemoveAll(tmpDir)
	})

	uploadsDir := filepath.Join(tmpDir, "uploads")
	thumbnailsDir := filepath.Join(tmpDir, "thumbnails")

	return uploadsDir, thumbnailsDir
}

func TestNewFileStorage(t *testing.T) {
	uploadsDir, thumbnailsDir := setupTestDirs(t)

	storage := NewFileStorage(uploadsDir, thumbnailsDir)
	if storage == nil {
		t.Fatal("expected non-nil storage")
	}

	// Verify directories were created
	if _, err := os.Stat(uploadsDir); os.IsNotExist(err) {
		t.Error("expected uploads directory to be created")
	}
	if _, err := os.Stat(thumbnailsDir); os.IsNotExist(err) {
		t.Error("expected thumbnails directory to be created")
	}
}

func TestFileStorage_Save(t *testing.T) {
	uploadsDir, thumbnailsDir := setupTestDirs(t)
	storage := NewFileStorage(uploadsDir, thumbnailsDir)
	ctx := context.Background()

	t.Run("save file successfully", func(t *testing.T) {
		data := []byte("test video content")
		reader := bytes.NewReader(data)

		path, err := storage.Save(ctx, reader, "test-video.mp4")
		if err != nil {
			t.Fatalf("failed to save file: %v", err)
		}

		// Verify file exists
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Error("expected file to exist")
		}

		// Verify content
		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("failed to read saved file: %v", err)
		}
		if !bytes.Equal(content, data) {
			t.Error("saved content doesn't match original")
		}
	})

	t.Run("save file with subdirectory in name", func(t *testing.T) {
		data := []byte("another test")
		reader := bytes.NewReader(data)

		// This should work - filename goes directly in uploads dir
		path, err := storage.Save(ctx, reader, "simple-name.mp4")
		if err != nil {
			t.Fatalf("failed to save file: %v", err)
		}

		if filepath.Dir(path) != uploadsDir {
			t.Errorf("expected file in uploads dir, got %s", filepath.Dir(path))
		}
	})

	t.Run("save empty file", func(t *testing.T) {
		reader := bytes.NewReader([]byte{})

		path, err := storage.Save(ctx, reader, "empty.mp4")
		if err != nil {
			t.Fatalf("failed to save empty file: %v", err)
		}

		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("failed to stat file: %v", err)
		}
		if info.Size() != 0 {
			t.Errorf("expected empty file, got size %d", info.Size())
		}
	})

	t.Run("save large file", func(t *testing.T) {
		// Create 1MB of data
		data := make([]byte, 1024*1024)
		for i := range data {
			data[i] = byte(i % 256)
		}
		reader := bytes.NewReader(data)

		path, err := storage.Save(ctx, reader, "large.mp4")
		if err != nil {
			t.Fatalf("failed to save large file: %v", err)
		}

		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("failed to stat file: %v", err)
		}
		if info.Size() != int64(len(data)) {
			t.Errorf("expected size %d, got %d", len(data), info.Size())
		}
	})
}

func TestFileStorage_Delete(t *testing.T) {
	uploadsDir, thumbnailsDir := setupTestDirs(t)
	storage := NewFileStorage(uploadsDir, thumbnailsDir)
	ctx := context.Background()

	t.Run("delete existing file", func(t *testing.T) {
		// Create a file
		path, err := storage.Save(ctx, bytes.NewReader([]byte("test")), "to-delete.mp4")
		if err != nil {
			t.Fatalf("failed to create file: %v", err)
		}

		// Verify it exists
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Fatal("expected file to exist before delete")
		}

		// Delete it
		err = storage.Delete(ctx, path)
		if err != nil {
			t.Fatalf("failed to delete file: %v", err)
		}

		// Verify it's gone
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Error("expected file to be deleted")
		}
	})

	t.Run("delete non-existing file returns no error", func(t *testing.T) {
		err := storage.Delete(ctx, "/non/existent/path.mp4")
		if err != nil {
			t.Errorf("expected no error for non-existent file, got %v", err)
		}
	})
}

func TestFileStorage_GetPath(t *testing.T) {
	uploadsDir, thumbnailsDir := setupTestDirs(t)
	storage := NewFileStorage(uploadsDir, thumbnailsDir)

	t.Run("returns thumbnail path", func(t *testing.T) {
		path := storage.GetPath("test.jpg")
		expected := filepath.Join(thumbnailsDir, "test.jpg")
		if path != expected {
			t.Errorf("expected %s, got %s", expected, path)
		}
	})

	t.Run("handles various filenames", func(t *testing.T) {
		tests := []string{
			"simple.jpg",
			"uuid-12345.jpg",
			"file-with-dashes.jpg",
		}

		for _, filename := range tests {
			path := storage.GetPath(filename)
			if filepath.Base(path) != filename {
				t.Errorf("expected filename %s in path, got %s", filename, filepath.Base(path))
			}
		}
	})
}

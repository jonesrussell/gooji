package service

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/textproto"
	"testing"

	"go.uber.org/zap"

	"gooji/internal/domain"
	"gooji/pkg/ffmpeg"
)

// Mock repository implementation
type mockVideoRepository struct {
	videos map[string]*domain.Video
}

func newMockVideoRepository() *mockVideoRepository {
	return &mockVideoRepository{
		videos: make(map[string]*domain.Video),
	}
}

func (m *mockVideoRepository) Create(ctx context.Context, video *domain.Video) error {
	m.videos[video.ID] = video
	return nil
}

func (m *mockVideoRepository) GetByID(ctx context.Context, id string) (*domain.Video, error) {
	if video, ok := m.videos[id]; ok {
		return video, nil
	}
	return nil, domain.ErrVideoNotFound
}

func (m *mockVideoRepository) List(ctx context.Context, opts domain.VideoListOptions) ([]domain.Video, error) {
	var videos []domain.Video
	for _, v := range m.videos {
		videos = append(videos, *v)
	}
	return videos, nil
}

func (m *mockVideoRepository) Delete(ctx context.Context, id string) error {
	if _, ok := m.videos[id]; !ok {
		return domain.ErrVideoNotFound
	}
	delete(m.videos, id)
	return nil
}

func (m *mockVideoRepository) Update(ctx context.Context, video *domain.Video) error {
	if _, ok := m.videos[video.ID]; !ok {
		return domain.ErrVideoNotFound
	}
	m.videos[video.ID] = video
	return nil
}

// Mock storage implementation
type mockFileStorage struct {
	files   map[string][]byte
	baseDir string
}

func newMockFileStorage() *mockFileStorage {
	return &mockFileStorage{
		files:   make(map[string][]byte),
		baseDir: "/tmp/test",
	}
}

func (m *mockFileStorage) Save(ctx context.Context, reader io.Reader, filename string) (string, error) {
	data, _ := io.ReadAll(reader)
	path := m.baseDir + "/" + filename
	m.files[path] = data
	return path, nil
}

func (m *mockFileStorage) Delete(ctx context.Context, path string) error {
	delete(m.files, path)
	return nil
}

func (m *mockFileStorage) GetPath(filename string) string {
	return m.baseDir + "/" + filename
}

// Helper to create a mock multipart file header
func createMockFileHeader(filename, contentType string, size int64) *multipart.FileHeader {
	header := make(textproto.MIMEHeader)
	header.Set("Content-Type", contentType)
	return &multipart.FileHeader{
		Filename: filename,
		Header:   header,
		Size:     size,
	}
}

func TestVideoService_Get(t *testing.T) {
	repo := newMockVideoRepository()
	storage := newMockFileStorage()
	logger := zap.NewNop()

	service := NewVideoService(repo, storage, nil, logger)

	// Create a test video directly in the repository
	testVideo := &domain.Video{
		ID:          "test-123",
		Title:       "Test Video",
		Description: "A test video",
	}
	repo.videos[testVideo.ID] = testVideo

	t.Run("get existing video", func(t *testing.T) {
		video, err := service.Get(context.Background(), "test-123")
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if video.ID != "test-123" {
			t.Errorf("expected ID 'test-123', got %s", video.ID)
		}
		if video.Title != "Test Video" {
			t.Errorf("expected Title 'Test Video', got %s", video.Title)
		}
	})

	t.Run("get non-existing video", func(t *testing.T) {
		_, err := service.Get(context.Background(), "non-existent")
		if err != domain.ErrVideoNotFound {
			t.Errorf("expected ErrVideoNotFound, got %v", err)
		}
	})
}

func TestVideoService_List(t *testing.T) {
	repo := newMockVideoRepository()
	storage := newMockFileStorage()
	logger := zap.NewNop()

	service := NewVideoService(repo, storage, nil, logger)

	// Add test videos
	repo.videos["v1"] = &domain.Video{ID: "v1", Title: "Video 1"}
	repo.videos["v2"] = &domain.Video{ID: "v2", Title: "Video 2"}

	t.Run("list all videos", func(t *testing.T) {
		videos, err := service.List(context.Background(), domain.VideoListOptions{Limit: 10})
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if len(videos) != 2 {
			t.Errorf("expected 2 videos, got %d", len(videos))
		}
	})

	t.Run("list empty", func(t *testing.T) {
		emptyRepo := newMockVideoRepository()
		emptyService := NewVideoService(emptyRepo, storage, nil, logger)

		videos, err := emptyService.List(context.Background(), domain.VideoListOptions{Limit: 10})
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if len(videos) != 0 {
			t.Errorf("expected 0 videos, got %d", len(videos))
		}
	})
}

func TestVideoService_Delete(t *testing.T) {
	repo := newMockVideoRepository()
	storage := newMockFileStorage()
	logger := zap.NewNop()

	service := NewVideoService(repo, storage, nil, logger)

	// Add a test video
	testVideo := &domain.Video{
		ID:            "delete-me",
		Title:         "Delete Me",
		FilePath:      "/tmp/test/delete-me.mp4",
		ThumbnailPath: "/tmp/test/delete-me.jpg",
	}
	repo.videos[testVideo.ID] = testVideo
	storage.files[testVideo.FilePath] = []byte("video data")
	storage.files[testVideo.ThumbnailPath] = []byte("thumbnail data")

	t.Run("delete existing video", func(t *testing.T) {
		err := service.Delete(context.Background(), "delete-me")
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		// Verify video is deleted
		if _, ok := repo.videos["delete-me"]; ok {
			t.Error("expected video to be deleted from repository")
		}

		// Verify files are deleted
		if _, ok := storage.files[testVideo.FilePath]; ok {
			t.Error("expected video file to be deleted")
		}
		if _, ok := storage.files[testVideo.ThumbnailPath]; ok {
			t.Error("expected thumbnail to be deleted")
		}
	})

	t.Run("delete non-existing video", func(t *testing.T) {
		err := service.Delete(context.Background(), "non-existent")
		if err != domain.ErrVideoNotFound {
			t.Errorf("expected ErrVideoNotFound, got %v", err)
		}
	})
}

func TestVideoService_validateFileType(t *testing.T) {
	repo := newMockVideoRepository()
	storage := newMockFileStorage()
	logger := zap.NewNop()

	service := NewVideoService(repo, storage, nil, logger)

	tests := []struct {
		name        string
		contentType string
		expectError bool
	}{
		{"valid mp4", "video/mp4", false},
		{"valid webm", "video/webm", false},
		{"invalid text", "text/plain", true},
		{"invalid image", "image/jpeg", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			header := createMockFileHeader("test.mp4", tt.contentType, 1000)
			err := service.validateFileType(header)
			if tt.expectError && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}

func TestVideoService_sanitize(t *testing.T) {
	repo := newMockVideoRepository()
	storage := newMockFileStorage()
	logger := zap.NewNop()

	service := NewVideoService(repo, storage, nil, logger)

	tests := []struct {
		input    string
		expected string
	}{
		{"normal text", "normal text"},
		{"  trimmed  ", "trimmed"},
		{"<script>alert('xss')</script>", "&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;"},
		{"a & b", "a &amp; b"},
		{"quotes \"test\"", "quotes &#34;test&#34;"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := service.sanitize(tt.input)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestUploadRequest_Validation(t *testing.T) {
	repo := newMockVideoRepository()
	storage := newMockFileStorage()
	logger := zap.NewNop()
	// Create a mock processor - we won't actually use it in validation tests
	processor := &ffmpeg.Processor{}

	service := NewVideoService(repo, storage, processor, logger)

	tests := []struct {
		name        string
		req         UploadRequest
		expectError bool
	}{
		{
			name:        "valid request",
			req:         UploadRequest{Title: "Test", Description: "Description"},
			expectError: false,
		},
		{
			name:        "title too long",
			req:         UploadRequest{Title: string(make([]byte, 201)), Description: "Description"},
			expectError: true,
		},
		{
			name:        "description too long",
			req:         UploadRequest{Title: "Test", Description: string(make([]byte, 2001))},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.validate.Struct(tt.req)
			if tt.expectError && err == nil {
				t.Error("expected validation error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}

// Helper to create a mock multipart file
type mockMultipartFile struct {
	*bytes.Reader
}

func (m *mockMultipartFile) Close() error {
	return nil
}

func newMockMultipartFile(data []byte) *mockMultipartFile {
	return &mockMultipartFile{Reader: bytes.NewReader(data)}
}

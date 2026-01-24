package filesystem

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gooji/internal/repository"
)

type fileStorage struct {
	uploadsDir    string
	thumbnailsDir string
}

func NewFileStorage(uploadsDir, thumbnailsDir string) repository.FileStorage {
	os.MkdirAll(uploadsDir, 0o750)
	os.MkdirAll(thumbnailsDir, 0o750)
	return &fileStorage{
		uploadsDir:    uploadsDir,
		thumbnailsDir: thumbnailsDir,
	}
}

func (s *fileStorage) Save(ctx context.Context, reader io.Reader, filename string) (string, error) {
	path := filepath.Join(s.uploadsDir, filename)

	file, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, reader); err != nil {
		os.Remove(path)
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return path, nil
}

func (s *fileStorage) Delete(ctx context.Context, path string) error {
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

func (s *fileStorage) GetPath(filename string) string {
	return filepath.Join(s.thumbnailsDir, filename)
}

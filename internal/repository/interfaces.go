package repository

import (
	"context"
	"io"

	"gooji/internal/domain"
)

type VideoRepository interface {
	Create(ctx context.Context, video *domain.Video) error
	GetByID(ctx context.Context, id string) (*domain.Video, error)
	List(ctx context.Context, opts domain.VideoListOptions) ([]domain.Video, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, video *domain.Video) error
}

type FileStorage interface {
	Save(ctx context.Context, reader io.Reader, filename string) (string, error)
	Delete(ctx context.Context, path string) error
	GetPath(filename string) string
}

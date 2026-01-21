package repository

import (
	"context"

	"gooji/internal/domain"
)

type VideoRepository interface {
	Create(ctx context.Context, video *domain.Video) error
	GetByID(ctx context.Context, id string) (*domain.Video, error)
	List(ctx context.Context, status *domain.VideoStatus) ([]*domain.Video, error)
	Update(ctx context.Context, video *domain.Video) error
	Delete(ctx context.Context, id string) error
}

type DeviceRepository interface {
	Create(ctx context.Context, device *domain.Device) error
	GetByID(ctx context.Context, id string) (*domain.Device, error)
	GetByAPIKeyHash(ctx context.Context, hash string) (*domain.Device, error)
	List(ctx context.Context) ([]*domain.Device, error)
	Update(ctx context.Context, device *domain.Device) error
	Delete(ctx context.Context, id string) error
}

type TagRepository interface {
	Create(ctx context.Context, name string) (int64, error)
	GetByName(ctx context.Context, name string) (int64, error)
	GetOrCreate(ctx context.Context, name string) (int64, error)
	AddToVideo(ctx context.Context, videoID string, tagID int64) error
	GetForVideo(ctx context.Context, videoID string) ([]string, error)
	RemoveFromVideo(ctx context.Context, videoID string, tagID int64) error
}

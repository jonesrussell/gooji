package service

import (
	"context"
	"fmt"

	"gooji/internal/domain"
	"gooji/internal/repository"
)

type VideoService struct {
	repo repository.VideoRepository
}

func NewVideoService(repo repository.VideoRepository) *VideoService {
	return &VideoService{repo: repo}
}

func (s *VideoService) Create(ctx context.Context, video *domain.Video) error {
	return s.repo.Create(ctx, video)
}

func (s *VideoService) GetByID(ctx context.Context, id string) (*domain.Video, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *VideoService) ListPublished(ctx context.Context) ([]*domain.Video, error) {
	status := domain.StatusPublished
	return s.repo.List(ctx, &status)
}

func (s *VideoService) ListPending(ctx context.Context) ([]*domain.Video, error) {
	status := domain.StatusPending
	return s.repo.List(ctx, &status)
}

func (s *VideoService) ListAll(ctx context.Context) ([]*domain.Video, error) {
	return s.repo.List(ctx, nil)
}

func (s *VideoService) Approve(ctx context.Context, id string) error {
	video, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get video: %w", err)
	}

	if err := video.Publish(); err != nil {
		return fmt.Errorf("failed to publish: %w", err)
	}

	return s.repo.Update(ctx, video)
}

func (s *VideoService) Reject(ctx context.Context, id string, reason string) error {
	video, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get video: %w", err)
	}

	if err := video.Reject(reason); err != nil {
		return fmt.Errorf("failed to reject: %w", err)
	}

	return s.repo.Update(ctx, video)
}

func (s *VideoService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

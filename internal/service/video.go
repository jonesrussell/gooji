package service

import (
	"context"
	"fmt"
	"html"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"gooji/internal/domain"
	"gooji/internal/repository"
	"gooji/pkg/ffmpeg"
)

type VideoService struct {
	repo      repository.VideoRepository
	storage   repository.FileStorage
	processor *ffmpeg.Processor
	logger    *zap.Logger
	validate  *validator.Validate
}

type UploadRequest struct {
	Title       string `validate:"max=200"`
	Description string `validate:"max=2000"`
}

func NewVideoService(
	repo repository.VideoRepository,
	storage repository.FileStorage,
	processor *ffmpeg.Processor,
	logger *zap.Logger,
) *VideoService {
	return &VideoService{
		repo:      repo,
		storage:   storage,
		processor: processor,
		logger:    logger,
		validate:  validator.New(),
	}
}

func (s *VideoService) Upload(ctx context.Context, file multipart.File, header *multipart.FileHeader, req UploadRequest) (*domain.Video, error) {
	if err := s.validate.Struct(req); err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInvalidInput, err)
	}

	// Validate file type
	if err := s.validateFileType(header); err != nil {
		return nil, err
	}

	// Generate secure filename
	id := uuid.New().String()
	ext := filepath.Ext(header.Filename)
	filename := id + ext

	// Save file
	filePath, err := s.storage.Save(ctx, file, filename)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrStorageFailed, err)
	}

	// Get video metadata
	info, err := s.processor.GetVideoInfo(filePath)
	if err != nil {
		s.storage.Delete(ctx, filePath)
		return nil, fmt.Errorf("%w: %v", domain.ErrProcessingFailed, err)
	}

	// Generate thumbnail
	thumbnailPath := s.storage.GetPath(id + ".jpg")
	if err := s.processor.GenerateThumbnail(filePath, thumbnailPath, 1.0); err != nil {
		s.logger.Warn("failed to generate thumbnail", zap.Error(err))
	}

	video := &domain.Video{
		ID:              id,
		Title:           s.sanitize(req.Title),
		Description:     s.sanitize(req.Description),
		Language:        "oji",
		DurationSeconds: info.Duration,
		FilePath:        filePath,
		ThumbnailPath:   thumbnailPath,
		FileSizeBytes:   header.Size,
		MimeType:        header.Header.Get("Content-Type"),
	}

	if err := s.repo.Create(ctx, video); err != nil {
		s.storage.Delete(ctx, filePath)
		return nil, fmt.Errorf("%w: %v", domain.ErrStorageFailed, err)
	}

	s.logger.Info("video uploaded", zap.String("id", video.ID))
	return video, nil
}

func (s *VideoService) Get(ctx context.Context, id string) (*domain.Video, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *VideoService) List(ctx context.Context, opts domain.VideoListOptions) ([]domain.Video, error) {
	return s.repo.List(ctx, opts)
}

func (s *VideoService) Delete(ctx context.Context, id string) error {
	video, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	s.storage.Delete(ctx, video.FilePath)
	s.storage.Delete(ctx, video.ThumbnailPath)

	s.logger.Info("video deleted", zap.String("id", id))
	return nil
}

func (s *VideoService) validateFileType(header *multipart.FileHeader) error {
	allowed := map[string]bool{
		"video/mp4":  true,
		"video/webm": true,
	}
	contentType := header.Header.Get("Content-Type")
	if !allowed[contentType] {
		return fmt.Errorf("%w: unsupported content type %s", domain.ErrInvalidInput, contentType)
	}
	return nil
}

func (s *VideoService) sanitize(input string) string {
	input = strings.TrimSpace(input)
	return html.EscapeString(input)
}

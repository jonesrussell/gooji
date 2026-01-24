package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"gooji/internal/domain"
	"gooji/internal/repository"
)

type videoRepository struct {
	db *sql.DB
}

func NewVideoRepository(db *sql.DB) repository.VideoRepository {
	return &videoRepository{db: db}
}

func (r *videoRepository) Create(ctx context.Context, video *domain.Video) error {
	query := `
		INSERT INTO videos (id, title, description, language, duration_seconds,
			file_path, thumbnail_path, file_size_bytes, mime_type, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	video.CreatedAt = now
	video.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, query,
		video.ID, video.Title, video.Description, video.Language,
		video.DurationSeconds, video.FilePath, video.ThumbnailPath,
		video.FileSizeBytes, video.MimeType, video.CreatedAt, video.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert video: %w", err)
	}

	return nil
}

func (r *videoRepository) GetByID(ctx context.Context, id string) (*domain.Video, error) {
	query := `
		SELECT id, title, description, language, duration_seconds,
			file_path, thumbnail_path, file_size_bytes, mime_type, created_at, updated_at
		FROM videos WHERE id = ?`

	var video domain.Video
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&video.ID, &video.Title, &video.Description, &video.Language,
		&video.DurationSeconds, &video.FilePath, &video.ThumbnailPath,
		&video.FileSizeBytes, &video.MimeType, &video.CreatedAt, &video.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrVideoNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get video: %w", err)
	}

	return &video, nil
}

func (r *videoRepository) List(ctx context.Context, opts domain.VideoListOptions) ([]domain.Video, error) {
	query := `
		SELECT id, title, description, language, duration_seconds,
			file_path, thumbnail_path, file_size_bytes, mime_type, created_at, updated_at
		FROM videos ORDER BY created_at DESC LIMIT ? OFFSET ?`

	limit := opts.Limit
	if limit <= 0 {
		limit = 50
	}

	rows, err := r.db.QueryContext(ctx, query, limit, opts.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list videos: %w", err)
	}
	defer rows.Close()

	var videos []domain.Video
	for rows.Next() {
		var v domain.Video
		if err := rows.Scan(&v.ID, &v.Title, &v.Description, &v.Language,
			&v.DurationSeconds, &v.FilePath, &v.ThumbnailPath,
			&v.FileSizeBytes, &v.MimeType, &v.CreatedAt, &v.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan video: %w", err)
		}
		videos = append(videos, v)
	}

	return videos, nil
}

func (r *videoRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM videos WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete video: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return domain.ErrVideoNotFound
	}

	return nil
}

func (r *videoRepository) Update(ctx context.Context, video *domain.Video) error {
	query := `
		UPDATE videos SET title = ?, description = ?, updated_at = ?
		WHERE id = ?`

	video.UpdatedAt = time.Now()
	result, err := r.db.ExecContext(ctx, query, video.Title, video.Description, video.UpdatedAt, video.ID)
	if err != nil {
		return fmt.Errorf("failed to update video: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return domain.ErrVideoNotFound
	}

	return nil
}

package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"gooji/internal/domain"
)

type VideoRepository struct {
	db *DB
}

func NewVideoRepository(db *DB) *VideoRepository {
	return &VideoRepository{db: db}
}

func (r *VideoRepository) Create(ctx context.Context, video *domain.Video) error {
	query := `
		INSERT INTO videos (
			id, title, description, language, duration_seconds,
			file_path, thumbnail_path, file_size_bytes, mime_type,
			status, rejection_reason, device_id, created_at, published_at, moderated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query,
		video.ID, video.Title, video.Description, video.Language, video.DurationSeconds,
		video.FilePath, video.ThumbnailPath, video.FileSizeBytes, video.MimeType,
		video.Status, video.RejectionReason, video.DeviceID, video.CreatedAt,
		video.PublishedAt, video.ModeratedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert video: %w", err)
	}
	return nil
}

func (r *VideoRepository) GetByID(ctx context.Context, id string) (*domain.Video, error) {
	query := `
		SELECT id, title, description, language, duration_seconds,
			file_path, thumbnail_path, file_size_bytes, mime_type,
			status, rejection_reason, device_id, created_at, published_at, moderated_at
		FROM videos WHERE id = ?`

	video := &domain.Video{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&video.ID, &video.Title, &video.Description, &video.Language, &video.DurationSeconds,
		&video.FilePath, &video.ThumbnailPath, &video.FileSizeBytes, &video.MimeType,
		&video.Status, &video.RejectionReason, &video.DeviceID, &video.CreatedAt,
		&video.PublishedAt, &video.ModeratedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrVideoNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get video: %w", err)
	}
	return video, nil
}

func (r *VideoRepository) List(ctx context.Context, status *domain.VideoStatus) ([]*domain.Video, error) {
	var query string
	var args []any

	if status != nil {
		query = `
			SELECT id, title, description, language, duration_seconds,
				file_path, thumbnail_path, file_size_bytes, mime_type,
				status, rejection_reason, device_id, created_at, published_at, moderated_at
			FROM videos WHERE status = ? ORDER BY created_at DESC`
		args = append(args, *status)
	} else {
		query = `
			SELECT id, title, description, language, duration_seconds,
				file_path, thumbnail_path, file_size_bytes, mime_type,
				status, rejection_reason, device_id, created_at, published_at, moderated_at
			FROM videos ORDER BY created_at DESC`
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list videos: %w", err)
	}
	defer rows.Close()

	var videos []*domain.Video
	for rows.Next() {
		video := &domain.Video{}
		err := rows.Scan(
			&video.ID, &video.Title, &video.Description, &video.Language, &video.DurationSeconds,
			&video.FilePath, &video.ThumbnailPath, &video.FileSizeBytes, &video.MimeType,
			&video.Status, &video.RejectionReason, &video.DeviceID, &video.CreatedAt,
			&video.PublishedAt, &video.ModeratedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan video: %w", err)
		}
		videos = append(videos, video)
	}

	return videos, rows.Err()
}

func (r *VideoRepository) Update(ctx context.Context, video *domain.Video) error {
	query := `
		UPDATE videos SET
			title = ?, description = ?, language = ?, duration_seconds = ?,
			file_path = ?, thumbnail_path = ?, file_size_bytes = ?, mime_type = ?,
			status = ?, rejection_reason = ?, device_id = ?,
			published_at = ?, moderated_at = ?
		WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query,
		video.Title, video.Description, video.Language, video.DurationSeconds,
		video.FilePath, video.ThumbnailPath, video.FileSizeBytes, video.MimeType,
		video.Status, video.RejectionReason, video.DeviceID,
		video.PublishedAt, video.ModeratedAt, video.ID,
	)
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

func (r *VideoRepository) Delete(ctx context.Context, id string) error {
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

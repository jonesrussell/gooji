package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type DB struct {
	*sql.DB
}

func NewDB(dsn string) (*DB, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enable foreign keys and WAL mode for better concurrency
	pragmas := []string{
		"PRAGMA foreign_keys = ON",
		"PRAGMA journal_mode = WAL",
		"PRAGMA busy_timeout = 5000",
	}
	for _, pragma := range pragmas {
		if _, err := db.Exec(pragma); err != nil {
			db.Close()
			return nil, fmt.Errorf("failed to set pragma: %w", err)
		}
	}

	// Run migrations
	if err := runMigrations(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return &DB{db}, nil
}

func runMigrations(db *sql.DB) error {
	if _, err := db.Exec(initialMigration); err != nil {
		return fmt.Errorf("failed to execute migration: %w", err)
	}
	return nil
}

func (db *DB) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return db.DB.QueryRowContext(ctx, query, args...)
}

const initialMigration = `
CREATE TABLE IF NOT EXISTS videos (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT DEFAULT '',
    language TEXT DEFAULT 'oji',
    duration_seconds INTEGER DEFAULT 0,
    file_path TEXT NOT NULL,
    thumbnail_path TEXT DEFAULT '',
    file_size_bytes INTEGER DEFAULT 0,
    mime_type TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    rejection_reason TEXT DEFAULT '',
    device_id TEXT,
    created_at DATETIME NOT NULL,
    published_at DATETIME,
    moderated_at DATETIME,
    FOREIGN KEY (device_id) REFERENCES devices(id)
);

CREATE TABLE IF NOT EXISTS devices (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    api_key_hash TEXT NOT NULL,
    trust_level TEXT NOT NULL DEFAULT 'trusted',
    last_seen_at DATETIME,
    created_at DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS tags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS video_tags (
    video_id TEXT NOT NULL,
    tag_id INTEGER NOT NULL,
    PRIMARY KEY (video_id, tag_id),
    FOREIGN KEY (video_id) REFERENCES videos(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_videos_status ON videos(status);
CREATE INDEX IF NOT EXISTS idx_videos_created ON videos(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_videos_device ON videos(device_id);
`

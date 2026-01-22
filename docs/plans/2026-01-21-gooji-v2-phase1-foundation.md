# Gooji v2 Phase 1: Foundation Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Scaffold the clean architecture foundation with Go 1.25, SQLite, Vue 3 + Inertia.js, serving an empty gallery page.

**Architecture:** Clean architecture with domain → service → repository → handler layers. Inertia.js bridges Go backend and Vue frontend, enabling SPA-like experience with server-driven routing. SQLite for persistence, Vite for frontend bundling.

**Tech Stack:** Go 1.25, SQLite (modernc.org/sqlite - pure Go), Inertia.js (petaki/inertia-go), Vue 3 Composition API, TypeScript, Vite, Tailwind CSS.

---

## Phase 1 Overview

This phase establishes the foundation. At the end, the app boots and serves an empty gallery via Inertia/Vue.

**Tasks:**
1. Clean project scaffold
2. Domain entities
3. SQLite repository setup
4. Video repository implementation
5. Video service
6. Inertia.js handler setup
7. Vue 3 + Vite frontend scaffold
8. Gallery page (empty state)
9. Health check endpoint
10. Integration test

---

### Task 1: Clean Project Scaffold

**Files:**
- Delete: All existing Go files, templates, static assets
- Create: `cmd/server/main.go`
- Create: `cmd/admin/main.go`
- Modify: `go.mod`

**Step 1: Remove old code structure**

```bash
rm -rf internal/ pkg/ web/ main.go
rm -f config/config.json
```

**Step 2: Create new directory structure**

```bash
mkdir -p cmd/server cmd/admin
mkdir -p internal/domain internal/service internal/repository/sqlite internal/repository/filesystem
mkdir -p internal/handler internal/middleware internal/ffmpeg
mkdir -p migrations
mkdir -p web/src/pages web/src/components web/src/composables web/src/types web/src/assets/css
```

**Step 3: Update go.mod for Go 1.25**

```go
module gooji

go 1.25

require (
	github.com/petaki/inertia-go v1.6.0
	go.uber.org/zap v1.27.0
	modernc.org/sqlite v1.34.5
	github.com/spf13/cobra v1.8.1
	github.com/go-playground/validator/v10 v10.23.0
	github.com/google/uuid v1.6.0
	golang.org/x/crypto v0.31.0
)
```

**Step 4: Create minimal server entry point**

Create `cmd/server/main.go`:

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
```

**Step 5: Create minimal admin CLI entry point**

Create `cmd/admin/main.go`:

```go
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "gooji-admin",
		Short: "Gooji administration tool",
	}

	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("gooji-admin v2.0.0")
		},
	})

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
```

**Step 6: Download dependencies and verify**

```bash
go mod tidy
go build ./cmd/server
go build ./cmd/admin
```

Expected: Both binaries compile without errors.

**Step 7: Run server and test health endpoint**

```bash
./server &
curl http://localhost:8080/health
pkill server
```

Expected: `{"status":"ok"}`

**Step 8: Commit**

```bash
git add -A
git commit -m "feat: scaffold clean project structure for v2

- Remove v1 code (templates, handlers, old architecture)
- Add cmd/server and cmd/admin entry points
- Update go.mod for Go 1.25 with new dependencies
- Add basic health check endpoint"
```

---

### Task 2: Domain Entities

**Files:**
- Create: `internal/domain/video.go`
- Create: `internal/domain/device.go`
- Create: `internal/domain/errors.go`
- Test: `internal/domain/video_test.go`

**Step 1: Write failing test for Video entity**

Create `internal/domain/video_test.go`:

```go
package domain_test

import (
	"testing"
	"time"

	"gooji/internal/domain"
)

func TestNewVideo(t *testing.T) {
	v := domain.NewVideo("Test Title", "test.mp4", "video/mp4")

	if v.ID == "" {
		t.Error("expected ID to be generated")
	}
	if v.Title != "Test Title" {
		t.Errorf("expected title 'Test Title', got '%s'", v.Title)
	}
	if v.Status != domain.StatusPending {
		t.Errorf("expected status Pending, got %s", v.Status)
	}
	if v.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}

func TestVideoPublish(t *testing.T) {
	v := domain.NewVideo("Test", "test.mp4", "video/mp4")

	err := v.Publish()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if v.Status != domain.StatusPublished {
		t.Errorf("expected status Published, got %s", v.Status)
	}
	if v.PublishedAt == nil {
		t.Error("expected PublishedAt to be set")
	}
}

func TestVideoReject(t *testing.T) {
	v := domain.NewVideo("Test", "test.mp4", "video/mp4")

	err := v.Reject("Inappropriate content")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if v.Status != domain.StatusRejected {
		t.Errorf("expected status Rejected, got %s", v.Status)
	}
	if v.RejectionReason != "Inappropriate content" {
		t.Errorf("expected rejection reason, got '%s'", v.RejectionReason)
	}
}

func TestVideoCannotPublishIfAlreadyRejected(t *testing.T) {
	v := domain.NewVideo("Test", "test.mp4", "video/mp4")
	v.Reject("Bad")

	err := v.Publish()
	if err == nil {
		t.Error("expected error when publishing rejected video")
	}
}
```

**Step 2: Run test to verify it fails**

```bash
go test ./internal/domain/... -v
```

Expected: FAIL - package doesn't exist yet.

**Step 3: Implement Video entity**

Create `internal/domain/video.go`:

```go
package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type VideoStatus string

const (
	StatusPending   VideoStatus = "pending"
	StatusPublished VideoStatus = "published"
	StatusRejected  VideoStatus = "rejected"
)

type Video struct {
	ID              string
	Title           string
	Description     string
	Language        string
	DurationSeconds int
	FilePath        string
	ThumbnailPath   string
	FileSizeBytes   int64
	MimeType        string
	Status          VideoStatus
	RejectionReason string
	DeviceID        *string
	CreatedAt       time.Time
	PublishedAt     *time.Time
	ModeratedAt     *time.Time
}

func NewVideo(title, filePath, mimeType string) *Video {
	now := time.Now()
	return &Video{
		ID:        uuid.New().String(),
		Title:     title,
		FilePath:  filePath,
		MimeType:  mimeType,
		Language:  "oji",
		Status:    StatusPending,
		CreatedAt: now,
	}
}

func (v *Video) Publish() error {
	if v.Status == StatusRejected {
		return errors.New("cannot publish rejected video")
	}
	now := time.Now()
	v.Status = StatusPublished
	v.PublishedAt = &now
	v.ModeratedAt = &now
	return nil
}

func (v *Video) Reject(reason string) error {
	if v.Status == StatusPublished {
		return errors.New("cannot reject published video")
	}
	now := time.Now()
	v.Status = StatusRejected
	v.RejectionReason = reason
	v.ModeratedAt = &now
	return nil
}

func (v *Video) IsPending() bool {
	return v.Status == StatusPending
}

func (v *Video) IsPublished() bool {
	return v.Status == StatusPublished
}
```

**Step 4: Run test to verify it passes**

```bash
go test ./internal/domain/... -v
```

Expected: PASS (4 tests)

**Step 5: Create Device entity**

Create `internal/domain/device.go`:

```go
package domain

import (
	"time"

	"github.com/google/uuid"
)

type TrustLevel string

const (
	TrustLevelTrusted   TrustLevel = "trusted"
	TrustLevelUntrusted TrustLevel = "untrusted"
)

type Device struct {
	ID          string
	Name        string
	APIKeyHash  string
	TrustLevel  TrustLevel
	LastSeenAt  *time.Time
	CreatedAt   time.Time
}

func NewDevice(name string, apiKeyHash string) *Device {
	return &Device{
		ID:         uuid.New().String(),
		Name:       name,
		APIKeyHash: apiKeyHash,
		TrustLevel: TrustLevelTrusted,
		CreatedAt:  time.Now(),
	}
}

func (d *Device) IsTrusted() bool {
	return d.TrustLevel == TrustLevelTrusted
}

func (d *Device) UpdateLastSeen() {
	now := time.Now()
	d.LastSeenAt = &now
}
```

**Step 6: Create domain errors**

Create `internal/domain/errors.go`:

```go
package domain

import "errors"

var (
	ErrVideoNotFound    = errors.New("video not found")
	ErrDeviceNotFound   = errors.New("device not found")
	ErrInvalidAPIKey    = errors.New("invalid API key")
	ErrInvalidStatus    = errors.New("invalid status transition")
	ErrValidationFailed = errors.New("validation failed")
)
```

**Step 7: Commit**

```bash
git add internal/domain/
git commit -m "feat: add domain entities (Video, Device)

- Video with status transitions (pending → published/rejected)
- Device with trust levels for kiosk authentication
- Domain errors for common failure cases
- Unit tests for Video state machine"
```

---

### Task 3: SQLite Repository Setup

**Files:**
- Create: `internal/repository/interfaces.go`
- Create: `internal/repository/sqlite/db.go`
- Create: `migrations/001_initial.sql`
- Test: `internal/repository/sqlite/db_test.go`

**Step 1: Create repository interfaces**

Create `internal/repository/interfaces.go`:

```go
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
```

**Step 2: Create migration file**

Create `migrations/001_initial.sql`:

```sql
-- Videos: core content entity
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

-- Devices: trust tracking for kiosks
CREATE TABLE IF NOT EXISTS devices (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    api_key_hash TEXT NOT NULL,
    trust_level TEXT NOT NULL DEFAULT 'trusted',
    last_seen_at DATETIME,
    created_at DATETIME NOT NULL
);

-- Tags: categorization
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

-- Indexes for common queries
CREATE INDEX IF NOT EXISTS idx_videos_status ON videos(status);
CREATE INDEX IF NOT EXISTS idx_videos_created ON videos(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_videos_device ON videos(device_id);
```

**Step 3: Write failing test for DB connection**

Create `internal/repository/sqlite/db_test.go`:

```go
package sqlite_test

import (
	"context"
	"testing"

	"gooji/internal/repository/sqlite"
)

func TestNewDB(t *testing.T) {
	db, err := sqlite.NewDB(":memory:")
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}
	defer db.Close()

	// Verify tables exist
	ctx := context.Background()
	var tableName string
	err = db.QueryRowContext(ctx, "SELECT name FROM sqlite_master WHERE type='table' AND name='videos'").Scan(&tableName)
	if err != nil {
		t.Errorf("videos table not found: %v", err)
	}
}
```

**Step 4: Run test to verify it fails**

```bash
go test ./internal/repository/sqlite/... -v
```

Expected: FAIL - package doesn't exist.

**Step 5: Implement DB setup**

Create `internal/repository/sqlite/db.go`:

```go
package sqlite

import (
	"context"
	"database/sql"
	"embed"
	"fmt"

	_ "modernc.org/sqlite"
)

//go:embed ../../../migrations/*.sql
var migrationsFS embed.FS

type DB struct {
	*sql.DB
}

func NewDB(dsn string) (*DB, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	// Run migrations
	if err := runMigrations(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return &DB{db}, nil
}

func runMigrations(db *sql.DB) error {
	migration, err := migrationsFS.ReadFile("migrations/001_initial.sql")
	if err != nil {
		return fmt.Errorf("failed to read migration: %w", err)
	}

	if _, err := db.Exec(string(migration)); err != nil {
		return fmt.Errorf("failed to execute migration: %w", err)
	}

	return nil
}

func (db *DB) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return db.DB.QueryRowContext(ctx, query, args...)
}
```

**Step 6: Fix embed path issue**

The embed path is relative to the file location. Update `internal/repository/sqlite/db.go`:

```go
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

func NewDBWithMigrations(dsn string, migrations string) (*DB, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

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

	if _, err := db.Exec(migrations); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to execute migrations: %w", err)
	}

	return &DB{db}, nil
}

func runMigrations(db *sql.DB) error {
	migration := initialMigration
	if _, err := db.Exec(migration); err != nil {
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
```

**Step 7: Run test to verify it passes**

```bash
go test ./internal/repository/sqlite/... -v
```

Expected: PASS

**Step 8: Commit**

```bash
git add internal/repository/ migrations/
git commit -m "feat: add SQLite repository setup with migrations

- Repository interfaces (Video, Device, Tag)
- SQLite connection with WAL mode and foreign keys
- Initial migration schema
- Test for DB initialization"
```

---

### Task 4: Video Repository Implementation

**Files:**
- Create: `internal/repository/sqlite/video.go`
- Test: `internal/repository/sqlite/video_test.go`

**Step 1: Write failing tests for VideoRepository**

Create `internal/repository/sqlite/video_test.go`:

```go
package sqlite_test

import (
	"context"
	"testing"

	"gooji/internal/domain"
	"gooji/internal/repository/sqlite"
)

func setupTestDB(t *testing.T) *sqlite.DB {
	db, err := sqlite.NewDB(":memory:")
	if err != nil {
		t.Fatalf("failed to create test db: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func TestVideoRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := sqlite.NewVideoRepository(db)
	ctx := context.Background()

	video := domain.NewVideo("Test Video", "/videos/test.mp4", "video/mp4")
	video.Description = "A test video"

	err := repo.Create(ctx, video)
	if err != nil {
		t.Fatalf("failed to create video: %v", err)
	}

	// Verify it was saved
	saved, err := repo.GetByID(ctx, video.ID)
	if err != nil {
		t.Fatalf("failed to get video: %v", err)
	}
	if saved.Title != "Test Video" {
		t.Errorf("expected title 'Test Video', got '%s'", saved.Title)
	}
	if saved.Description != "A test video" {
		t.Errorf("expected description 'A test video', got '%s'", saved.Description)
	}
}

func TestVideoRepository_List(t *testing.T) {
	db := setupTestDB(t)
	repo := sqlite.NewVideoRepository(db)
	ctx := context.Background()

	// Create videos with different statuses
	v1 := domain.NewVideo("Video 1", "/v1.mp4", "video/mp4")
	v2 := domain.NewVideo("Video 2", "/v2.mp4", "video/mp4")
	v2.Publish()
	v3 := domain.NewVideo("Video 3", "/v3.mp4", "video/mp4")

	repo.Create(ctx, v1)
	repo.Create(ctx, v2)
	repo.Create(ctx, v3)

	// List all
	all, err := repo.List(ctx, nil)
	if err != nil {
		t.Fatalf("failed to list all: %v", err)
	}
	if len(all) != 3 {
		t.Errorf("expected 3 videos, got %d", len(all))
	}

	// List pending only
	pending := domain.StatusPending
	pendingVideos, err := repo.List(ctx, &pending)
	if err != nil {
		t.Fatalf("failed to list pending: %v", err)
	}
	if len(pendingVideos) != 2 {
		t.Errorf("expected 2 pending videos, got %d", len(pendingVideos))
	}
}

func TestVideoRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := sqlite.NewVideoRepository(db)
	ctx := context.Background()

	video := domain.NewVideo("Original", "/test.mp4", "video/mp4")
	repo.Create(ctx, video)

	video.Title = "Updated"
	video.Publish()
	err := repo.Update(ctx, video)
	if err != nil {
		t.Fatalf("failed to update: %v", err)
	}

	updated, _ := repo.GetByID(ctx, video.ID)
	if updated.Title != "Updated" {
		t.Errorf("expected 'Updated', got '%s'", updated.Title)
	}
	if updated.Status != domain.StatusPublished {
		t.Errorf("expected Published, got %s", updated.Status)
	}
}

func TestVideoRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := sqlite.NewVideoRepository(db)
	ctx := context.Background()

	video := domain.NewVideo("To Delete", "/test.mp4", "video/mp4")
	repo.Create(ctx, video)

	err := repo.Delete(ctx, video.ID)
	if err != nil {
		t.Fatalf("failed to delete: %v", err)
	}

	_, err = repo.GetByID(ctx, video.ID)
	if err != domain.ErrVideoNotFound {
		t.Errorf("expected ErrVideoNotFound, got %v", err)
	}
}
```

**Step 2: Run test to verify it fails**

```bash
go test ./internal/repository/sqlite/... -v
```

Expected: FAIL - NewVideoRepository doesn't exist.

**Step 3: Implement VideoRepository**

Create `internal/repository/sqlite/video.go`:

```go
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
```

**Step 4: Run tests to verify they pass**

```bash
go test ./internal/repository/sqlite/... -v
```

Expected: PASS (5 tests)

**Step 5: Commit**

```bash
git add internal/repository/sqlite/video.go internal/repository/sqlite/video_test.go
git commit -m "feat: implement SQLite video repository

- CRUD operations for videos
- Filter by status in List()
- Proper error handling with domain errors
- Comprehensive unit tests"
```

---

### Task 5: Video Service

**Files:**
- Create: `internal/service/video.go`
- Test: `internal/service/video_test.go`

**Step 1: Write failing tests for VideoService**

Create `internal/service/video_test.go`:

```go
package service_test

import (
	"context"
	"testing"

	"gooji/internal/domain"
	"gooji/internal/repository/sqlite"
	"gooji/internal/service"
)

func setupTestService(t *testing.T) *service.VideoService {
	db, err := sqlite.NewDB(":memory:")
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	repo := sqlite.NewVideoRepository(db)
	return service.NewVideoService(repo)
}

func TestVideoService_ListPublished(t *testing.T) {
	svc := setupTestService(t)
	ctx := context.Background()

	// Create mix of published and pending
	v1 := domain.NewVideo("Published", "/v1.mp4", "video/mp4")
	v1.Publish()
	v2 := domain.NewVideo("Pending", "/v2.mp4", "video/mp4")

	svc.Create(ctx, v1)
	svc.Create(ctx, v2)

	published, err := svc.ListPublished(ctx)
	if err != nil {
		t.Fatalf("failed to list: %v", err)
	}
	if len(published) != 1 {
		t.Errorf("expected 1 published, got %d", len(published))
	}
}

func TestVideoService_ListPending(t *testing.T) {
	svc := setupTestService(t)
	ctx := context.Background()

	v1 := domain.NewVideo("Pending 1", "/v1.mp4", "video/mp4")
	v2 := domain.NewVideo("Pending 2", "/v2.mp4", "video/mp4")
	v3 := domain.NewVideo("Published", "/v3.mp4", "video/mp4")
	v3.Publish()

	svc.Create(ctx, v1)
	svc.Create(ctx, v2)
	svc.Create(ctx, v3)

	pending, err := svc.ListPending(ctx)
	if err != nil {
		t.Fatalf("failed to list: %v", err)
	}
	if len(pending) != 2 {
		t.Errorf("expected 2 pending, got %d", len(pending))
	}
}

func TestVideoService_Approve(t *testing.T) {
	svc := setupTestService(t)
	ctx := context.Background()

	video := domain.NewVideo("Test", "/test.mp4", "video/mp4")
	svc.Create(ctx, video)

	err := svc.Approve(ctx, video.ID)
	if err != nil {
		t.Fatalf("failed to approve: %v", err)
	}

	approved, _ := svc.GetByID(ctx, video.ID)
	if approved.Status != domain.StatusPublished {
		t.Errorf("expected Published, got %s", approved.Status)
	}
}

func TestVideoService_Reject(t *testing.T) {
	svc := setupTestService(t)
	ctx := context.Background()

	video := domain.NewVideo("Test", "/test.mp4", "video/mp4")
	svc.Create(ctx, video)

	err := svc.Reject(ctx, video.ID, "Inappropriate content")
	if err != nil {
		t.Fatalf("failed to reject: %v", err)
	}

	rejected, _ := svc.GetByID(ctx, video.ID)
	if rejected.Status != domain.StatusRejected {
		t.Errorf("expected Rejected, got %s", rejected.Status)
	}
	if rejected.RejectionReason != "Inappropriate content" {
		t.Errorf("expected reason, got '%s'", rejected.RejectionReason)
	}
}
```

**Step 2: Run test to verify it fails**

```bash
go test ./internal/service/... -v
```

Expected: FAIL - package doesn't exist.

**Step 3: Implement VideoService**

Create `internal/service/video.go`:

```go
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
```

**Step 4: Run tests to verify they pass**

```bash
go test ./internal/service/... -v
```

Expected: PASS (4 tests)

**Step 5: Commit**

```bash
git add internal/service/
git commit -m "feat: add video service with moderation operations

- ListPublished/ListPending for gallery and moderation queue
- Approve/Reject for moderation workflow
- Business logic layer between handlers and repository"
```

---

### Task 6: Inertia.js Handler Setup

**Files:**
- Create: `internal/handler/handler.go`
- Create: `internal/handler/gallery.go`
- Modify: `cmd/server/main.go`

**Step 1: Create base handler with Inertia setup**

Create `internal/handler/handler.go`:

```go
package handler

import (
	"net/http"

	inertia "github.com/petaki/inertia-go"

	"gooji/internal/service"
)

type Handler struct {
	inertia      *inertia.Inertia
	videoService *service.VideoService
}

func New(videoService *service.VideoService, viteManifestPath string) *Handler {
	i := inertia.New("http://localhost:5173", "app.html", viteManifestPath)

	return &Handler{
		inertia:      i,
		videoService: videoService,
	}
}

func (h *Handler) SetRootTemplate(tmpl string) {
	h.inertia.SetRootTemplate(tmpl)
}

func (h *Handler) ShareFunc(r *http.Request) map[string]any {
	return map[string]any{
		"appName": "Gooji",
	}
}

func (h *Handler) Inertia() *inertia.Inertia {
	return h.inertia
}
```

**Step 2: Create gallery handler**

Create `internal/handler/gallery.go`:

```go
package handler

import (
	"net/http"
)

type VideoDTO struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	ThumbnailPath string `json:"thumbnailPath"`
	DurationSeconds int  `json:"durationSeconds"`
}

func (h *Handler) Gallery(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	videos, err := h.videoService.ListPublished(ctx)
	if err != nil {
		http.Error(w, "Failed to load videos", http.StatusInternalServerError)
		return
	}

	// Convert to DTOs
	dtos := make([]VideoDTO, 0, len(videos))
	for _, v := range videos {
		dtos = append(dtos, VideoDTO{
			ID:              v.ID,
			Title:           v.Title,
			Description:     v.Description,
			ThumbnailPath:   v.ThumbnailPath,
			DurationSeconds: v.DurationSeconds,
		})
	}

	err = h.inertia.Render(w, r, "Gallery", map[string]any{
		"videos": dtos,
	})
	if err != nil {
		http.Error(w, "Failed to render", http.StatusInternalServerError)
	}
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	err := h.inertia.Render(w, r, "Home", map[string]any{})
	if err != nil {
		http.Error(w, "Failed to render", http.StatusInternalServerError)
	}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
```

**Step 3: Update main.go to wire everything together**

Update `cmd/server/main.go`:

```go
package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"gooji/internal/handler"
	"gooji/internal/repository/sqlite"
	"gooji/internal/service"
)

func main() {
	// Configuration
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "gooji.db"
	}

	// Database
	db, err := sqlite.NewDB(dbPath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Repositories
	videoRepo := sqlite.NewVideoRepository(db)

	// Services
	videoService := service.NewVideoService(videoRepo)

	// Handlers
	manifestPath := filepath.Join("web", "dist", ".vite", "manifest.json")
	h := handler.New(videoService, manifestPath)

	// Routes
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", h.Home)
	mux.HandleFunc("GET /gallery", h.Gallery)
	mux.HandleFunc("GET /health", h.Health)

	// Static files (Vite build output)
	staticDir := filepath.Join("web", "dist", "assets")
	mux.Handle("GET /assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(staticDir))))

	// Start server
	addr := ":" + port
	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
```

**Step 4: Verify it compiles**

```bash
go build ./cmd/server
```

Expected: Compiles without errors.

**Step 5: Commit**

```bash
git add internal/handler/ cmd/server/main.go
git commit -m "feat: add Inertia.js handler setup with gallery route

- Base handler with Inertia integration
- Gallery handler serving published videos
- Home and health endpoints
- Wired up SQLite -> Repository -> Service -> Handler"
```

---

### Task 7: Vue 3 + Vite Frontend Scaffold

**Files:**
- Create: `web/package.json`
- Create: `web/vite.config.ts`
- Create: `web/tsconfig.json`
- Create: `web/tailwind.config.js`
- Create: `web/postcss.config.js`
- Create: `web/index.html`
- Create: `web/src/app.ts`
- Create: `web/src/types/index.ts`
- Create: `web/src/assets/css/app.css`

**Step 1: Create package.json**

Create `web/package.json`:

```json
{
  "name": "gooji-web",
  "private": true,
  "version": "2.0.0",
  "type": "module",
  "scripts": {
    "dev": "vite",
    "build": "vue-tsc && vite build",
    "preview": "vite preview",
    "typecheck": "vue-tsc --noEmit"
  },
  "dependencies": {
    "@inertiajs/vue3": "^1.2.0",
    "vue": "^3.5.0"
  },
  "devDependencies": {
    "@vitejs/plugin-vue": "^5.2.0",
    "@types/node": "^22.10.0",
    "autoprefixer": "^10.4.20",
    "postcss": "^8.4.49",
    "tailwindcss": "^3.4.17",
    "typescript": "^5.7.0",
    "vite": "^6.0.0",
    "vue-tsc": "^2.2.0"
  }
}
```

**Step 2: Create vite.config.ts**

Create `web/vite.config.ts`:

```typescript
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src'),
    },
  },
  build: {
    manifest: true,
    outDir: 'dist',
    rollupOptions: {
      input: resolve(__dirname, 'src/app.ts'),
    },
  },
  server: {
    origin: 'http://localhost:5173',
  },
})
```

**Step 3: Create tsconfig.json**

Create `web/tsconfig.json`:

```json
{
  "compilerOptions": {
    "target": "ES2022",
    "useDefineForClassFields": true,
    "module": "ESNext",
    "lib": ["ES2022", "DOM", "DOM.Iterable"],
    "skipLibCheck": true,
    "moduleResolution": "bundler",
    "allowImportingTsExtensions": true,
    "resolveJsonModule": true,
    "isolatedModules": true,
    "noEmit": true,
    "jsx": "preserve",
    "strict": true,
    "noUnusedLocals": true,
    "noUnusedParameters": true,
    "noFallthroughCasesInSwitch": true,
    "paths": {
      "@/*": ["./src/*"]
    }
  },
  "include": ["src/**/*.ts", "src/**/*.tsx", "src/**/*.vue"],
  "references": [{ "path": "./tsconfig.node.json" }]
}
```

Create `web/tsconfig.node.json`:

```json
{
  "compilerOptions": {
    "composite": true,
    "skipLibCheck": true,
    "module": "ESNext",
    "moduleResolution": "bundler",
    "allowSyntheticDefaultImports": true
  },
  "include": ["vite.config.ts"]
}
```

**Step 4: Create Tailwind config**

Create `web/tailwind.config.js`:

```javascript
/** @type {import('tailwindcss').Config} */
export default {
  content: [
    './index.html',
    './src/**/*.{vue,js,ts,jsx,tsx}',
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          50: '#f0f9ff',
          100: '#e0f2fe',
          500: '#0ea5e9',
          600: '#0284c7',
          700: '#0369a1',
        },
      },
    },
  },
  plugins: [],
}
```

Create `web/postcss.config.js`:

```javascript
export default {
  plugins: {
    tailwindcss: {},
    autoprefixer: {},
  },
}
```

**Step 5: Create base CSS**

Create `web/src/assets/css/app.css`:

```css
@tailwind base;
@tailwind components;
@tailwind utilities;

@layer base {
  body {
    @apply bg-gray-50 text-gray-900 antialiased;
  }
}

@layer components {
  .btn {
    @apply px-4 py-2 rounded-lg font-medium transition-colors;
  }
  .btn-primary {
    @apply bg-primary-600 text-white hover:bg-primary-700;
  }
}
```

**Step 6: Create TypeScript types**

Create `web/src/types/index.ts`:

```typescript
export interface Video {
  id: string
  title: string
  description: string
  thumbnailPath: string
  durationSeconds: number
}

export interface PageProps {
  appName: string
}
```

**Step 7: Create Inertia app entry**

Create `web/src/app.ts`:

```typescript
import { createApp, h } from 'vue'
import { createInertiaApp } from '@inertiajs/vue3'
import './assets/css/app.css'

createInertiaApp({
  resolve: (name: string) => {
    const pages = import.meta.glob('./pages/**/*.vue', { eager: true })
    return pages[`./pages/${name}.vue`]
  },
  setup({ el, App, props, plugin }) {
    createApp({ render: () => h(App, props) })
      .use(plugin)
      .mount(el)
  },
})
```

**Step 8: Create index.html (Inertia root template)**

Create `web/index.html`:

```html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Gooji - Ojibwe Language & Culture</title>
</head>
<body>
  <div id="app" data-page='{{ .page }}'></div>
  <script type="module" src="http://localhost:5173/src/app.ts"></script>
</body>
</html>
```

Also create `web/app.html` for Go template:

```html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Gooji - Ojibwe Language & Culture</title>
  {{ if .isDev }}
  <script type="module" src="http://localhost:5173/@vite/client"></script>
  <script type="module" src="http://localhost:5173/src/app.ts"></script>
  {{ else }}
  {{ range .cssFiles }}<link rel="stylesheet" href="/assets/{{ . }}">{{ end }}
  {{ range .jsFiles }}<script type="module" src="/assets/{{ . }}"></script>{{ end }}
  {{ end }}
</head>
<body>
  <div id="app" data-page='{{ .page }}'></div>
</body>
</html>
```

**Step 9: Install dependencies and verify**

```bash
cd web && npm install && cd ..
```

Expected: Dependencies installed without errors.

**Step 10: Commit**

```bash
git add web/
git commit -m "feat: scaffold Vue 3 + Vite + Tailwind frontend

- Vite config with Vue plugin and manifest generation
- TypeScript config with strict mode
- Tailwind CSS with custom color palette
- Inertia.js Vue 3 app entry point
- Base types for Video entity"
```

---

### Task 8: Gallery Page (Empty State)

**Files:**
- Create: `web/src/pages/Home.vue`
- Create: `web/src/pages/Gallery.vue`
- Create: `web/src/components/VideoCard.vue`
- Create: `web/src/components/ui/Container.vue`

**Step 1: Create Container component**

Create `web/src/components/ui/Container.vue`:

```vue
<script setup lang="ts">
defineProps<{
  size?: 'sm' | 'md' | 'lg' | 'xl'
}>()
</script>

<template>
  <div
    class="mx-auto px-4 sm:px-6 lg:px-8"
    :class="{
      'max-w-2xl': size === 'sm',
      'max-w-4xl': size === 'md',
      'max-w-6xl': size === 'lg' || !size,
      'max-w-7xl': size === 'xl',
    }"
  >
    <slot />
  </div>
</template>
```

**Step 2: Create VideoCard component**

Create `web/src/components/VideoCard.vue`:

```vue
<script setup lang="ts">
import type { Video } from '@/types'

defineProps<{
  video: Video
}>()

function formatDuration(seconds: number): string {
  const mins = Math.floor(seconds / 60)
  const secs = seconds % 60
  return `${mins}:${secs.toString().padStart(2, '0')}`
}
</script>

<template>
  <article class="group relative bg-white rounded-lg shadow-sm overflow-hidden hover:shadow-md transition-shadow">
    <div class="aspect-video bg-gray-200 relative">
      <img
        v-if="video.thumbnailPath"
        :src="video.thumbnailPath"
        :alt="video.title"
        class="w-full h-full object-cover"
      />
      <div v-else class="w-full h-full flex items-center justify-center text-gray-400">
        <svg class="w-12 h-12" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" />
        </svg>
      </div>
      <div class="absolute bottom-2 right-2 bg-black/75 text-white text-xs px-1.5 py-0.5 rounded">
        {{ formatDuration(video.durationSeconds) }}
      </div>
    </div>
    <div class="p-4">
      <h3 class="font-medium text-gray-900 group-hover:text-primary-600 transition-colors line-clamp-1">
        {{ video.title }}
      </h3>
      <p v-if="video.description" class="mt-1 text-sm text-gray-500 line-clamp-2">
        {{ video.description }}
      </p>
    </div>
  </article>
</template>
```

**Step 3: Create Home page**

Create `web/src/pages/Home.vue`:

```vue
<script setup lang="ts">
import { Link } from '@inertiajs/vue3'
import Container from '@/components/ui/Container.vue'
</script>

<template>
  <div class="min-h-screen">
    <!-- Hero -->
    <div class="bg-gradient-to-br from-primary-600 to-primary-800 text-white py-20">
      <Container>
        <h1 class="text-4xl md:text-5xl font-bold">Gooji</h1>
        <p class="mt-4 text-xl text-primary-100 max-w-2xl">
          Preserving Ojibwe language and culture through community video sharing.
        </p>
        <div class="mt-8 flex gap-4">
          <Link href="/gallery" class="btn btn-primary bg-white text-primary-700 hover:bg-primary-50">
            Browse Videos
          </Link>
          <Link href="/upload" class="btn border border-white/30 hover:bg-white/10">
            Share Your Story
          </Link>
        </div>
      </Container>
    </div>

    <!-- Features -->
    <Container class="py-16">
      <h2 class="text-2xl font-bold text-gray-900">Community Video Sharing</h2>
      <p class="mt-2 text-gray-600 max-w-2xl">
        Record or upload videos to share Ojibwe language, stories, and cultural knowledge with the community.
      </p>
    </Container>
  </div>
</template>
```

**Step 4: Create Gallery page**

Create `web/src/pages/Gallery.vue`:

```vue
<script setup lang="ts">
import type { Video } from '@/types'
import Container from '@/components/ui/Container.vue'
import VideoCard from '@/components/VideoCard.vue'

defineProps<{
  videos: Video[]
}>()
</script>

<template>
  <div class="min-h-screen bg-gray-50 py-8">
    <Container>
      <header class="mb-8">
        <h1 class="text-3xl font-bold text-gray-900">Video Gallery</h1>
        <p class="mt-2 text-gray-600">
          Browse videos shared by the community.
        </p>
      </header>

      <!-- Empty state -->
      <div
        v-if="videos.length === 0"
        class="text-center py-16 bg-white rounded-lg border-2 border-dashed border-gray-300"
      >
        <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" />
        </svg>
        <h3 class="mt-4 text-lg font-medium text-gray-900">No videos yet</h3>
        <p class="mt-2 text-gray-500">Be the first to share a video with the community.</p>
      </div>

      <!-- Video grid -->
      <div
        v-else
        class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6"
      >
        <VideoCard
          v-for="video in videos"
          :key="video.id"
          :video="video"
        />
      </div>
    </Container>
  </div>
</template>
```

**Step 5: Add vue shim for TypeScript**

Create `web/src/shims-vue.d.ts`:

```typescript
declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent<{}, {}, any>
  export default component
}
```

**Step 6: Verify frontend builds**

```bash
cd web && npm run build && cd ..
```

Expected: Build completes, dist/ folder created with manifest.json.

**Step 7: Commit**

```bash
git add web/src/
git commit -m "feat: add Gallery and Home pages with empty state

- Home page with hero section and navigation
- Gallery page with video grid and empty state
- VideoCard component for video display
- Container utility component"
```

---

### Task 9: Health Check Endpoint

Already implemented in Task 6. Verify it works.

**Step 1: Start the server and test**

```bash
go run ./cmd/server &
sleep 2
curl http://localhost:8080/health
pkill -f "go run"
```

Expected: `{"status":"ok"}`

**Step 2: No commit needed (already done)**

---

### Task 10: Integration Test

**Files:**
- Create: `tests/integration/gallery_test.go`

**Step 1: Write integration test**

Create `tests/integration/gallery_test.go`:

```go
package integration_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"gooji/internal/domain"
	"gooji/internal/handler"
	"gooji/internal/repository/sqlite"
	"gooji/internal/service"
)

func setupTestServer(t *testing.T) (*httptest.Server, *service.VideoService) {
	db, err := sqlite.NewDB(":memory:")
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	videoRepo := sqlite.NewVideoRepository(db)
	videoService := service.NewVideoService(videoRepo)
	h := handler.New(videoService, "")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", h.Health)

	return httptest.NewServer(mux), videoService
}

func TestHealthEndpoint(t *testing.T) {
	server, _ := setupTestServer(t)
	defer server.Close()

	resp, err := http.Get(server.URL + "/health")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var result map[string]string
	json.Unmarshal(body, &result)

	if result["status"] != "ok" {
		t.Errorf("expected status ok, got %s", result["status"])
	}
}

func TestVideoServiceIntegration(t *testing.T) {
	_, videoService := setupTestServer(t)
	ctx := t.Context()

	// Create a video
	video := domain.NewVideo("Integration Test", "/test.mp4", "video/mp4")
	err := videoService.Create(ctx, video)
	if err != nil {
		t.Fatalf("failed to create: %v", err)
	}

	// List pending (should have 1)
	pending, _ := videoService.ListPending(ctx)
	if len(pending) != 1 {
		t.Errorf("expected 1 pending, got %d", len(pending))
	}

	// Approve it
	err = videoService.Approve(ctx, video.ID)
	if err != nil {
		t.Fatalf("failed to approve: %v", err)
	}

	// List published (should have 1)
	published, _ := videoService.ListPublished(ctx)
	if len(published) != 1 {
		t.Errorf("expected 1 published, got %d", len(published))
	}

	// List pending (should have 0)
	pending, _ = videoService.ListPending(ctx)
	if len(pending) != 0 {
		t.Errorf("expected 0 pending, got %d", len(pending))
	}
}
```

**Step 2: Run integration tests**

```bash
go test ./tests/integration/... -v
```

Expected: PASS (2 tests)

**Step 3: Run all tests**

```bash
go test ./... -v
```

Expected: All tests pass.

**Step 4: Commit**

```bash
git add tests/
git commit -m "feat: add integration tests for Phase 1

- Health endpoint test
- Video service integration test (create → approve → verify)"
```

---

## Phase 1 Complete

**Final verification:**

```bash
# Run all tests
go test ./... -v

# Build backend
go build ./cmd/server
go build ./cmd/admin

# Build frontend
cd web && npm run build && cd ..

# Start server (requires Vite dev server for development)
./server
```

**What's working:**
- Go 1.25 backend with clean architecture
- SQLite database with migrations
- Domain entities with business rules
- Video repository and service
- Inertia.js handler setup
- Vue 3 + Vite frontend with Gallery page
- Health check endpoint
- All tests passing

**Next Phase:** Core Video Flow (upload, thumbnail generation, playback)

# Vue + Gonertia Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Replace Go templates with Vue 3 + Gonertia + Vite, migrate to SQLite, implement shadcn-vue UI with medicine wheel colors.

**Architecture:** Clean architecture with domain → service → repository layers. Gonertia handles server-side rendering with Vue 3 SPA feel. SQLite for persistence.

**Tech Stack:** Go 1.25, Gonertia, SQLite (modernc.org/sqlite), Vue 3, Vite, shadcn-vue, Tailwind CSS v4, TypeScript

---

## Phase 1: Frontend Foundation

### Task 1: Initialize Vue + Vite Project

**Files:**
- Create: `web/package.json`
- Create: `web/vite.config.ts`
- Create: `web/tsconfig.json`
- Create: `web/index.html`

**Step 1: Create web directory and initialize npm**

```bash
cd /home/jones/dev/gooji
rm -rf web/templates web/static
mkdir -p web/src
cd web
npm init -y
```

**Step 2: Install core dependencies**

```bash
cd /home/jones/dev/gooji/web
npm install vue@latest @inertiajs/vue3@latest
npm install -D vite@latest @vitejs/plugin-vue@latest typescript vue-tsc
```

**Step 3: Create vite.config.ts**

```typescript
// web/vite.config.ts
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
      input: resolve(__dirname, 'src/main.ts'),
    },
  },
  server: {
    origin: 'http://localhost:5173',
  },
})
```

**Step 4: Create tsconfig.json**

```json
// web/tsconfig.json
{
  "compilerOptions": {
    "target": "ES2020",
    "useDefineForClassFields": true,
    "module": "ESNext",
    "lib": ["ES2020", "DOM", "DOM.Iterable"],
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

**Step 5: Create tsconfig.node.json**

```json
// web/tsconfig.node.json
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

**Step 6: Create index.html**

```html
<!-- web/index.html -->
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Gooji</title>
  </head>
  <body>
    <div id="app"></div>
    <script type="module" src="/src/main.ts"></script>
  </body>
</html>
```

**Step 7: Commit**

```bash
git add web/
git commit -m "feat: initialize Vue 3 + Vite project structure"
```

---

### Task 2: Set Up Inertia.js Client

**Files:**
- Create: `web/src/main.ts`
- Create: `web/src/vite-env.d.ts`

**Step 1: Create vite-env.d.ts**

```typescript
// web/src/vite-env.d.ts
/// <reference types="vite/client" />

declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent<{}, {}, any>
  export default component
}
```

**Step 2: Create main.ts with Inertia setup**

```typescript
// web/src/main.ts
import { createApp, h } from 'vue'
import { createInertiaApp } from '@inertiajs/vue3'

createInertiaApp({
  resolve: (name) => {
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

**Step 3: Create a test page**

```vue
<!-- web/src/pages/Home.vue -->
<script setup lang="ts">
defineProps<{
  message: string
}>()
</script>

<template>
  <div>
    <h1>Gooji</h1>
    <p>{{ message }}</p>
  </div>
</template>
```

**Step 4: Commit**

```bash
git add web/src/
git commit -m "feat: set up Inertia.js client with Vue 3"
```

---

### Task 3: Install and Configure Tailwind CSS v4

**Files:**
- Create: `web/src/assets/index.css`
- Modify: `web/package.json`
- Modify: `web/vite.config.ts`

**Step 1: Install Tailwind CSS v4**

```bash
cd /home/jones/dev/gooji/web
npm install -D tailwindcss@latest @tailwindcss/vite@latest
```

**Step 2: Update vite.config.ts**

```typescript
// web/vite.config.ts
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite'
import { resolve } from 'path'

export default defineConfig({
  plugins: [vue(), tailwindcss()],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src'),
    },
  },
  build: {
    manifest: true,
    outDir: 'dist',
    rollupOptions: {
      input: resolve(__dirname, 'src/main.ts'),
    },
  },
  server: {
    origin: 'http://localhost:5173',
  },
})
```

**Step 3: Create index.css with medicine wheel theme**

```css
/* web/src/assets/index.css */
@import "tailwindcss";

@theme {
  /* Medicine wheel colors - used subtly */
  --color-earth-black: oklch(0.15 0.01 260);
  --color-snow-white: oklch(0.97 0.01 90);
  --color-fire-red: oklch(0.55 0.18 25);
  --color-sun-yellow: oklch(0.80 0.14 85);

  /* Map to semantic colors */
  --color-primary: var(--color-fire-red);
  --color-accent: var(--color-sun-yellow);
  --color-background: var(--color-snow-white);
  --color-foreground: var(--color-earth-black);
}
```

**Step 4: Import CSS in main.ts**

```typescript
// web/src/main.ts
import { createApp, h } from 'vue'
import { createInertiaApp } from '@inertiajs/vue3'
import './assets/index.css'

createInertiaApp({
  resolve: (name) => {
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

**Step 5: Commit**

```bash
git add web/
git commit -m "feat: add Tailwind CSS v4 with medicine wheel color theme"
```

---

### Task 4: Install and Configure shadcn-vue

**Files:**
- Create: `web/components.json`
- Create: `web/src/lib/utils.ts`
- Create: `web/src/components/ui/button.vue`

**Step 1: Install shadcn-vue dependencies**

```bash
cd /home/jones/dev/gooji/web
npm install -D radix-vue class-variance-authority clsx tailwind-merge
npm install lucide-vue-next
```

**Step 2: Create utils.ts**

```typescript
// web/src/lib/utils.ts
import { type ClassValue, clsx } from 'clsx'
import { twMerge } from 'tailwind-merge'

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}
```

**Step 3: Create components.json**

```json
// web/components.json
{
  "$schema": "https://shadcn-vue.com/schema.json",
  "style": "default",
  "typescript": true,
  "tailwind": {
    "config": "",
    "css": "src/assets/index.css",
    "baseColor": "neutral"
  },
  "aliases": {
    "components": "@/components",
    "utils": "@/lib/utils"
  }
}
```

**Step 4: Create Button component**

```vue
<!-- web/src/components/ui/Button.vue -->
<script setup lang="ts">
import { computed } from 'vue'
import { cn } from '@/lib/utils'

interface Props {
  variant?: 'default' | 'outline' | 'ghost' | 'destructive'
  size?: 'default' | 'sm' | 'lg' | 'icon'
  disabled?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'default',
  size: 'default',
  disabled: false,
})

const classes = computed(() =>
  cn(
    'inline-flex items-center justify-center rounded-md font-medium transition-colors',
    'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2',
    'disabled:pointer-events-none disabled:opacity-50',
    {
      'bg-primary text-snow-white hover:bg-primary/90': props.variant === 'default',
      'border border-foreground/20 bg-transparent hover:bg-foreground/5': props.variant === 'outline',
      'hover:bg-foreground/5': props.variant === 'ghost',
      'bg-red-600 text-white hover:bg-red-700': props.variant === 'destructive',
    },
    {
      'h-10 px-4 py-2': props.size === 'default',
      'h-8 px-3 text-sm': props.size === 'sm',
      'h-12 px-6 text-lg': props.size === 'lg',
      'h-10 w-10': props.size === 'icon',
    }
  )
)
</script>

<template>
  <button :class="classes" :disabled="disabled">
    <slot />
  </button>
</template>
```

**Step 5: Commit**

```bash
git add web/
git commit -m "feat: add shadcn-vue foundation with Button component"
```

---

## Phase 2: Go Backend with Gonertia

### Task 5: Add Gonertia and SQLite Dependencies

**Files:**
- Modify: `go.mod`

**Step 1: Add Go dependencies**

```bash
cd /home/jones/dev/gooji
go get github.com/romsar/gonertia
go get modernc.org/sqlite
go get github.com/golang-migrate/migrate/v4
go get github.com/google/uuid
go get github.com/go-playground/validator/v10
go mod tidy
```

**Step 2: Commit**

```bash
git add go.mod go.sum
git commit -m "feat: add Gonertia, SQLite, and validation dependencies"
```

---

### Task 6: Create Domain Types

**Files:**
- Create: `internal/domain/video.go`
- Create: `internal/domain/errors.go`

**Step 1: Create video domain types**

```go
// internal/domain/video.go
package domain

import "time"

type Video struct {
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	Language        string    `json:"language"`
	DurationSeconds float64   `json:"duration_seconds"`
	FilePath        string    `json:"-"`
	ThumbnailPath   string    `json:"-"`
	FileSizeBytes   int64     `json:"file_size_bytes"`
	MimeType        string    `json:"mime_type"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type Tag struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type VideoListOptions struct {
	Limit  int
	Offset int
}
```

**Step 2: Create domain errors**

```go
// internal/domain/errors.go
package domain

import "errors"

var (
	ErrVideoNotFound   = errors.New("video not found")
	ErrInvalidVideo    = errors.New("invalid video")
	ErrInvalidInput    = errors.New("invalid input")
	ErrStorageFailed   = errors.New("storage operation failed")
	ErrProcessingFailed = errors.New("video processing failed")
)
```

**Step 3: Commit**

```bash
git add internal/domain/
git commit -m "feat: add domain types for Video and Tag"
```

---

### Task 7: Create SQLite Migrations

**Files:**
- Create: `migrations/001_initial.up.sql`
- Create: `migrations/001_initial.down.sql`

**Step 1: Create up migration**

```sql
-- migrations/001_initial.up.sql
CREATE TABLE IF NOT EXISTS videos (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT DEFAULT '',
    language TEXT DEFAULT 'oji',
    duration_seconds REAL,
    file_path TEXT NOT NULL,
    thumbnail_path TEXT,
    file_size_bytes INTEGER NOT NULL,
    mime_type TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS video_tags (
    video_id TEXT REFERENCES videos(id) ON DELETE CASCADE,
    tag_id INTEGER REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (video_id, tag_id)
);

CREATE INDEX IF NOT EXISTS idx_videos_created ON videos(created_at DESC);
```

**Step 2: Create down migration**

```sql
-- migrations/001_initial.down.sql
DROP INDEX IF EXISTS idx_videos_created;
DROP TABLE IF EXISTS video_tags;
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS videos;
```

**Step 3: Commit**

```bash
git add migrations/
git commit -m "feat: add SQLite schema migrations"
```

---

### Task 8: Create SQLite Repository

**Files:**
- Create: `internal/repository/interfaces.go`
- Create: `internal/repository/sqlite/video.go`
- Create: `internal/repository/sqlite/db.go`

**Step 1: Create repository interfaces**

```go
// internal/repository/interfaces.go
package repository

import (
	"context"
	"gooji/internal/domain"
	"io"
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
```

**Step 2: Create database connection**

```go
// internal/repository/sqlite/db.go
package sqlite

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func NewDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	return db, nil
}

func RunMigrations(db *sql.DB, migrationsPath string) error {
	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	source, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("failed to create migration source: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", source, "sqlite", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
```

**Step 3: Create video repository**

```go
// internal/repository/sqlite/video.go
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
```

**Step 4: Commit**

```bash
git add internal/repository/
git commit -m "feat: add SQLite video repository implementation"
```

---

### Task 9: Create Video Service

**Files:**
- Create: `internal/service/video.go`

**Step 1: Create video service**

```go
// internal/service/video.go
package service

import (
	"context"
	"fmt"
	"html"
	"io"
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
```

**Step 2: Commit**

```bash
git add internal/service/
git commit -m "feat: add video service with upload, list, delete operations"
```

---

### Task 10: Create Gonertia Handlers

**Files:**
- Create: `internal/handler/inertia.go`
- Create: `internal/handler/video.go`
- Create: `internal/handler/pages.go`

**Step 1: Create Inertia helper**

```go
// internal/handler/inertia.go
package handler

import (
	"html/template"
	"io/fs"
	"net/http"
	"os"

	"github.com/romsar/gonertia"
)

func NewInertia(isDev bool) (*gonertia.Inertia, error) {
	var opts []gonertia.Option

	if isDev {
		// Development: use Vite dev server
		opts = append(opts, gonertia.WithTemplateFunc(func(w http.ResponseWriter, r *http.Request, page *gonertia.Page) error {
			tmpl := template.Must(template.New("app").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Gooji</title>
    <script type="module" src="http://localhost:5173/@vite/client"></script>
    <script type="module" src="http://localhost:5173/src/main.ts"></script>
</head>
<body>
    <div id="app" data-page='{{ .PageJSON }}'></div>
</body>
</html>
`))
			pageJSON, _ := page.MarshalJSON()
			return tmpl.Execute(w, map[string]any{"PageJSON": template.JS(pageJSON)})
		}))
	} else {
		// Production: use built assets
		distFS := os.DirFS("web/dist")
		opts = append(opts, gonertia.WithFS(distFS.(fs.FS)))
	}

	return gonertia.New("/", opts...)
}
```

**Step 2: Create video handlers**

```go
// internal/handler/video.go
package handler

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/romsar/gonertia"
	"go.uber.org/zap"

	"gooji/internal/domain"
	"gooji/internal/service"
)

type VideoHandler struct {
	inertia *gonertia.Inertia
	service *service.VideoService
	logger  *zap.Logger
	storage string
}

func NewVideoHandler(i *gonertia.Inertia, s *service.VideoService, l *zap.Logger, storagePath string) *VideoHandler {
	return &VideoHandler{
		inertia: i,
		service: s,
		logger:  l,
		storage: storagePath,
	}
}

func (h *VideoHandler) Upload(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(100 << 20); err != nil {
		h.jsonError(w, "failed to parse form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("video")
	if err != nil {
		h.jsonError(w, "missing video file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	req := service.UploadRequest{
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
	}

	video, err := h.service.Upload(r.Context(), file, header, req)
	if err != nil {
		h.logger.Error("upload failed", zap.Error(err))
		h.jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.jsonResponse(w, video, http.StatusCreated)
}

func (h *VideoHandler) List(w http.ResponseWriter, r *http.Request) {
	videos, err := h.service.List(r.Context(), domain.VideoListOptions{Limit: 50})
	if err != nil {
		h.logger.Error("list failed", zap.Error(err))
		h.jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.jsonResponse(w, videos, http.StatusOK)
}

func (h *VideoHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := h.extractID(r.URL.Path, "/api/videos/")
	if id == "" {
		h.jsonError(w, "missing video ID", http.StatusBadRequest)
		return
	}

	video, err := h.service.Get(r.Context(), id)
	if err != nil {
		if err == domain.ErrVideoNotFound {
			h.jsonError(w, "video not found", http.StatusNotFound)
			return
		}
		h.jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.ServeFile(w, r, video.FilePath)
}

func (h *VideoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := h.extractID(r.URL.Path, "/api/videos/")
	if id == "" {
		h.jsonError(w, "missing video ID", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		if err == domain.ErrVideoNotFound {
			h.jsonError(w, "video not found", http.StatusNotFound)
			return
		}
		h.jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.jsonResponse(w, map[string]string{"status": "deleted"}, http.StatusOK)
}

func (h *VideoHandler) Thumbnail(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	// Remove extension if present and add .jpg
	base := strings.TrimSuffix(id, filepath.Ext(id))
	thumbnailPath := filepath.Join(h.storage, "thumbnails", base+".jpg")

	http.ServeFile(w, r, thumbnailPath)
}

func (h *VideoHandler) extractID(path, prefix string) string {
	return strings.TrimPrefix(path, prefix)
}

func (h *VideoHandler) jsonResponse(w http.ResponseWriter, data any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *VideoHandler) jsonError(w http.ResponseWriter, msg string, status int) {
	h.jsonResponse(w, map[string]string{"error": msg}, status)
}
```

**Step 3: Create page handlers**

```go
// internal/handler/pages.go
package handler

import (
	"net/http"

	"github.com/romsar/gonertia"
	"go.uber.org/zap"

	"gooji/internal/domain"
	"gooji/internal/service"
)

type PageHandler struct {
	inertia *gonertia.Inertia
	service *service.VideoService
	logger  *zap.Logger
}

func NewPageHandler(i *gonertia.Inertia, s *service.VideoService, l *zap.Logger) *PageHandler {
	return &PageHandler{
		inertia: i,
		service: s,
		logger:  l,
	}
}

func (h *PageHandler) Home(w http.ResponseWriter, r *http.Request) {
	h.inertia.Render(w, r, "Home", gonertia.Props{
		"message": "Welcome to Gooji",
	})
}

func (h *PageHandler) Gallery(w http.ResponseWriter, r *http.Request) {
	videos, err := h.service.List(r.Context(), domain.VideoListOptions{Limit: 50})
	if err != nil {
		h.logger.Error("failed to list videos", zap.Error(err))
		videos = []domain.Video{}
	}

	h.inertia.Render(w, r, "Gallery", gonertia.Props{
		"videos": videos,
	})
}

func (h *PageHandler) Upload(w http.ResponseWriter, r *http.Request) {
	h.inertia.Render(w, r, "Upload", gonertia.Props{})
}

func (h *PageHandler) Record(w http.ResponseWriter, r *http.Request) {
	h.inertia.Render(w, r, "Record", gonertia.Props{})
}

func (h *PageHandler) VideoShow(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/video/"):]
	video, err := h.service.Get(r.Context(), id)
	if err != nil {
		h.inertia.Render(w, r, "Error", gonertia.Props{
			"message": "Video not found",
		})
		return
	}

	h.inertia.Render(w, r, "Video/Show", gonertia.Props{
		"video": video,
	})
}

func (h *PageHandler) Edit(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/edit/"):]
	video, err := h.service.Get(r.Context(), id)
	if err != nil {
		h.inertia.Render(w, r, "Error", gonertia.Props{
			"message": "Video not found",
		})
		return
	}

	h.inertia.Render(w, r, "Edit", gonertia.Props{
		"video": video,
	})
}
```

**Step 4: Commit**

```bash
git add internal/handler/
git commit -m "feat: add Gonertia handlers for pages and API"
```

---

## Phase 3: Wire Everything Together

### Task 11: Update Main Entry Point

**Files:**
- Modify: `main.go`

**Step 1: Rewrite main.go**

```go
// main.go
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"gooji/internal/config"
	"gooji/internal/handler"
	"gooji/internal/middleware"
	"gooji/internal/repository"
	"gooji/internal/repository/sqlite"
	"gooji/internal/service"
	"gooji/pkg/ffmpeg"
)

func main() {
	// Initialize logger
	logger, _ := zap.NewProduction()
	if os.Getenv("APP_ENV") == "development" {
		logger, _ = zap.NewDevelopment()
	}
	defer logger.Sync()

	// Load configuration
	cfg, err := config.Load("config/config.json")
	if err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
	}

	// Initialize database
	db, err := sqlite.NewDB("storage/gooji.db")
	if err != nil {
		logger.Fatal("failed to open database", zap.Error(err))
	}
	defer db.Close()

	if err := sqlite.RunMigrations(db, "migrations"); err != nil {
		logger.Fatal("failed to run migrations", zap.Error(err))
	}

	// Initialize repositories
	videoRepo := sqlite.NewVideoRepository(db)
	fileStorage := repository.NewFileStorage(cfg.Storage.Uploads, cfg.Storage.Thumbnails)

	// Initialize FFmpeg processor
	processor := ffmpeg.NewProcessorWithSecurity(cfg.FFmpeg.Path, cfg.Storage.BasePath)

	// Initialize services
	videoService := service.NewVideoService(videoRepo, fileStorage, processor, logger)

	// Initialize Inertia
	isDev := os.Getenv("APP_ENV") == "development" || os.Getenv("VITE_DEV") == "true"
	inertia, err := handler.NewInertia(isDev)
	if err != nil {
		logger.Fatal("failed to create inertia", zap.Error(err))
	}

	// Initialize handlers
	pageHandler := handler.NewPageHandler(inertia, videoService, logger)
	videoHandler := handler.NewVideoHandler(inertia, videoService, logger, cfg.Storage.BasePath)

	// Set up routes
	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("POST /api/videos", videoHandler.Upload)
	mux.HandleFunc("GET /api/videos", videoHandler.List)
	mux.HandleFunc("GET /api/videos/{id}", videoHandler.Get)
	mux.HandleFunc("DELETE /api/videos/{id}", videoHandler.Delete)
	mux.HandleFunc("GET /api/thumbnails", videoHandler.Thumbnail)

	// Page routes
	mux.HandleFunc("GET /", pageHandler.Home)
	mux.HandleFunc("GET /gallery", pageHandler.Gallery)
	mux.HandleFunc("GET /upload", pageHandler.Upload)
	mux.HandleFunc("GET /record", pageHandler.Record)
	mux.HandleFunc("GET /video/{id}", pageHandler.VideoShow)
	mux.HandleFunc("GET /edit/{id}", pageHandler.Edit)

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	})

	// Middleware
	var h http.Handler = mux
	h = middleware.Recovery(logger)(h)
	h = middleware.Logging(logger)(h)
	h = middleware.CORS()(h)
	h = inertia.Middleware(h)

	// Start server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: h,
	}

	go func() {
		logger.Info("starting server", zap.Int("port", cfg.Server.Port))
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			logger.Fatal("server error", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.Info("shutting down server")
	server.Shutdown(ctx)
}
```

**Step 2: Commit**

```bash
git add main.go
git commit -m "feat: rewrite main.go with Gonertia and SQLite"
```

---

### Task 12: Update Taskfile

**Files:**
- Modify: `Taskfile.yml`

**Step 1: Update Taskfile.yml**

```yaml
# Taskfile.yml
version: '3'

vars:
  BINARY_NAME: gooji

tasks:
  default:
    cmds:
      - task --list

  dev:
    desc: Start Go server and Vite dev server concurrently
    deps: [dev:frontend]
    cmds:
      - APP_ENV=development VITE_DEV=true go run github.com/air-verse/air@latest

  dev:frontend:
    desc: Run Vite dev server
    dir: web
    cmds:
      - npm run dev
    ignore_error: true

  build:
    desc: Build production binary and frontend
    cmds:
      - task: build:frontend
      - task: build:backend

  build:frontend:
    desc: Build frontend for production
    dir: web
    cmds:
      - npm run build

  build:backend:
    desc: Build Go binary
    cmds:
      - go build -o bin/{{.BINARY_NAME}} main.go

  test:
    desc: Run all tests
    cmds:
      - go test -v -race -cover ./internal/...

  test:coverage:
    desc: Generate coverage report
    cmds:
      - go test -coverprofile=coverage.out ./internal/...
      - go tool cover -html=coverage.out -o coverage.html

  lint:
    desc: Run linters
    cmds:
      - go vet ./...
      - go run honnef.co/go/tools/cmd/staticcheck@latest ./...

  migrate:data:
    desc: Migrate JSON metadata to SQLite (one-time)
    cmds:
      - go run cmd/migrate/main.go

  clean:
    desc: Clean build artifacts
    cmds:
      - rm -rf bin/ web/dist/ coverage.out coverage.html
```

**Step 2: Update web/package.json scripts**

```json
{
  "scripts": {
    "dev": "vite",
    "build": "vue-tsc && vite build",
    "preview": "vite preview",
    "type-check": "vue-tsc --noEmit"
  }
}
```

**Step 3: Commit**

```bash
git add Taskfile.yml web/package.json
git commit -m "feat: update Taskfile for Vue + Gonertia dev workflow"
```

---

## Phase 4: Vue Pages

### Task 13: Create Layout Component

**Files:**
- Create: `web/src/components/Layout.vue`

**Step 1: Create Layout.vue**

```vue
<!-- web/src/components/Layout.vue -->
<script setup lang="ts">
import { Link } from '@inertiajs/vue3'
</script>

<template>
  <div class="min-h-screen bg-background">
    <nav class="border-b border-foreground/10 bg-snow-white">
      <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div class="flex h-16 items-center justify-between">
          <Link href="/" class="text-2xl font-bold text-foreground">
            Gooji
          </Link>
          <div class="flex gap-6">
            <Link
              href="/gallery"
              class="text-foreground/70 hover:text-foreground border-b-2 border-transparent hover:border-accent transition-colors"
            >
              Gallery
            </Link>
            <Link
              href="/record"
              class="text-foreground/70 hover:text-foreground border-b-2 border-transparent hover:border-accent transition-colors"
            >
              Record
            </Link>
            <Link
              href="/upload"
              class="text-foreground/70 hover:text-foreground border-b-2 border-transparent hover:border-accent transition-colors"
            >
              Upload
            </Link>
          </div>
        </div>
      </div>
    </nav>

    <main class="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
      <slot />
    </main>

    <footer class="mt-auto border-t border-foreground/10 bg-earth-black text-snow-white">
      <div class="mx-auto max-w-7xl px-4 py-6 sm:px-6 lg:px-8">
        <p class="text-center text-sm">
          Gooji - Preserving Ojibwe language and culture through video
        </p>
      </div>
    </footer>
  </div>
</template>
```

**Step 2: Commit**

```bash
git add web/src/components/Layout.vue
git commit -m "feat: add Layout component with navigation"
```

---

### Task 14: Create All Vue Pages

**Files:**
- Modify: `web/src/pages/Home.vue`
- Create: `web/src/pages/Gallery.vue`
- Create: `web/src/pages/Upload.vue`
- Create: `web/src/pages/Record.vue`
- Create: `web/src/pages/Video/Show.vue`
- Create: `web/src/pages/Edit.vue`
- Create: `web/src/pages/Error.vue`

**Step 1: Create Home.vue**

```vue
<!-- web/src/pages/Home.vue -->
<script setup lang="ts">
import Layout from '@/components/Layout.vue'
import Button from '@/components/ui/Button.vue'
import { Link } from '@inertiajs/vue3'

defineProps<{
  message: string
}>()
</script>

<template>
  <Layout>
    <div class="text-center py-16">
      <h1 class="text-4xl font-bold text-foreground mb-4">
        {{ message }}
      </h1>
      <p class="text-lg text-foreground/70 mb-8 max-w-2xl mx-auto">
        Record, share, and preserve Ojibwe language and cultural content for future generations.
      </p>
      <div class="flex gap-4 justify-center">
        <Link href="/record">
          <Button variant="default">Start Recording</Button>
        </Link>
        <Link href="/gallery">
          <Button variant="outline">Browse Gallery</Button>
        </Link>
      </div>
    </div>
  </Layout>
</template>
```

**Step 2: Create Gallery.vue**

```vue
<!-- web/src/pages/Gallery.vue -->
<script setup lang="ts">
import Layout from '@/components/Layout.vue'
import { Link } from '@inertiajs/vue3'

interface Video {
  id: string
  title: string
  description: string
  duration_seconds: number
  created_at: string
}

defineProps<{
  videos: Video[]
}>()

function formatDuration(seconds: number): string {
  const mins = Math.floor(seconds / 60)
  const secs = Math.floor(seconds % 60)
  return `${mins}:${secs.toString().padStart(2, '0')}`
}
</script>

<template>
  <Layout>
    <h1 class="text-3xl font-bold text-foreground mb-8">Video Gallery</h1>

    <div v-if="videos.length === 0" class="text-center py-16 text-foreground/60">
      <p>No videos yet. Be the first to contribute!</p>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <Link
        v-for="video in videos"
        :key="video.id"
        :href="`/video/${video.id}`"
        class="group block rounded-lg border border-foreground/10 overflow-hidden hover:border-accent transition-colors"
      >
        <div class="aspect-video bg-foreground/5 relative">
          <img
            :src="`/api/thumbnails?id=${video.id}`"
            :alt="video.title"
            class="w-full h-full object-cover"
          />
          <span class="absolute bottom-2 right-2 bg-earth-black/80 text-snow-white text-xs px-2 py-1 rounded">
            {{ formatDuration(video.duration_seconds) }}
          </span>
        </div>
        <div class="p-4">
          <h3 class="font-medium text-foreground group-hover:text-primary transition-colors">
            {{ video.title || 'Untitled' }}
          </h3>
          <p class="text-sm text-foreground/60 mt-1 line-clamp-2">
            {{ video.description }}
          </p>
        </div>
      </Link>
    </div>
  </Layout>
</template>
```

**Step 3: Create Upload.vue**

```vue
<!-- web/src/pages/Upload.vue -->
<script setup lang="ts">
import { ref } from 'vue'
import { router } from '@inertiajs/vue3'
import Layout from '@/components/Layout.vue'
import Button from '@/components/ui/Button.vue'

const title = ref('')
const description = ref('')
const file = ref<File | null>(null)
const uploading = ref(false)
const error = ref('')
const dragOver = ref(false)

function onFileSelect(event: Event) {
  const input = event.target as HTMLInputElement
  if (input.files?.[0]) {
    file.value = input.files[0]
  }
}

function onDrop(event: DragEvent) {
  dragOver.value = false
  if (event.dataTransfer?.files?.[0]) {
    file.value = event.dataTransfer.files[0]
  }
}

async function submit() {
  if (!file.value) return

  uploading.value = true
  error.value = ''

  const formData = new FormData()
  formData.append('video', file.value)
  formData.append('title', title.value)
  formData.append('description', description.value)

  try {
    const response = await fetch('/api/videos', {
      method: 'POST',
      body: formData,
    })

    if (!response.ok) {
      throw new Error('Upload failed')
    }

    const data = await response.json()
    router.visit(`/video/${data.id}`)
  } catch (e) {
    error.value = 'Failed to upload video. Please try again.'
  } finally {
    uploading.value = false
  }
}
</script>

<template>
  <Layout>
    <div class="max-w-2xl mx-auto">
      <h1 class="text-3xl font-bold text-foreground mb-8">Upload Video</h1>

      <form @submit.prevent="submit" class="space-y-6">
        <div
          class="border-2 border-dashed rounded-lg p-8 text-center transition-colors"
          :class="dragOver ? 'border-primary bg-primary/5' : 'border-foreground/20'"
          @dragover.prevent="dragOver = true"
          @dragleave="dragOver = false"
          @drop.prevent="onDrop"
        >
          <input
            type="file"
            accept="video/mp4,video/webm"
            class="hidden"
            id="video-input"
            @change="onFileSelect"
          />
          <label for="video-input" class="cursor-pointer">
            <div v-if="file" class="text-foreground">
              <p class="font-medium">{{ file.name }}</p>
              <p class="text-sm text-foreground/60">{{ (file.size / 1024 / 1024).toFixed(2) }} MB</p>
            </div>
            <div v-else class="text-foreground/60">
              <p>Drag and drop a video file here, or click to select</p>
              <p class="text-sm mt-2">MP4 or WebM, up to 100MB</p>
            </div>
          </label>
        </div>

        <div>
          <label class="block text-sm font-medium text-foreground mb-2">Title</label>
          <input
            v-model="title"
            type="text"
            class="w-full rounded-md border border-foreground/20 px-3 py-2 focus:border-primary focus:outline-none"
            placeholder="Give your video a title"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-foreground mb-2">Description</label>
          <textarea
            v-model="description"
            rows="3"
            class="w-full rounded-md border border-foreground/20 px-3 py-2 focus:border-primary focus:outline-none"
            placeholder="Describe what's in this video"
          />
        </div>

        <p v-if="error" class="text-red-600 text-sm">{{ error }}</p>

        <Button type="submit" :disabled="!file || uploading" class="w-full">
          {{ uploading ? 'Uploading...' : 'Upload Video' }}
        </Button>
      </form>
    </div>
  </Layout>
</template>
```

**Step 4: Create Record.vue**

```vue
<!-- web/src/pages/Record.vue -->
<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { router } from '@inertiajs/vue3'
import Layout from '@/components/Layout.vue'
import Button from '@/components/ui/Button.vue'

const videoRef = ref<HTMLVideoElement | null>(null)
const stream = ref<MediaStream | null>(null)
const mediaRecorder = ref<MediaRecorder | null>(null)
const chunks = ref<Blob[]>([])
const recording = ref(false)
const recorded = ref<Blob | null>(null)
const error = ref('')

onMounted(async () => {
  try {
    stream.value = await navigator.mediaDevices.getUserMedia({
      video: true,
      audio: true,
    })
    if (videoRef.value) {
      videoRef.value.srcObject = stream.value
    }
  } catch (e) {
    error.value = 'Could not access camera. Please grant permission.'
  }
})

onUnmounted(() => {
  stream.value?.getTracks().forEach((track) => track.stop())
})

function startRecording() {
  if (!stream.value) return

  chunks.value = []
  mediaRecorder.value = new MediaRecorder(stream.value, {
    mimeType: 'video/webm',
  })

  mediaRecorder.value.ondataavailable = (e) => {
    if (e.data.size > 0) {
      chunks.value.push(e.data)
    }
  }

  mediaRecorder.value.onstop = () => {
    recorded.value = new Blob(chunks.value, { type: 'video/webm' })
  }

  mediaRecorder.value.start()
  recording.value = true
}

function stopRecording() {
  mediaRecorder.value?.stop()
  recording.value = false
}

function retake() {
  recorded.value = null
}

async function upload() {
  if (!recorded.value) return

  const formData = new FormData()
  formData.append('video', recorded.value, 'recording.webm')
  formData.append('title', '')
  formData.append('description', '')

  const response = await fetch('/api/videos', {
    method: 'POST',
    body: formData,
  })

  if (response.ok) {
    const data = await response.json()
    router.visit(`/edit/${data.id}`)
  }
}
</script>

<template>
  <Layout>
    <div class="max-w-3xl mx-auto">
      <h1 class="text-3xl font-bold text-foreground mb-8">Record Video</h1>

      <p v-if="error" class="text-red-600 mb-4">{{ error }}</p>

      <div class="relative aspect-video bg-earth-black rounded-lg overflow-hidden mb-6">
        <video
          v-if="!recorded"
          ref="videoRef"
          autoplay
          muted
          playsinline
          class="w-full h-full object-cover"
        />
        <video
          v-else
          :src="URL.createObjectURL(recorded)"
          controls
          class="w-full h-full object-cover"
        />

        <div
          v-if="recording"
          class="absolute top-4 left-4 flex items-center gap-2 bg-primary px-3 py-1 rounded-full"
        >
          <span class="w-2 h-2 bg-snow-white rounded-full animate-pulse" />
          <span class="text-snow-white text-sm font-medium">Recording</span>
        </div>
      </div>

      <div class="flex justify-center gap-4">
        <template v-if="!recorded">
          <Button
            v-if="!recording"
            @click="startRecording"
            variant="default"
            class="bg-primary"
          >
            Start Recording
          </Button>
          <Button
            v-else
            @click="stopRecording"
            variant="destructive"
          >
            Stop Recording
          </Button>
        </template>
        <template v-else>
          <Button @click="retake" variant="outline">
            Retake
          </Button>
          <Button @click="upload" variant="default">
            Use This Video
          </Button>
        </template>
      </div>
    </div>
  </Layout>
</template>
```

**Step 5: Create Video/Show.vue**

```vue
<!-- web/src/pages/Video/Show.vue -->
<script setup lang="ts">
import Layout from '@/components/Layout.vue'
import Button from '@/components/ui/Button.vue'
import { Link } from '@inertiajs/vue3'

interface Video {
  id: string
  title: string
  description: string
  duration_seconds: number
  created_at: string
}

defineProps<{
  video: Video
}>()
</script>

<template>
  <Layout>
    <div class="max-w-4xl mx-auto">
      <div class="aspect-video bg-earth-black rounded-lg overflow-hidden mb-6">
        <video
          :src="`/api/videos/${video.id}`"
          controls
          class="w-full h-full"
        />
      </div>

      <div class="flex items-start justify-between mb-4">
        <div>
          <h1 class="text-2xl font-bold text-foreground">
            {{ video.title || 'Untitled Video' }}
          </h1>
          <p class="text-foreground/60 text-sm mt-1">
            {{ new Date(video.created_at).toLocaleDateString() }}
          </p>
        </div>
        <Link :href="`/edit/${video.id}`">
          <Button variant="outline" size="sm">Edit</Button>
        </Link>
      </div>

      <p v-if="video.description" class="text-foreground/80">
        {{ video.description }}
      </p>
    </div>
  </Layout>
</template>
```

**Step 6: Create Edit.vue**

```vue
<!-- web/src/pages/Edit.vue -->
<script setup lang="ts">
import { ref } from 'vue'
import { router } from '@inertiajs/vue3'
import Layout from '@/components/Layout.vue'
import Button from '@/components/ui/Button.vue'

interface Video {
  id: string
  title: string
  description: string
}

const props = defineProps<{
  video: Video
}>()

const title = ref(props.video.title)
const description = ref(props.video.description)
const saving = ref(false)

async function save() {
  saving.value = true
  // TODO: Implement update API
  router.visit(`/video/${props.video.id}`)
}

async function deleteVideo() {
  if (!confirm('Are you sure you want to delete this video?')) return

  await fetch(`/api/videos/${props.video.id}`, { method: 'DELETE' })
  router.visit('/gallery')
}
</script>

<template>
  <Layout>
    <div class="max-w-2xl mx-auto">
      <h1 class="text-3xl font-bold text-foreground mb-8">Edit Video</h1>

      <form @submit.prevent="save" class="space-y-6">
        <div>
          <label class="block text-sm font-medium text-foreground mb-2">Title</label>
          <input
            v-model="title"
            type="text"
            class="w-full rounded-md border border-foreground/20 px-3 py-2 focus:border-primary focus:outline-none"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-foreground mb-2">Description</label>
          <textarea
            v-model="description"
            rows="4"
            class="w-full rounded-md border border-foreground/20 px-3 py-2 focus:border-primary focus:outline-none"
          />
        </div>

        <div class="flex gap-4">
          <Button type="submit" :disabled="saving" class="flex-1">
            {{ saving ? 'Saving...' : 'Save Changes' }}
          </Button>
          <Button type="button" variant="destructive" @click="deleteVideo">
            Delete
          </Button>
        </div>
      </form>
    </div>
  </Layout>
</template>
```

**Step 7: Create Error.vue**

```vue
<!-- web/src/pages/Error.vue -->
<script setup lang="ts">
import Layout from '@/components/Layout.vue'
import Button from '@/components/ui/Button.vue'
import { Link } from '@inertiajs/vue3'

defineProps<{
  message: string
}>()
</script>

<template>
  <Layout>
    <div class="text-center py-16">
      <h1 class="text-4xl font-bold text-foreground mb-4">Oops!</h1>
      <p class="text-lg text-foreground/70 mb-8">{{ message }}</p>
      <Link href="/">
        <Button>Go Home</Button>
      </Link>
    </div>
  </Layout>
</template>
```

**Step 8: Commit**

```bash
git add web/src/pages/
git commit -m "feat: add all Vue pages with Inertia integration"
```

---

## Phase 5: Cleanup and Tests

### Task 15: Delete Old Template Code

**Files:**
- Delete: `web/templates/` (entire directory)
- Delete: `web/static/` (entire directory)
- Delete: `internal/video/` (entire directory - replaced by new structure)
- Delete: `internal/ui/` (if exists)

**Step 1: Remove old directories**

```bash
cd /home/jones/dev/gooji
rm -rf web/templates web/static
rm -rf internal/video internal/ui
```

**Step 2: Commit**

```bash
git add -A
git commit -m "chore: remove old template and handler code"
```

---

### Task 16: Add File Storage Implementation

**Files:**
- Create: `internal/repository/filesystem/storage.go`

**Step 1: Create file storage**

```go
// internal/repository/filesystem/storage.go
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
```

**Step 2: Commit**

```bash
git add internal/repository/filesystem/
git commit -m "feat: add filesystem storage implementation"
```

---

### Task 17: Write Service Tests

**Files:**
- Create: `internal/service/video_test.go`

**Step 1: Create video service tests**

```go
// internal/service/video_test.go
package service_test

import (
	"context"
	"testing"

	"go.uber.org/zap"

	"gooji/internal/domain"
	"gooji/internal/service"
)

type mockVideoRepo struct {
	videos map[string]*domain.Video
}

func newMockVideoRepo() *mockVideoRepo {
	return &mockVideoRepo{videos: make(map[string]*domain.Video)}
}

func (m *mockVideoRepo) Create(ctx context.Context, video *domain.Video) error {
	m.videos[video.ID] = video
	return nil
}

func (m *mockVideoRepo) GetByID(ctx context.Context, id string) (*domain.Video, error) {
	if v, ok := m.videos[id]; ok {
		return v, nil
	}
	return nil, domain.ErrVideoNotFound
}

func (m *mockVideoRepo) List(ctx context.Context, opts domain.VideoListOptions) ([]domain.Video, error) {
	var result []domain.Video
	for _, v := range m.videos {
		result = append(result, *v)
	}
	return result, nil
}

func (m *mockVideoRepo) Delete(ctx context.Context, id string) error {
	if _, ok := m.videos[id]; !ok {
		return domain.ErrVideoNotFound
	}
	delete(m.videos, id)
	return nil
}

func (m *mockVideoRepo) Update(ctx context.Context, video *domain.Video) error {
	if _, ok := m.videos[video.ID]; !ok {
		return domain.ErrVideoNotFound
	}
	m.videos[video.ID] = video
	return nil
}

func TestVideoService_List(t *testing.T) {
	repo := newMockVideoRepo()
	repo.videos["1"] = &domain.Video{ID: "1", Title: "Test Video"}

	logger, _ := zap.NewDevelopment()
	svc := service.NewVideoService(repo, nil, nil, logger)

	videos, err := svc.List(context.Background(), domain.VideoListOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(videos) != 1 {
		t.Errorf("expected 1 video, got %d", len(videos))
	}
}

func TestVideoService_Get_NotFound(t *testing.T) {
	repo := newMockVideoRepo()
	logger, _ := zap.NewDevelopment()
	svc := service.NewVideoService(repo, nil, nil, logger)

	_, err := svc.Get(context.Background(), "nonexistent")
	if err != domain.ErrVideoNotFound {
		t.Errorf("expected ErrVideoNotFound, got %v", err)
	}
}

func TestVideoService_Delete(t *testing.T) {
	repo := newMockVideoRepo()
	repo.videos["1"] = &domain.Video{ID: "1", Title: "Test", FilePath: "/tmp/test.mp4"}

	logger, _ := zap.NewDevelopment()
	// Note: In real test, mock file storage too
	svc := service.NewVideoService(repo, nil, nil, logger)

	err := svc.Delete(context.Background(), "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(repo.videos) != 0 {
		t.Error("video should be deleted")
	}
}
```

**Step 2: Run tests**

```bash
go test -v ./internal/service/...
```
Expected: All tests pass

**Step 3: Commit**

```bash
git add internal/service/video_test.go
git commit -m "test: add video service unit tests"
```

---

### Task 18: Write Repository Tests

**Files:**
- Create: `internal/repository/sqlite/video_test.go`

**Step 1: Create repository tests**

```go
// internal/repository/sqlite/video_test.go
package sqlite_test

import (
	"context"
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"

	"gooji/internal/domain"
	"gooji/internal/repository/sqlite"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE videos (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			description TEXT DEFAULT '',
			language TEXT DEFAULT 'oji',
			duration_seconds REAL,
			file_path TEXT NOT NULL,
			thumbnail_path TEXT,
			file_size_bytes INTEGER NOT NULL,
			mime_type TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}

	return db
}

func TestVideoRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := sqlite.NewVideoRepository(db)

	video := &domain.Video{
		ID:            "test-id",
		Title:         "Test Video",
		Description:   "A test video",
		FilePath:      "/path/to/video.mp4",
		FileSizeBytes: 1024,
		MimeType:      "video/mp4",
	}

	err := repo.Create(context.Background(), video)
	if err != nil {
		t.Fatalf("failed to create video: %v", err)
	}

	retrieved, err := repo.GetByID(context.Background(), "test-id")
	if err != nil {
		t.Fatalf("failed to get video: %v", err)
	}

	if retrieved.Title != "Test Video" {
		t.Errorf("expected title 'Test Video', got '%s'", retrieved.Title)
	}
}

func TestVideoRepository_List(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := sqlite.NewVideoRepository(db)

	for i := 0; i < 3; i++ {
		video := &domain.Video{
			ID:            fmt.Sprintf("id-%d", i),
			Title:         fmt.Sprintf("Video %d", i),
			FilePath:      "/path/video.mp4",
			FileSizeBytes: 1024,
			MimeType:      "video/mp4",
		}
		repo.Create(context.Background(), video)
	}

	videos, err := repo.List(context.Background(), domain.VideoListOptions{Limit: 10})
	if err != nil {
		t.Fatalf("failed to list videos: %v", err)
	}

	if len(videos) != 3 {
		t.Errorf("expected 3 videos, got %d", len(videos))
	}
}

func TestVideoRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := sqlite.NewVideoRepository(db)

	video := &domain.Video{
		ID:            "delete-me",
		Title:         "To Delete",
		FilePath:      "/path/video.mp4",
		FileSizeBytes: 1024,
		MimeType:      "video/mp4",
	}
	repo.Create(context.Background(), video)

	err := repo.Delete(context.Background(), "delete-me")
	if err != nil {
		t.Fatalf("failed to delete: %v", err)
	}

	_, err = repo.GetByID(context.Background(), "delete-me")
	if err != domain.ErrVideoNotFound {
		t.Error("expected video to be deleted")
	}
}
```

**Step 2: Run tests**

```bash
go test -v ./internal/repository/sqlite/...
```
Expected: All tests pass

**Step 3: Commit**

```bash
git add internal/repository/sqlite/video_test.go
git commit -m "test: add SQLite repository unit tests"
```

---

### Task 19: Update CLAUDE.md

**Files:**
- Modify: `CLAUDE.md`

**Step 1: Update CLAUDE.md with new architecture**

```markdown
# CLAUDE.md

This file provides guidance to Claude Code when working with this repository.

## Project Overview

Gooji (meaning "video" in Ojibwe) is a Go-based video platform for recording and sharing Ojibwe language and cultural content. It uses Vue 3 with Inertia.js for the frontend and SQLite for data storage.

**Tech Stack:**
- **Backend:** Go 1.25, net/http, Gonertia, SQLite (modernc.org/sqlite), Uber Zap logging
- **Frontend:** Vue 3, TypeScript, Vite, shadcn-vue, Tailwind CSS v4
- **Video Processing:** FFmpeg

## Common Commands

```bash
# Development (starts Go server + Vite dev server)
task dev

# Build for production
task build

# Run tests
task test

# Run linters
task lint
```

## Architecture

Clean Architecture with handler → service → repository layers:

```
main.go                          # Entry point, DI wiring
internal/
├── domain/                      # Core entities, no dependencies
│   ├── video.go                 # Video, Tag types
│   └── errors.go                # Domain errors
├── service/                     # Business logic
│   └── video.go                 # Upload, list, delete operations
├── repository/                  # Data access
│   ├── interfaces.go            # Repository interfaces
│   ├── sqlite/                  # SQLite implementations
│   └── filesystem/              # File storage
├── handler/                     # HTTP handlers
│   ├── inertia.go               # Gonertia setup
│   ├── video.go                 # API handlers
│   └── pages.go                 # Page handlers
└── middleware/                  # CORS, logging, recovery
pkg/ffmpeg/                      # FFmpeg wrapper
web/                             # Vue 3 frontend
├── src/
│   ├── pages/                   # Inertia pages
│   ├── components/              # Vue components
│   │   └── ui/                  # shadcn-vue components
│   └── composables/             # Shared logic
└── vite.config.ts
migrations/                      # SQLite migrations
storage/                         # Runtime data (gitignored)
```

## Key Patterns

**Dependency Injection:** All dependencies injected via constructors in main.go. No global state.

**Inertia.js:** Server renders page name + props, Vue handles the view. Pages in `web/src/pages/` map to handler render calls.

**Repository Pattern:** Interfaces in `internal/repository/interfaces.go`, implementations in subpackages.

**Error Handling:** Domain errors in `internal/domain/errors.go`. Wrap with context: `fmt.Errorf("context: %w", err)`

## API Routes

- `GET /` - Home page
- `GET /gallery` - Video gallery
- `GET /video/{id}` - Video player
- `GET /record` - Recording page
- `GET /upload` - Upload page
- `GET /edit/{id}` - Edit video metadata
- `POST /api/videos` - Upload video (multipart)
- `GET /api/videos` - List videos (JSON)
- `GET /api/videos/{id}` - Stream video file
- `DELETE /api/videos/{id}` - Delete video
- `GET /api/thumbnails?id={id}` - Get thumbnail
- `GET /health` - Health check

## Security

- Path validation on all file operations (no `..`, no dangerous chars)
- UUID filenames for uploads
- HTML escaping on user input
- FFmpeg operations restricted to allowed directories

## Color Theme

Medicine wheel colors used subtly:
- **Black (earth-black):** Text, headings
- **White (snow-white):** Backgrounds
- **Red (fire-red):** Primary actions, recording indicator
- **Yellow (sun-yellow):** Accents, hover states
```

**Step 2: Commit**

```bash
git add CLAUDE.md
git commit -m "docs: update CLAUDE.md for Vue + Gonertia architecture"
```

---

### Task 20: Final Verification

**Step 1: Install all dependencies**

```bash
cd /home/jones/dev/gooji
go mod tidy
cd web && npm install
```

**Step 2: Run linters**

```bash
cd /home/jones/dev/gooji
task lint
```
Expected: No errors

**Step 3: Run tests**

```bash
task test
```
Expected: All tests pass

**Step 4: Start dev server**

```bash
task dev
```
Expected: Server starts, frontend builds, app accessible at localhost:8080

**Step 5: Final commit**

```bash
git add -A
git commit -m "chore: final cleanup and verification"
```

---

## Summary

**Total Tasks:** 20
**Estimated Commits:** 20+

**What was built:**
1. Vue 3 + Vite + TypeScript frontend
2. shadcn-vue components with medicine wheel theme
3. Gonertia integration for Inertia.js protocol
4. SQLite database with migrations
5. Clean architecture backend (domain/service/repository)
6. All pages: Home, Gallery, Upload, Record, Video, Edit
7. Unit tests for services and repositories
8. Updated CLAUDE.md documentation

**What was removed:**
- Go HTML templates
- Vanilla JavaScript files
- Old CSS files
- Old video handler/service/repository in `internal/video/`

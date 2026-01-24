# Vue + Gonertia Modernization Design

**Date:** 2026-01-24
**Status:** Approved
**Approach:** In-place replacement

## Overview

Replace Go backend templating with Vue 3 + Gonertia + Vite. Full feature parity with existing pages, SQLite for data storage, shadcn-vue for UI components, medicine wheel color theme.

## Key Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Transition approach | In-place replacement | Clean break, no hybrid state |
| Vue setup | Fresh start | No legacy to preserve |
| Page scope | Full feature parity | All pages migrated at once |
| Testing | Backend unit tests | Focus on service layer |
| Data storage | SQLite | Cleaner than JSON files |
| Moderation | Skip for now | Add in follow-up phase |
| UI components | shadcn-vue | Accessible, customizable |

## Tech Stack

**Backend:**
- Go 1.25 with `net/http` (stdlib router)
- Gonertia for Inertia.js protocol
- SQLite via `modernc.org/sqlite` (pure Go, no CGO)
- Uber Zap for logging
- `go-playground/validator` for input validation

**Frontend:**
- Vue 3 Composition API + TypeScript
- Vite for dev server and builds
- `@inertiajs/vue3` client adapter
- shadcn-vue for components
- Tailwind CSS v4

## Project Structure

```
gooji/
├── cmd/server/main.go           # Entry point
├── internal/
│   ├── domain/                  # Entities, no dependencies
│   ├── service/                 # Business logic
│   ├── repository/              # Data access (SQLite, filesystem)
│   ├── handler/                 # HTTP handlers (Gonertia responses)
│   └── middleware/              # Logging, recovery, CORS
├── pkg/ffmpeg/                  # FFmpeg wrapper
├── web/                         # Vue frontend (Vite project)
│   ├── src/
│   │   ├── pages/               # Inertia pages
│   │   ├── components/          # shadcn-vue + custom
│   │   └── composables/         # Shared logic
│   └── vite.config.ts
├── migrations/                  # SQLite migrations
├── storage/                     # Uploads, thumbnails (gitignored)
└── Taskfile.yml
```

## Backend Architecture

**Clean Architecture layers:**

```
Handler (HTTP) → Service (Business Logic) → Repository (Data Access)
     ↓                    ↓                        ↓
 Gonertia            Domain types              SQLite + FS
```

**Handler pattern with Gonertia:**
```go
func (h *VideoHandler) Gallery(w http.ResponseWriter, r *http.Request) {
    videos, err := h.service.ListPublished(r.Context())
    if err != nil {
        h.inertia.Render(w, r, "Error", map[string]any{"message": "Failed to load videos"})
        return
    }

    h.inertia.Render(w, r, "Gallery", map[string]any{
        "videos": videos,
    })
}
```

**Service layer** owns all business logic - validation, orchestration, error handling. Handlers only parse requests and call services.

**Repository interfaces** defined in `internal/repository/interfaces.go`, implementations in subpackages. Enables testing services with mocks.

```go
type VideoRepository interface {
    Create(ctx context.Context, video *domain.Video) error
    GetByID(ctx context.Context, id string) (*domain.Video, error)
    List(ctx context.Context, opts ListOptions) ([]domain.Video, error)
    Delete(ctx context.Context, id string) error
}
```

**Dependency injection** via constructors - no global state. Wire everything in `main.go`.

## Frontend Architecture

```
web/src/
├── app.ts                       # Inertia app bootstrap
├── pages/
│   ├── Home.vue                 # Landing page
│   ├── Gallery.vue              # Video grid
│   ├── Video/Show.vue           # Single video player
│   ├── Upload.vue               # File upload dropzone
│   ├── Record.vue               # WebRTC recorder
│   └── Edit.vue                 # Video metadata editor
├── components/
│   ├── ui/                      # shadcn-vue components
│   ├── VideoCard.vue            # Gallery thumbnail card
│   ├── VideoPlayer.vue          # Playback wrapper
│   ├── VideoRecorder.vue        # Camera capture
│   └── Layout.vue               # Shared nav/footer
├── composables/
│   ├── useMediaDevices.ts       # Camera/mic access
│   ├── useVideoRecorder.ts      # Recording logic
│   └── useUpload.ts             # Upload progress tracking
├── lib/
│   └── utils.ts                 # shadcn-vue utilities (cn helper)
└── assets/
    └── index.css                # Tailwind + custom properties
```

**shadcn-vue components:** Button, Card, Input, Textarea, Dialog, Progress, DropdownMenu, Toast.

**Inertia shared data** for navigation state, flash messages - passed from Go via `inertia.Share()`.

## Database Schema

```sql
-- 001_initial.up.sql

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
);

CREATE TABLE tags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE video_tags (
    video_id TEXT REFERENCES videos(id) ON DELETE CASCADE,
    tag_id INTEGER REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (video_id, tag_id)
);

CREATE INDEX idx_videos_created ON videos(created_at DESC);
```

**Migration on startup** - app runs pending migrations automatically.

**File storage** on filesystem (`storage/uploads/`, `storage/thumbnails/`). Database stores paths only. UUIDs for all filenames.

## Taskfile & Dev Workflow

```yaml
version: '3'

tasks:
  dev:
    desc: Start Go server and Vite dev server
    deps: [dev:backend, dev:frontend]

  dev:backend:
    desc: Run Go server with hot reload
    cmds:
      - air -c .air.toml
    env:
      APP_ENV: development
      VITE_DEV: "true"

  dev:frontend:
    desc: Run Vite dev server
    dir: web
    cmds:
      - npm run dev

  build:
    desc: Build production binary and frontend assets
    cmds:
      - task: build:frontend
      - task: build:backend

  build:frontend:
    dir: web
    cmds:
      - npm run build

  build:backend:
    cmds:
      - go build -o bin/gooji ./cmd/server

  test:
    desc: Run all tests
    cmds:
      - go test -v -race -cover ./internal/...

  lint:
    cmds:
      - go vet ./...
      - staticcheck ./...
```

**In development:** Go server proxies to Vite for HMR.
**In production:** Go serves pre-built assets from `web/dist/`.

## Testing Strategy

**Focus: Service layer unit tests**

```go
func TestVideoService_Upload(t *testing.T) {
    mockRepo := &mocks.VideoRepository{}
    mockFS := &mocks.FileStorage{}
    mockFFmpeg := &mocks.Processor{}

    svc := service.NewVideoService(mockRepo, mockFS, mockFFmpeg)

    mockFS.On("Save", mock.Anything, mock.Anything).Return("/path/to/file.webm", nil)
    mockFFmpeg.On("GetMetadata", mock.Anything).Return(&ffmpeg.Metadata{Duration: 30.5}, nil)
    mockFFmpeg.On("GenerateThumbnail", mock.Anything, mock.Anything).Return(nil)
    mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

    video, err := svc.Upload(ctx, uploadRequest)

    assert.NoError(t, err)
    assert.Equal(t, 30.5, video.DurationSeconds)
}
```

**What gets tested:**
- `VideoService`: Upload validation, metadata extraction, error handling
- `Repository` (SQLite): CRUD operations against test database
- `FFmpegProcessor`: Command building (actual FFmpeg mocked)

**Test utilities:**
- In-memory SQLite for repository tests
- Mock interfaces generated with `mockery`
- Table-driven tests for validation edge cases

## Medicine Wheel Color Theme

**Colors used subtly:**

```css
@theme {
  --color-earth-black: oklch(0.15 0.01 260);    /* Deep charcoal */
  --color-snow-white: oklch(0.97 0.01 90);      /* Warm white */
  --color-fire-red: oklch(0.55 0.18 25);        /* Muted terracotta */
  --color-sun-yellow: oklch(0.80 0.14 85);      /* Soft gold */
}
```

**Application:**
- **Black**: Primary text, headings, footer background
- **White**: Page backgrounds, card surfaces
- **Red**: Primary action buttons, active states, recording indicator
- **Yellow**: Accent highlights, hover states, success feedback

**Subtle usage:**
- Thin colored borders on cards, not filled backgrounds
- Yellow underline on active nav item
- Red glow on record button when active
- Black/white base with color accents only

## Cleanup Plan

**Remove entirely:**
- `web/templates/` - all Go HTML templates
- `web/static/js/` - vanilla JS files
- `web/static/css/` - old stylesheets
- Template rendering code in handlers
- `config/config.go` duplicate

**Refactor:**
- `internal/video/handler.go` → split into `internal/handler/video.go`, `gallery.go`, `upload.go`
- `internal/video/service.go` → `internal/service/video.go` with interfaces
- `internal/video/repository.go` → `internal/repository/sqlite/video.go`
- `pkg/ffmpeg/processor.go` → add interface, improve error handling

**Data migration:**
```bash
task migrate:data  # One-time: JSON → SQLite
```

## Implementation Phases

### Phase 1: Foundation
1. Set up Vue 3 + Vite + TypeScript in `web/`
2. Install and configure shadcn-vue with medicine wheel theme
3. Add Gonertia to Go backend
4. Create `task dev` with concurrent backend/frontend
5. Single test page proving Inertia round-trip works

### Phase 2: Database
6. Add SQLite with migrations
7. Create repository interfaces and SQLite implementations
8. Write data migration script (JSON → SQLite)
9. Repository unit tests

### Phase 3: Backend Restructure
10. Extract domain types to `internal/domain/`
11. Create service layer with proper interfaces
12. Convert handlers to use Gonertia
13. Service unit tests
14. Delete old template code

### Phase 4: Frontend Pages
15. Layout component with navigation
16. Home page
17. Gallery + VideoCard
18. Video playback page
19. Upload page with dropzone
20. Record page with WebRTC
21. Edit page

### Phase 5: Polish
22. Clean up unused code, comments
23. Update CLAUDE.md
24. Final lint pass, test coverage check

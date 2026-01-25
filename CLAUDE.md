# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Gooji (meaning "video" in Ojibwe) is a Go-based video platform for recording and sharing videos focused on Ojibwe language and culture. It provides a web interface for community members to record, upload, and browse cultural content.

**Tech Stack:**
- **Backend:** Go 1.24+, net/http, Inertia.js adapter, SQLite, Uber Zap logging
- **Frontend:** Vue 3 + TypeScript, Inertia.js, Vite, Tailwind CSS v4
- **Video Processing:** FFmpeg for thumbnails and metadata extraction

## Common Commands

```bash
# Development
task dev              # Run Go backend with hot reload
task frontend:dev     # Run Vite dev server for frontend

# Build
task build            # Build Go binary to bin/gooji
task frontend:build   # Build frontend to web/dist/

# Testing
task test             # Run Go tests with coverage
task test-coverage    # Generate HTML coverage report

# Code Quality
task lint             # Run linting (vet, staticcheck, govulncheck)
task format           # Format Go code

# Setup
task setup            # First-time development environment setup
task check-ffmpeg     # Verify FFmpeg installation
```

## Architecture

The project follows Clean Architecture with handler → service → repository layers.

### Backend Structure

```
main.go                           # Entry point, HTTP server, Inertia setup
internal/
├── domain/                       # Core business entities
│   ├── video.go                  # Video, Tag, VideoListOptions
│   └── errors.go                 # Domain errors (ErrVideoNotFound, etc.)
├── service/                      # Business logic layer
│   └── video.go                  # VideoService (upload, list, delete)
├── repository/                   # Data access layer
│   ├── interfaces.go             # VideoRepository, FileStorage interfaces
│   ├── sqlite/                   # SQLite implementation
│   │   ├── db.go                 # Database connection, migrations
│   │   └── video.go              # Video CRUD operations
│   └── filesystem/               # File storage implementation
│       └── storage.go            # Save, delete, path management
├── handler/                      # HTTP handlers
│   ├── inertia.go                # Inertia.js adapter setup
│   ├── pages.go                  # Page handlers (Home, Gallery, etc.)
│   └── video.go                  # Video API handlers
├── config/                       # Configuration with path validation
├── middleware/                   # CORS, logging, panic recovery
└── logger/                       # Zap structured logging
pkg/ffmpeg/processor.go           # FFmpeg wrapper with security checks
```

### Frontend Structure

```
web/
├── src/
│   ├── main.ts                   # Inertia app initialization
│   ├── types.ts                  # Shared TypeScript types
│   ├── pages/                    # Inertia pages (route components)
│   │   ├── Home.vue              # Landing page
│   │   ├── Gallery.vue           # Video gallery with search/filter
│   │   ├── Upload.vue            # Drag-drop video upload
│   │   ├── Record.vue            # WebRTC recording studio
│   │   └── Video/Show.vue        # Video player page
│   ├── components/               # Reusable Vue components
│   │   ├── Layout.vue            # App layout with nav/footer
│   │   └── ui/                   # UI primitives (Button, etc.)
│   └── lib/
│       └── utils.ts              # Utility functions
├── vite.config.ts                # Vite configuration
└── package.json
```

### Key Patterns

- **Inertia.js:** Server-driven SPA - Go renders Inertia responses, Vue handles UI
- **Dependency injection** via constructors for testability
- **Interface-based design** - Repository interfaces allow mock testing
- **SQLite** for metadata, filesystem for video/thumbnail storage

## API Routes

### Pages (Inertia)
- `GET /` - Home page
- `GET /gallery` - Video gallery
- `GET /record` - Recording interface
- `GET /upload` - Video upload page
- `GET /video/{id}` - Video player

### API (JSON)
- `GET /api/videos` - List videos
- `POST /api/videos` - Upload video (multipart form)
- `GET /api/videos/{id}` - Stream video file
- `DELETE /api/videos/{id}` - Delete video
- `GET /api/thumbnails?id={id}` - Get thumbnail
- `GET /health` - Health check

## Security Requirements

**Path validation is mandatory for all file operations.** Use validation functions from `internal/config/config.go`:
- Prevent path traversal (no `..` in paths)
- Block dangerous characters: `|`, `&`, `;`, `` ` ``, `$`, `(`, `)`, `{`, `}`, `[`, `]`, `*`, `?`, `\`
- Use secure permissions: 0o600 for files, 0o750 for directories

**FFmpeg security:** Use `NewProcessorWithSecurity()` for operations involving user uploads. The processor validates paths and restricts operations to allowed directories.

**Input sanitization:** All user text input is HTML-escaped in the service layer before storage.

## Code Standards

### Go
- Use wrapped errors: `fmt.Errorf("context: %w", err)`
- Never ignore errors with `_`
- Use structured logging via `internal/logger` (not fmt.Printf)
- Keep packages in `internal/` private; use `pkg/` for reusable code
- Run `task lint` before commits

### Vue/TypeScript
- Use Vue 3 Composition API with `<script setup>`
- TypeScript strict mode enabled
- Tailwind CSS for styling (no custom CSS unless necessary)
- Components use props interfaces via `defineProps<T>()`

## Testing

```bash
# Run all Go tests
go test ./...

# Run with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific package tests
go test ./internal/service/...
go test ./internal/repository/...
```

Test files are co-located with source: `video.go` → `video_test.go`

## Configuration

Configuration loads from `config/config.json` with environment variable overrides:
- `APP_DEBUG=true` - Enable debug logging
- `GOOJI_PORT=8080` - Override server port

## External Dependencies

- **FFmpeg 6.1.1+** - Required for video processing (thumbnail generation, transcoding)
- **Task runner** - For running build/dev commands
- **Node.js 18+** - For frontend build tooling

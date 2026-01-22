# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Gooji (meaning "video" in Ojibwe) is a Go-based video kiosk application for recording and editing videos focused on Ojibwe language and culture. It provides a web interface for community members to record, edit, and share cultural content.

**Tech Stack:** Go 1.24+, built-in `net/http`, Uber Zap logging, FFmpeg for video processing, HTML templates with vanilla JavaScript.

## Common Commands

```bash
# Development with hot reload
task dev

# Production mode
task run

# Build binary to bin/gooji
task build

# Run tests with coverage
task test

# Generate HTML coverage report
task test-coverage

# Run linting (vet, staticcheck, govulncheck)
task lint

# Format code
task format

# Security vulnerability scan
task security-scan

# Setup development environment (first time)
task setup

# Check FFmpeg installation
task check-ffmpeg
```

## Architecture

The project follows Clean Architecture with handler → service → repository layers:

```
main.go                     # Entry point, HTTP server setup
internal/
├── video/                  # Core domain - video management
│   ├── handler.go          # HTTP handlers, template rendering
│   ├── service.go          # Business logic, validation
│   └── repository.go       # File-based persistence
├── config/                 # Configuration with path validation
├── logger/                 # Zap structured logging
├── middleware/             # CORS, logging, panic recovery
└── errors/                 # Domain error types
pkg/ffmpeg/processor.go     # FFmpeg wrapper with security checks
web/
├── templates/              # HTML templates (base.html + content blocks)
└── static/                 # CSS and JavaScript
storage/                    # Runtime data (uploads, logs, metadata, thumbnails)
```

**Key patterns:**
- Dependency injection via constructors
- Interface-based design for testability
- Template system uses base template with content blocks per page
- File-based storage (JSON metadata, video files in storage/uploads)

## API Routes

- `GET /` - Home page
- `GET /record` - Recording interface
- `GET /upload` - Video upload page
- `GET /gallery` - Video gallery
- `GET /edit/{id}` - Video editor
- `GET /health` - Health check
- `POST /api/videos` - Upload video (multipart form)
- `GET /api/videos` - List videos (JSON)
- `GET /api/videos/{id}` - Get video file
- `DELETE /api/videos/{id}` - Delete video

## Security Requirements

**Path validation is mandatory for all file operations.** Use validation functions from `internal/config/config.go`:
- Prevent path traversal (no `..` in paths)
- Block dangerous characters: `|`, `&`, `;`, `` ` ``, `$`, `(`, `)`, `{`, `}`, `[`, `]`, `*`, `?`, `\`
- Use secure permissions: 0o600 for files, 0o750 for directories

**FFmpeg security:** Use `NewProcessorWithSecurity()` for operations involving user uploads. The processor validates paths and restricts operations to allowed directories.

**When using `//nolint:gosec`:** Only allowed when path validation has been explicitly verified.

## Code Standards

- Use wrapped errors: `fmt.Errorf("context: %w", err)`
- Never ignore errors with `_`
- Use structured logging via `internal/logger` (not fmt.Printf)
- Keep packages in `internal/` private; use `pkg/` for reusable code
- Run `task lint` before commits

## Configuration

Configuration loads from `config/config.json` with environment variable overrides:
- `APP_DEBUG=true` - Enable debug logging
- `GOOJI_PORT=8080` - Override server port

## External Dependencies

- **FFmpeg 6.1.1+** - Required for video processing (thumbnail generation, transcoding)
- **Task runner** - For running build/dev commands
- **Air** - Hot reload during development (auto-installed via `task setup`)

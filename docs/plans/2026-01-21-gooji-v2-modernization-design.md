# Gooji v2 Modernization Design

**Date:** 2026-01-21
**Status:** Approved
**Approach:** Clean rebuild

## Overview

Transform Gooji from a kiosk-only prototype to a web platform while preserving the frictionless contribution model. Open contribution (no accounts), device-based trust (kiosks auto-publish, web uploads go to moderation queue), CLI-based admin tooling.

## Key Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Contribution model | Open (no accounts) | Keeps kiosk spirit, low barrier |
| Moderation | Hybrid (device-based trust) | Kiosk friction-free, web protected |
| Admin interface | CLI tool | Simple, no web admin to secure |
| Frontend | Vue 3 + Inertia.js | SPA feel, server-driven routing |
| Database | SQLite | Simple ops, sufficient scale |
| Video storage | Local filesystem | Simple deployment |
| Approach | Clean rebuild | Scope too large for incremental |

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        Clients                               │
├──────────────────┬──────────────────┬───────────────────────┤
│   Kiosk Device   │   Web Browser    │   (Future: Android)   │
│   (trusted)      │   (untrusted)    │                       │
└────────┬─────────┴────────┬─────────┴───────────────────────┘
         │                  │
         ▼                  ▼
┌─────────────────────────────────────────────────────────────┐
│                    Go 1.25 Backend                          │
├─────────────────────────────────────────────────────────────┤
│  Inertia.js Adapter  │  Device Auth  │  Moderation Queue   │
├─────────────────────────────────────────────────────────────┤
│  Video Service  │  Upload Service  │  Thumbnail Service    │
├─────────────────────────────────────────────────────────────┤
│                     SQLite + Local FS                       │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                    CLI Tools                                 │
├─────────────────────────────────────────────────────────────┤
│  gooji-admin: review, approve, reject, list-pending         │
└─────────────────────────────────────────────────────────────┘
```

**Key flows:**
- **Kiosk upload** → Device authenticated → Video published immediately
- **Web upload** → No auth required → Video queued for review
- **Admin moderation** → CLI tool reviews queue → Approve/reject

## Backend Structure

```
gooji/
├── cmd/
│   ├── server/main.go        # Web server entry point
│   └── admin/main.go         # CLI moderation tool
├── internal/
│   ├── domain/               # Core business entities (no dependencies)
│   │   ├── video.go          # Video, VideoStatus, VideoMetadata
│   │   ├── device.go         # Device, TrustLevel
│   │   └── errors.go         # Domain errors
│   ├── service/              # Business logic (depends only on domain)
│   │   ├── video.go          # Upload, list, get, delete
│   │   ├── moderation.go     # Queue, approve, reject
│   │   └── thumbnail.go      # Generation, caching
│   ├── repository/           # Data access interfaces + implementations
│   │   ├── interfaces.go     # VideoRepo, DeviceRepo interfaces
│   │   ├── sqlite/           # SQLite implementations
│   │   └── filesystem/       # Video file storage
│   ├── handler/              # HTTP handlers (Inertia responses)
│   │   ├── video.go
│   │   ├── gallery.go
│   │   └── upload.go
│   ├── middleware/           # Auth, logging, recovery, rate limiting
│   └── ffmpeg/               # FFmpeg wrapper
├── web/                      # Vue.js frontend (Vite build)
├── migrations/               # SQLite migrations
└── config/
```

**SRP enforcement:**
- Handlers only parse requests and return responses
- Services contain business logic, no HTTP awareness
- Repositories handle data persistence only
- FFmpeg wrapper handles only video processing commands

## Database Schema

```sql
-- Videos: core content entity
CREATE TABLE videos (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    language TEXT DEFAULT 'oji',
    duration_seconds INTEGER,
    file_path TEXT NOT NULL,
    thumbnail_path TEXT,
    file_size_bytes INTEGER,
    mime_type TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',  -- pending, published, rejected
    rejection_reason TEXT,
    device_id TEXT REFERENCES devices(id),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    published_at DATETIME,
    moderated_at DATETIME
);

-- Devices: trust tracking for kiosks
CREATE TABLE devices (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    api_key_hash TEXT NOT NULL,
    trust_level TEXT NOT NULL,  -- trusted, untrusted
    last_seen_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Tags: categorization
CREATE TABLE tags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE video_tags (
    video_id TEXT REFERENCES videos(id) ON DELETE CASCADE,
    tag_id INTEGER REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (video_id, tag_id)
);

-- Indexes
CREATE INDEX idx_videos_status ON videos(status);
CREATE INDEX idx_videos_created ON videos(created_at DESC);
CREATE INDEX idx_videos_device ON videos(device_id);
```

## Frontend Structure

```
web/
├── src/
│   ├── app.ts                    # Inertia app initialization
│   ├── pages/                    # Inertia pages (map to routes)
│   │   ├── Home.vue
│   │   ├── Gallery.vue
│   │   ├── Video/Show.vue
│   │   ├── Upload.vue
│   │   └── Record.vue
│   ├── components/               # Reusable components
│   │   ├── VideoCard.vue
│   │   ├── VideoPlayer.vue
│   │   ├── VideoRecorder.vue
│   │   ├── UploadDropzone.vue
│   │   └── ui/
│   │       ├── Button.vue
│   │       ├── Input.vue
│   │       └── Modal.vue
│   ├── composables/              # Shared logic
│   │   ├── useVideoRecorder.ts
│   │   ├── useUpload.ts
│   │   └── useMediaDevices.ts
│   ├── types/
│   │   └── index.ts
│   └── assets/
│       └── css/
│           └── app.css           # Tailwind CSS
├── index.html
├── vite.config.ts
├── tsconfig.json
├── tailwind.config.js
└── package.json
```

**Stack:** Vue 3 Composition API, TypeScript, Vite, Tailwind CSS, Inertia.js

## Device Authentication & Trust

```go
func DeviceAuth(deviceRepo repository.DeviceRepo) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            ctx := r.Context()

            apiKey := r.Header.Get("X-Device-Key")
            if apiKey == "" {
                // No key = untrusted web user
                ctx = context.WithValue(ctx, trustLevelKey, domain.Untrusted)
            } else {
                device, err := deviceRepo.ValidateAPIKey(ctx, apiKey)
                if err != nil {
                    ctx = context.WithValue(ctx, trustLevelKey, domain.Untrusted)
                } else {
                    ctx = context.WithValue(ctx, trustLevelKey, device.TrustLevel)
                    ctx = context.WithValue(ctx, deviceKey, device)
                }
            }

            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

**Kiosk setup flow:**
1. Admin runs `gooji-admin device add --name "Library Kiosk 1"`
2. CLI generates UUID + API key, displays once, stores hash
3. API key configured on kiosk device (environment variable)
4. Kiosk includes `X-Device-Key` header on all requests

**Upload behavior:**
- Trusted device → `status = published`
- Untrusted (web) → `status = pending`

## CLI Moderation Tool

```bash
# Queue management
$ gooji-admin queue list
$ gooji-admin queue review <id>    # Opens in browser
$ gooji-admin queue approve <id>
$ gooji-admin queue reject <id> --reason "..."

# Device management
$ gooji-admin device add --name "..."
$ gooji-admin device list
$ gooji-admin device revoke <id>
```

Uses `cobra` for CLI framework, shares `internal/` packages with server.

## Security

**Input Validation:**
- Centralized validation using `go-playground/validator`
- HTML sanitization on all text inputs
- Validation in service layer

**File Upload:**
- Magic byte validation
- Max file size at middleware level
- UUID filenames only
- Videos stored outside web root

**FFmpeg:**
- Allowlist of permitted operations
- No user input in command arguments
- Timeout enforcement
- Sandboxed execution directory

**Rate Limiting:**
- Token bucket per IP (5 uploads/hour untrusted)
- Separate limits for trusted devices

**CORS:**
- Configured allowed origins (no `*`)
- Environment-specific settings

## Implementation Phases

### Phase 1: Foundation
- Go 1.25 project scaffold with clean architecture
- SQLite setup with migrations
- Domain entities and repository interfaces
- Basic Inertia.js + Vue 3 + Vite setup
- Health check endpoint
- *Deliverable: App boots, serves empty gallery*

### Phase 2: Core Video Flow
- Video upload (file handling, FFmpeg metadata extraction)
- Thumbnail generation
- Gallery listing and video playback
- *Deliverable: Can upload and watch videos*

### Phase 3: Device Trust & Moderation
- Device authentication middleware
- Moderation queue
- CLI tool: queue commands
- CLI tool: device commands
- *Deliverable: Kiosk auto-publishes, web queued*

### Phase 4: Recording
- WebRTC video recorder component
- Recording preview and retry
- Direct upload from recorder
- *Deliverable: Full record-to-publish flow*

### Phase 5: Polish & Harden
- Rate limiting
- Input validation hardening
- Error handling and user feedback
- Responsive design pass
- *Deliverable: Production-ready*

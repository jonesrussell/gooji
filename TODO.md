# Gooji Project - TODO & Technical Debt Roadmap

## 🚨 **CRITICAL PRIORITY** (Must Fix Immediately)

### Security Vulnerabilities

#### [SEC-001] Command Injection Risk in FFmpeg Integration
- **Status**: 🔴 Critical
- **Effort**: 1-2 days
- **Location**: `pkg/ffmpeg/processor.go`
- **Issue**: Direct use of user inputs in `exec.Command` without sanitization
- **Risk**: Potential command injection attacks
- **Tasks**:
  - [ ] Validate and sanitize all file paths before passing to FFmpeg
  - [ ] Implement allowlist for FFmpeg parameters
  - [ ] Add path traversal protection
  - [ ] Create secure wrapper for FFmpeg command execution
- **Implementation**:
```go
// Add to pkg/ffmpeg/processor.go
func (p *Processor) sanitizePath(path string) (string, error) {
    // Validate path is within allowed directory
    // Check for path traversal attempts
    // Return sanitized path
}
```

#### [SEC-002] File Upload Security Gaps
- **Status**: 🔴 Critical
- **Effort**: 2-3 days
- **Location**: `internal/video/handler.go` - `HandleUpload()`
- **Issues**:
  - No MIME type validation beyond form data
  - No file size limits enforced in code
  - No malicious file scanning
  - User-controlled filenames could cause path traversal
- **Tasks**:
  - [ ] Implement proper MIME type validation
  - [ ] Add file size enforcement in handler
  - [ ] Use UUIDs for filenames instead of user input
  - [ ] Add file content validation
  - [ ] Implement virus scanning (optional)
- **Implementation**:
```go
// Add to internal/video/handler.go
func (h *Handler) validateUpload(file multipart.File, header *multipart.FileHeader) error {
    // Check file size
    // Validate MIME type
    // Check file extension
    // Scan for malicious content
}
```

#### [SEC-003] CORS Misconfiguration
- **Status**: 🔴 Critical
- **Effort**: 0.5 days
- **Location**: `internal/middleware/middleware.go`
- **Issue**: `Access-Control-Allow-Origin: "*"` allows any domain
- **Tasks**:
  - [ ] Configure specific allowed origins in config
  - [ ] Update CORS middleware to use config values
  - [ ] Add environment-specific CORS settings
- **Implementation**:
```go
// Update internal/middleware/middleware.go
func CORS(allowedOrigins []string) Middleware {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            origin := r.Header.Get("Origin")
            if contains(allowedOrigins, origin) {
                w.Header().Set("Access-Control-Allow-Origin", origin)
            }
            // ... rest of CORS logic
        })
    }
}
```

#### [SEC-004] Input Validation for Video Metadata
- **Status**: 🔴 Critical
- **Effort**: 1 day
- **Location**: `internal/video/handler.go`
- **Issue**: Title, description, and tags not validated for XSS/injection
- **Tasks**:
  - [ ] Implement HTML sanitization for user inputs
  - [ ] Add length limits for metadata fields
  - [ ] Validate tag format and content
  - [ ] Implement content filtering
- **Implementation**:
```go
import "html"

func sanitizeMetadata(metadata *VideoMetadata) error {
    metadata.Title = html.EscapeString(strings.TrimSpace(metadata.Title))
    metadata.Description = html.EscapeString(strings.TrimSpace(metadata.Description))
    
    if len(metadata.Title) > 100 {
        return errors.New("title too long")
    }
    
    // Validate and sanitize tags
    for i, tag := range metadata.Tags {
        metadata.Tags[i] = html.EscapeString(strings.TrimSpace(tag))
    }
    
    return nil
}
```

#### [SEC-005] Rate Limiting for Upload Endpoints
- **Status**: 🔴 Critical
- **Effort**: 1-2 days
- **Location**: `internal/middleware/middleware.go`
- **Issue**: No protection against abuse/DoS attacks
- **Tasks**:
  - [ ] Implement rate limiting middleware
  - [ ] Add IP-based request limiting
  - [ ] Configure different limits for different endpoints
- **Implementation**:
```go
import "golang.org/x/time/rate"

func RateLimitMiddleware(r rate.Limit, b int) Middleware {
    limiter := rate.NewLimiter(r, b)
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
            if !limiter.Allow() {
                http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
                return
            }
            next.ServeHTTP(w, req)
        })
    }
}
```

### Testing Infrastructure

#### [TEST-001] Implement Basic Testing Framework
- **Status**: 🔴 Critical
- **Effort**: 3-4 days
- **Issue**: Zero test coverage across entire codebase
- **Tasks**:
  - [ ] Set up testing framework and tools
  - [ ] Create test utilities and helpers
  - [ ] Write unit tests for critical paths
  - [ ] Add integration tests for video processing
  - [ ] Set up CI/CD pipeline for tests
- **Implementation**:
```bash
# Add to Taskfile.yml
test-unit:
  desc: Run unit tests
  cmds:
    - go test -v -race -cover ./internal/...

test-integration:
  desc: Run integration tests
  cmds:
    - go test -v -tags=integration ./tests/...

test-coverage:
  desc: Generate test coverage report
  cmds:
    - go test -coverprofile=coverage.out ./...
    - go tool cover -html=coverage.out -o coverage.html
```

#### [TEST-002] Test Data Management
- **Status**: 🟡 High
- **Effort**: 1 day
- **Tasks**:
  - [ ] Create test video files of various formats/sizes
  - [ ] Set up test fixtures and helpers
  - [ ] Implement test cleanup utilities
- **Implementation**:
```go
// tests/fixtures/video_fixtures.go
func CreateTestVideo(t *testing.T, format string, duration int) string {
    // Generate test video file
    // Return path to test file
}

func CleanupTestFiles(t *testing.T, paths []string) {
    t.Cleanup(func() {
        for _, path := range paths {
            os.Remove(path)
        }
    })
}
```

#### [TEST-003] Mock FFmpeg for Testing
- **Status**: 🟡 High
- **Effort**: 1-2 days
- **Tasks**:
  - [ ] Create FFmpeg mock interface
  - [ ] Implement test doubles
  - [ ] Add dependency injection for testing
- **Implementation**:
```go
type FFmpegProcessor interface {
    GetVideoInfo(inputPath string) (*VideoInfo, error)
    GenerateThumbnail(inputPath, outputPath string, timestamp float64) error
}

type MockFFmpegProcessor struct {
    // Mock implementation
}
```

## ⚠️ **HIGH PRIORITY** (Fix Within 1-2 Weeks)

### Enhanced Security Measures

#### [SEC-006] Content Security Policy (CSP)
- **Status**: 🟡 High
- **Effort**: 0.5 days
- **Issue**: No CSP headers for XSS protection
- **Tasks**:
  - [ ] Add CSP middleware
  - [ ] Configure appropriate CSP directives
  - [ ] Test with frontend functionality
- **Implementation**:
```go
func CSPMiddleware() Middleware {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Content-Security-Policy", 
                "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline' https://cdn.jsdelivr.net")
            next.ServeHTTP(w, r)
        })
    }
}
```

#### [SEC-007] Secure Headers Package
- **Status**: 🟡 High
- **Effort**: 0.5 days
- **Tasks**:
  - [ ] Add security headers middleware
  - [ ] Implement HSTS, X-Frame-Options, etc.
- **Implementation**:
```go
func SecurityHeadersMiddleware() Middleware {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("X-Content-Type-Options", "nosniff")
            w.Header().Set("X-Frame-Options", "DENY")
            w.Header().Set("X-XSS-Protection", "1; mode=block")
            w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
            next.ServeHTTP(w, r)
        })
    }
}
```

### Code Architecture Issues

#### [ARCH-001] Violation of Single Responsibility Principle
- **Status**: 🟡 High
- **Effort**: 3-4 days
- **Location**: `internal/video/handler.go`
- **Issue**: `HandleUpload()` method handles multiple responsibilities
- **Tasks**:
  - [ ] Extract file parsing logic into separate method
  - [ ] Create dedicated video processing service
  - [ ] Separate metadata handling
  - [ ] Extract thumbnail generation logic
- **Implementation**:
```go
// Refactor internal/video/handler.go
type VideoService interface {
    ProcessUpload(ctx context.Context, upload *Upload) (*Video, error)
    GenerateThumbnail(ctx context.Context, videoPath string) error
    ValidateVideo(ctx context.Context, file io.Reader) error
}

type FileService interface {
    SaveVideo(ctx context.Context, file multipart.File, filename string) (string, error)
    SaveMetadata(ctx context.Context, metadata *VideoMetadata) error
}
```

#### [ARCH-002] Duplicate Configuration Systems
- **Status**: 🟡 High
- **Effort**: 1-2 days
- **Issue**: Two different config structs and loading mechanisms
- **Tasks**:
  - [ ] Consolidate `config/config.go` and `internal/config/config.go`
  - [ ] Choose single configuration approach
  - [ ] Update all references to use unified config
  - [ ] Add configuration validation
- **Implementation**:
```go
// Unified config structure
type Config struct {
    Server   ServerConfig   `json:"server" validate:"required"`
    Video    VideoConfig    `json:"video" validate:"required"`
    FFmpeg   FFmpegConfig   `json:"ffmpeg" validate:"required"`
    Security SecurityConfig `json:"security" validate:"required"`
}

func (c *Config) Validate() error {
    // Add validation logic
}
```

### Error Handling

#### [ERR-001] Standardize Error Handling
- **Status**: 🟡 High
- **Effort**: 2-3 days
- **Issue**: Mixed error handling patterns throughout codebase
- **Tasks**:
  - [ ] Replace `fmt.Printf` with structured logging
  - [ ] Create consistent error types
  - [ ] Add error context and wrapping
  - [ ] Implement proper error responses
- **Implementation**:
```go
// Add to internal/errors/errors.go
type AppError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Err     error  `json:"-"`
}

func (e *AppError) Error() string {
    return fmt.Sprintf("%s: %v", e.Message, e.Err)
}

// Update all error handling to use structured approach
```

#### [ERR-002] Add Error Context
- **Status**: 🟡 High
- **Effort**: 1-2 days
- **Issue**: Generic error messages lack context
- **Tasks**:
  - [ ] Add file paths to error messages
  - [ ] Include operation context in errors
  - [ ] Add request ID tracking
  - [ ] Implement error correlation
- **Implementation**:
```go
// Enhanced error context
return fmt.Errorf("failed to process video %s: %w", videoPath, err)
```

## 📊 **MEDIUM PRIORITY** (Fix Within 1 Month)

### Monitoring & Observability

#### [MON-001] Structured Logging
- **Status**: 🟡 High
- **Effort**: 1 day
- **Tasks**:
  - [ ] Replace custom logger with structured logging (logrus/zap)
  - [ ] Add request correlation IDs
  - [ ] Implement log levels and configuration
- **Implementation**:
```go
import "go.uber.org/zap"

type Logger struct {
    logger *zap.Logger
}

func (l *Logger) InfoWithFields(msg string, fields map[string]interface{}) {
    l.logger.Info(msg, zap.Any("fields", fields))
}
```

#### [MON-002] Health Check Endpoints
- **Status**: 🟡 Medium
- **Effort**: 0.5 days
- **Tasks**:
  - [ ] Add /health endpoint
  - [ ] Check FFmpeg availability
  - [ ] Check video directory accessibility
- **Implementation**:
```go
func (h *Handler) HandleHealth(w http.ResponseWriter, r *http.Request) {
    health := map[string]interface{}{
        "status": "healthy",
        "timestamp": time.Now(),
        "checks": map[string]bool{
            "ffmpeg": h.checkFFmpeg(),
            "video_dir": h.checkVideoDir(),
        },
    }
    
    json.NewEncoder(w).Encode(health)
}
```

### Resource Management

#### [RES-001] Missing Cleanup in Error Paths
- **Status**: 🟡 Medium
- **Effort**: 1-2 days
- **Location**: `internal/video/handler.go`
- **Issue**: Resources not cleaned up on errors
- **Tasks**:
  - [ ] Add defer statements for cleanup
  - [ ] Implement proper resource management
  - [ ] Add cleanup on partial failures
  - [ ] Create cleanup utilities
- **Implementation**:
```go
// Add cleanup utilities
func (h *Handler) cleanupOnError(paths []string) {
    for _, path := range paths {
        os.Remove(path)
    }
}
```

#### [RES-002] Graceful Shutdown
- **Status**: 🟡 Medium
- **Effort**: 1 day
- **Location**: `main.go`
- **Issue**: Resources not properly closed on shutdown
- **Tasks**:
  - [ ] Add cleanup in shutdown handler
  - [ ] Close logger properly
  - [ ] Clean up temporary files
  - [ ] Add shutdown timeout handling
- **Implementation**:
```go
// Enhanced shutdown in main.go
defer func() {
    log.Info("Cleaning up resources...")
    // Close logger
    // Clean up temp files
    // Close database connections
}()
```

### Code Organization

#### [ORG-001] Break Up Large Methods
- **Status**: 🟡 Medium
- **Effort**: 2-3 days
- **Issue**: Multiple methods exceed 30 lines
- **Tasks**:
  - [ ] Refactor `HandleUpload()` into smaller methods
  - [ ] Extract common patterns into utilities
  - [ ] Create helper functions for repeated logic
  - [ ] Add method documentation
- **Implementation**:
```go
// Break down HandleUpload
func (h *Handler) HandleUpload(w http.ResponseWriter, r *http.Request) {
    file, metadata, err := h.parseUploadRequest(r)
    if err != nil {
        h.handleError(w, err)
        return
    }
    
    savedPath, err := h.saveVideoFile(file, metadata)
    if err != nil {
        h.handleError(w, err)
        return
    }
    
    if err := h.processVideo(savedPath, metadata); err != nil {
        h.handleError(w, err)
        return
    }
    
    h.respondSuccess(w, metadata)
}
```

#### [ORG-002] Frontend Code Duplication
- **Status**: 🟡 Medium
- **Effort**: 1-2 days
- **Issue**: Repeated patterns across JavaScript files
- **Tasks**:
  - [ ] Create shared API client utility
  - [ ] Extract common UI components
  - [ ] Standardize error handling in frontend
  - [ ] Add frontend testing
- **Implementation**:
```javascript
// Create web/static/js/api-client.js
class ApiClient {
    async uploadVideo(formData) {
        const response = await fetch('/api/videos', {
            method: 'POST',
            body: formData
        });
        return this.handleResponse(response);
    }
    
    async getVideos() {
        const response = await fetch('/api/videos');
        return this.handleResponse(response);
    }
    
    handleResponse(response) {
        if (!response.ok) {
            throw new Error(`HTTP ${response.status}: ${response.statusText}`);
        }
        return response.json();
    }
}
```

## 💡 **LOW PRIORITY** (Fix Within 2-3 Months)

### DevOps & Deployment

#### [DEVOPS-001] Docker Support
- **Status**: 🟢 Low
- **Effort**: 1-2 days
- **Tasks**:
  - [ ] Create optimized Dockerfile
  - [ ] Add docker-compose for development
  - [ ] Configure volume mounts for videos
- **Implementation**:
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o gooji cmd/gooji/main.go

FROM alpine:latest
RUN apk add --no-cache ffmpeg
WORKDIR /app
COPY --from=builder /app/gooji .
COPY --from=builder /app/web ./web
EXPOSE 8080
CMD ["./gooji"]
```

#### [DEVOPS-002] CI/CD Pipeline
- **Status**: 🟢 Low
- **Effort**: 2-3 days
- **Tasks**:
  - [ ] Set up GitHub Actions workflow
  - [ ] Add automated testing
  - [ ] Configure security scanning
  - [ ] Add deployment automation

### Documentation

#### [DOC-001] Add Package Documentation
- **Status**: 🟢 Low
- **Effort**: 1-2 days
- **Issue**: Missing package-level documentation
- **Tasks**:
  - [ ] Add package comments to all packages
  - [ ] Document exported types and functions
  - [ ] Add usage examples
  - [ ] Create API documentation
- **Implementation**:
```go
// Package video provides video processing and management functionality
// for the Gooji kiosk application.
//
// This package handles:
// - Video file uploads and validation
// - Metadata management
// - Thumbnail generation
// - Video serving
package video
```

#### [DOC-002] API Documentation
- **Status**: 🟢 Low
- **Effort**: 2-3 days
- **Issue**: No API documentation
- **Tasks**:
  - [ ] Add OpenAPI/Swagger documentation
  - [ ] Document all endpoints
  - [ ] Add request/response examples
  - [ ] Create API testing tools
- **Implementation**:
```go
// Add swagger annotations
// @Summary Upload a video file
// @Description Upload a video file with metadata
// @Tags videos
// @Accept multipart/form-data
// @Produce json
// @Param video formData file true "Video file"
// @Param title formData string false "Video title"
// @Success 200 {object} VideoMetadata
// @Failure 400 {object} ErrorResponse
// @Router /api/videos [post]
func (h *Handler) HandleUpload(w http.ResponseWriter, r *http.Request) {
    // Implementation
}
```

### Performance Optimization

#### [PERF-001] Memory Usage Optimization
- **Status**: 🟢 Low
- **Effort**: 2-3 days
- **Issue**: Large files loaded entirely into memory
- **Tasks**:
  - [ ] Implement streaming for large files
  - [ ] Add memory usage monitoring
  - [ ] Optimize video processing pipeline
  - [ ] Add file size limits
- **Implementation**:
```go
// Streaming file upload
func (h *Handler) streamUpload(file multipart.File, dstPath string) error {
    dst, err := os.Create(dstPath)
    if err != nil {
        return err
    }
    defer dst.Close()
    
    buffer := make([]byte, 32*1024) // 32KB buffer
    for {
        n, err := file.Read(buffer)
        if n > 0 {
            if _, err := dst.Write(buffer[:n]); err != nil {
                return err
            }
        }
        if err == io.EOF {
            break
        }
        if err != nil {
            return err
        }
    }
    return nil
}
```

#### [PERF-002] Add Monitoring and Metrics
- **Status**: 🟢 Low
- **Effort**: 2-3 days
- **Issue**: No performance monitoring
- **Tasks**:
  - [ ] Add Prometheus metrics
  - [ ] Implement health checks
  - [ ] Add performance logging
  - [ ] Create monitoring dashboard
- **Implementation**:
```go
// Add to internal/monitoring/metrics.go
type Metrics struct {
    UploadCounter   prometheus.Counter
    ProcessingTime  prometheus.Histogram
    ErrorCounter    prometheus.Counter
    FileSizeGauge   prometheus.Gauge
}
```

## 🛠️ **IMPLEMENTATION ROADMAP**

### Phase 1: Critical Security & Testing (Week 1-2)
- [ ] Fix command injection vulnerabilities
- [ ] Implement file upload security
- [ ] Set up basic testing framework
- [ ] Add critical path tests

### Phase 2: Architecture & Error Handling (Week 3-4)
- [ ] Refactor large methods
- [ ] Standardize error handling
- [ ] Consolidate configuration
- [ ] Add proper resource cleanup

### Phase 3: Code Quality & Documentation (Week 5-6)
- [ ] Eliminate code duplication
- [ ] Add comprehensive documentation
- [ ] Implement frontend improvements
- [ ] Add API documentation

### Phase 4: Performance & Monitoring (Week 7-8)
- [ ] Optimize memory usage
- [ ] Add monitoring and metrics
- [ ] Implement performance testing
- [ ] Create deployment automation

## 📈 **SUCCESS METRICS**

### Code Quality
- [ ] Achieve >80% test coverage
- [ ] Reduce cyclomatic complexity to <10 per function
- [ ] Eliminate all security vulnerabilities
- [ ] Achieve 0 linting errors

### Performance
- [ ] Support file uploads up to 100MB
- [ ] Process videos within 30 seconds
- [ ] Maintain <100ms response time for API calls
- [ ] Support concurrent uploads

### Security Metrics
- [ ] Zero high-severity security vulnerabilities
- [ ] Pass OWASP security checklist
- [ ] Implement all recommended security headers
- [ ] Rate limiting prevents >99% of abuse attempts

### User Experience Metrics
- [ ] Video upload success rate >95%
- [ ] Page load time <2 seconds
- [ ] Mobile responsiveness score >90
- [ ] Accessibility compliance (WCAG 2.1 AA)

### Operational Metrics
- [ ] System uptime >99.5%
- [ ] Automated deployment success rate >95%
- [ ] Mean time to recovery <15 minutes
- [ ] Monitoring coverage >90%

### Maintainability
- [ ] Complete API documentation
- [ ] Add comprehensive logging
- [ ] Implement monitoring dashboard
- [ ] Create deployment automation

## 🔄 **ONGOING MAINTENANCE**

### Weekly Tasks
- [ ] Run security scans
- [ ] Update dependencies
- [ ] Review error logs
- [ ] Monitor performance metrics

### Monthly Tasks
- [ ] Review and update documentation
- [ ] Analyze test coverage
- [ ] Review security best practices
- [ ] Plan technical debt reduction

### Quarterly Tasks
- [ ] Major dependency updates
- [ ] Architecture review
- [ ] Performance optimization
- [ ] Security audit

## 🎯 **QUICK WINS TO PRIORITIZE**

### Week 1 Quick Wins (1-2 hours each)
1. **Fix CORS configuration** - immediate security improvement
2. **Add security headers** - easy security enhancement
3. **Implement input sanitization** - prevent XSS attacks
4. **Add rate limiting** - prevent abuse

### Week 2 Quick Wins
1. **Create health check endpoint** - operational visibility
2. **Add request logging** - debugging support
3. **Implement graceful shutdown** - reliability improvement
4. **Create test fixtures** - testing foundation

## 🚀 **IMPLEMENTATION TIPS**

### Development Workflow
1. **Branch Strategy**: Create feature branches for each TODO item
2. **Code Reviews**: Require peer review for all security-related changes
3. **Testing First**: Write tests before implementing fixes where possible
4. **Incremental Deployment**: Deploy fixes in small batches

### Validation Checklist
- [ ] All security fixes tested with penetration testing tools
- [ ] Performance impact measured before/after changes
- [ ] Documentation updated with each change
- [ ] Backward compatibility maintained where possible

---

**Last Updated**: $(date)
**Next Review**: $(date -d '+1 month')
**Priority**: Critical security fixes first, then systematic improvement 
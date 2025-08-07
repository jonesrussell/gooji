package ffmpeg

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// VideoInfo contains metadata about a video file
type VideoInfo struct {
	Duration   float64
	Width      int
	Height     int
	Format     string
	VideoCodec string
	AudioCodec string
	FrameRate  float64
	Bitrate    int
}

// Processor handles video processing operations using FFmpeg
type Processor struct {
	ffmpegPath string
	allowedDir string
}

// NewProcessor creates a new FFmpeg processor
func NewProcessor(ffmpegPath string) *Processor {
	if ffmpegPath == "" {
		ffmpegPath = "ffmpeg"
	}
	return &Processor{
		ffmpegPath: ffmpegPath,
	}
}

// NewProcessorWithSecurity creates a new FFmpeg processor with security restrictions
func NewProcessorWithSecurity(ffmpegPath, allowedDir string) *Processor {
	if ffmpegPath == "" {
		ffmpegPath = "ffmpeg"
	}
	return &Processor{
		ffmpegPath: ffmpegPath,
		allowedDir: allowedDir,
	}
}

// FFmpegPath returns the FFmpeg executable path
func (p *Processor) FFmpegPath() string {
	return p.ffmpegPath
}

// validatePath ensures a file path is secure and within allowed directory
func (p *Processor) validatePath(filePath string) error {
	if filePath == "" {
		return fmt.Errorf("file path cannot be empty")
	}

	// Check for path traversal attempts
	if strings.Contains(filePath, "..") {
		return fmt.Errorf("path traversal not allowed: %s", filePath)
	}

	// Check for absolute paths outside allowed directory
	if p.allowedDir != "" && filepath.IsAbs(filePath) {
		absAllowedDir, err := filepath.Abs(p.allowedDir)
		if err != nil {
			return fmt.Errorf("failed to resolve allowed directory: %w", err)
		}

		absFilePath, err := filepath.Abs(filePath)
		if err != nil {
			return fmt.Errorf("failed to resolve file path: %w", err)
		}

		if !strings.HasPrefix(absFilePath, absAllowedDir) {
			return fmt.Errorf("file path %s is outside allowed directory %s", absFilePath, absAllowedDir)
		}
	}

	// Check for dangerous characters
	dangerousChars := []string{"|", "&", ";", "`", "$", "(", ")", "{", "}", "[", "]", "*", "?", "\\"}
	for _, char := range dangerousChars {
		if strings.Contains(filePath, char) {
			return fmt.Errorf("dangerous character '%s' not allowed in path: %s", char, filePath)
		}
	}

	return nil
}

// validateFFmpegPath ensures the FFmpeg executable path is secure
func (p *Processor) validateFFmpegPath() error {
	if p.ffmpegPath == "" {
		return fmt.Errorf("FFmpeg path cannot be empty")
	}

	// Check for path traversal attempts
	if strings.Contains(p.ffmpegPath, "..") {
		return fmt.Errorf("path traversal not allowed in FFmpeg path: %s", p.ffmpegPath)
	}

	// Check for dangerous characters
	dangerousChars := []string{"|", "&", ";", "`", "$", "(", ")", "{", "}", "[", "]", "*", "?", "\\"}
	for _, char := range dangerousChars {
		if strings.Contains(p.ffmpegPath, char) {
			return fmt.Errorf("dangerous character '%s' not allowed in FFmpeg path: %s", char, p.ffmpegPath)
		}
	}

	// Check if FFmpeg executable exists and is executable
	if filepath.IsAbs(p.ffmpegPath) {
		if _, err := os.Stat(p.ffmpegPath); os.IsNotExist(err) {
			return fmt.Errorf("FFmpeg executable not found: %s", p.ffmpegPath)
		}
	} else {
		// For relative paths, check if it's in PATH
		if _, err := exec.LookPath(p.ffmpegPath); err != nil {
			return fmt.Errorf("FFmpeg executable not found in PATH: %s", p.ffmpegPath)
		}
	}

	return nil
}

// executeCommand executes a command with security validation
func (p *Processor) executeCommand(args []string) error {
	// Validate FFmpeg path
	if err := p.validateFFmpegPath(); err != nil {
		return fmt.Errorf("FFmpeg path validation failed: %w", err)
	}

	// Validate all arguments that are file paths
	for i, arg := range args {
		// Skip FFmpeg options that start with -
		if strings.HasPrefix(arg, "-") {
			continue
		}
		// Validate file paths
		if err := p.validatePath(arg); err != nil {
			return fmt.Errorf("invalid argument %d: %w", i, err)
		}
	}

	// Execute command with validated arguments
	// Note: All arguments have been validated above, so this is safe
	cmd := exec.Command(p.ffmpegPath, args...) //nolint:gosec // All arguments validated above
	return cmd.Run()
}

// executeCommandWithOutput executes a command with security validation and returns output
func (p *Processor) executeCommandWithOutput(args []string) (string, error) {
	// Validate FFmpeg path
	if err := p.validateFFmpegPath(); err != nil {
		return "", fmt.Errorf("FFmpeg path validation failed: %w", err)
	}

	// Validate all arguments that are file paths
	for i, arg := range args {
		// Skip FFmpeg options that start with -
		if strings.HasPrefix(arg, "-") {
			continue
		}
		// Validate file paths
		if err := p.validatePath(arg); err != nil {
			return "", fmt.Errorf("invalid argument %d: %w", i, err)
		}
	}

	// Execute command with validated arguments
	// Note: All arguments have been validated above, so this is safe
	cmd := exec.Command(p.ffmpegPath, args...) //nolint:gosec // All arguments validated above
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stderr.String(), err
}

// GetVideoInfo retrieves metadata about a video file
func (p *Processor) GetVideoInfo(inputPath string) (*VideoInfo, error) {
	// Validate input path
	if err := p.validatePath(inputPath); err != nil {
		return nil, fmt.Errorf("invalid input path: %w", err)
	}

	output, err := p.executeCommandWithOutput([]string{"-i", inputPath, "-f", "null", "-"})
	if err != nil {
		return nil, fmt.Errorf("failed to get video info: %w", err)
	}

	// Parse FFmpeg output
	info := &VideoInfo{}

	// Extract duration
	if durationStr := extractValue(output, "Duration: "); durationStr != "" {
		parts := strings.Split(durationStr, ":")
		if len(parts) == 3 {
			hours, err := strconv.ParseFloat(parts[0], 64)
			if err != nil {
				hours = 0
			}
			minutes, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				minutes = 0
			}
			seconds, err := strconv.ParseFloat(parts[2], 64)
			if err != nil {
				seconds = 0
			}
			info.Duration = hours*3600 + minutes*60 + seconds
		}
	}

	// Extract resolution
	if resStr := extractValue(output, "Stream #0:0"); resStr != "" {
		if strings.Contains(resStr, "Video:") {
			if res := extractValue(resStr, ", "); res != "" {
				if _, err := fmt.Sscanf(res, "%dx%d", &info.Width, &info.Height); err != nil {
					// If parsing fails, set default values
					info.Width = 0
					info.Height = 0
				}
			}
		}
	}

	// Extract codecs
	if videoStr := extractValue(output, "Video: "); videoStr != "" {
		parts := strings.Split(videoStr, ",")
		if len(parts) > 0 {
			info.VideoCodec = strings.TrimSpace(parts[0])
		}
	}
	if audioStr := extractValue(output, "Audio: "); audioStr != "" {
		parts := strings.Split(audioStr, ",")
		if len(parts) > 0 {
			info.AudioCodec = strings.TrimSpace(parts[0])
		}
	}

	return info, nil
}

// GenerateThumbnail creates a thumbnail image from a video file
func (p *Processor) GenerateThumbnail(inputPath, outputPath string, timestamp float64) error {
	// Validate input and output paths
	if err := p.validatePath(inputPath); err != nil {
		return fmt.Errorf("invalid input path: %w", err)
	}
	if err := p.validatePath(outputPath); err != nil {
		return fmt.Errorf("invalid output path: %w", err)
	}

	// Validate timestamp
	if timestamp < 0 {
		return fmt.Errorf("timestamp must be non-negative")
	}

	return p.executeCommand([]string{
		"-i", inputPath,
		"-ss", fmt.Sprintf("%.2f", timestamp),
		"-vframes", "1",
		"-q:v", "2",
		outputPath,
	})
}

// TrimVideo trims a video to the specified start and end times
func (p *Processor) TrimVideo(inputPath, outputPath string, startTime, endTime float64) error {
	// Validate input and output paths
	if err := p.validatePath(inputPath); err != nil {
		return fmt.Errorf("invalid input path: %w", err)
	}
	if err := p.validatePath(outputPath); err != nil {
		return fmt.Errorf("invalid output path: %w", err)
	}

	// Validate timestamps
	if startTime < 0 || endTime < 0 {
		return fmt.Errorf("timestamps must be non-negative")
	}
	if startTime >= endTime {
		return fmt.Errorf("start time must be less than end time")
	}

	return p.executeCommand([]string{
		"-i", inputPath,
		"-ss", fmt.Sprintf("%.2f", startTime),
		"-to", fmt.Sprintf("%.2f", endTime),
		"-c", "copy",
		outputPath,
	})
}

// ExtractValue extracts a value from FFmpeg output using a prefix
func extractValue(output, prefix string) string {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, prefix) {
			parts := strings.SplitN(line, prefix, 2)
			if len(parts) > 1 {
				return strings.TrimSpace(parts[1])
			}
		}
	}
	return ""
}

// AddWatermark adds a watermark to the video
func (p *Processor) AddWatermark(inputPath, outputPath, watermarkPath string) error {
	// Validate all paths
	if err := p.validatePath(inputPath); err != nil {
		return fmt.Errorf("invalid input path: %w", err)
	}
	if err := p.validatePath(outputPath); err != nil {
		return fmt.Errorf("invalid output path: %w", err)
	}
	if err := p.validatePath(watermarkPath); err != nil {
		return fmt.Errorf("invalid watermark path: %w", err)
	}

	return p.executeCommand([]string{
		"-i", inputPath,
		"-i", watermarkPath,
		"-filter_complex", "overlay=10:10",
		"-c:a", "copy",
		outputPath,
	})
}

// ConvertToMP4 converts a video to MP4 format
func (p *Processor) ConvertToMP4(inputPath, outputPath string) error {
	// Validate input and output paths
	if err := p.validatePath(inputPath); err != nil {
		return fmt.Errorf("invalid input path: %w", err)
	}
	if err := p.validatePath(outputPath); err != nil {
		return fmt.Errorf("invalid output path: %w", err)
	}

	return p.executeCommand([]string{
		"-i", inputPath,
		"-c:v", "libx264",
		"-c:a", "aac",
		"-strict", "experimental",
		outputPath,
	})
}

// ValidateVideo validates that a file is a valid video file
func (p *Processor) ValidateVideo(inputPath string) error {
	// Validate input path
	if err := p.validatePath(inputPath); err != nil {
		return fmt.Errorf("invalid input path: %w", err)
	}

	// Check if file exists
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return fmt.Errorf("video file does not exist: %s", inputPath)
	}

	// Try to get video info to validate it's a proper video file
	_, err := p.GetVideoInfo(inputPath)
	if err != nil {
		return fmt.Errorf("invalid video file: %w", err)
	}

	return nil
}

// EnsureDirectory ensures the output directory exists
func (p *Processor) EnsureDirectory(path string) error {
	return os.MkdirAll(filepath.Dir(path), 0750)
}

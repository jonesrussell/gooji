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

// GetVideoInfo retrieves metadata about a video file
func (p *Processor) GetVideoInfo(inputPath string) (*VideoInfo, error) {
	cmd := exec.Command(p.ffmpegPath, "-i", inputPath, "-f", "null", "-")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to get video info: %v", err)
	}

	// Parse FFmpeg output
	output := stderr.String()
	info := &VideoInfo{}

	// Extract duration
	if durationStr := extractValue(output, "Duration: "); durationStr != "" {
		parts := strings.Split(durationStr, ":")
		if len(parts) == 3 {
			hours, _ := strconv.ParseFloat(parts[0], 64)
			minutes, _ := strconv.ParseFloat(parts[1], 64)
			seconds, _ := strconv.ParseFloat(parts[2], 64)
			info.Duration = hours*3600 + minutes*60 + seconds
		}
	}

	// Extract resolution
	if resStr := extractValue(output, "Stream #0:0"); resStr != "" {
		if strings.Contains(resStr, "Video:") {
			if res := extractValue(resStr, ", "); res != "" {
				fmt.Sscanf(res, "%dx%d", &info.Width, &info.Height)
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
	cmd := exec.Command(p.ffmpegPath,
		"-i", inputPath,
		"-ss", fmt.Sprintf("%.2f", timestamp),
		"-vframes", "1",
		"-q:v", "2",
		outputPath,
	)
	return cmd.Run()
}

// TrimVideo trims a video to the specified start and end times
func (p *Processor) TrimVideo(inputPath, outputPath string, startTime, endTime float64) error {
	cmd := exec.Command(p.ffmpegPath,
		"-i", inputPath,
		"-ss", fmt.Sprintf("%.2f", startTime),
		"-to", fmt.Sprintf("%.2f", endTime),
		"-c", "copy",
		outputPath,
	)
	return cmd.Run()
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
	cmd := exec.Command(p.ffmpegPath,
		"-i", inputPath,
		"-i", watermarkPath,
		"-filter_complex", "overlay=10:10",
		"-c:a", "copy",
		outputPath,
	)
	return cmd.Run()
}

// ConvertToMP4 converts a video to MP4 format
func (p *Processor) ConvertToMP4(inputPath, outputPath string) error {
	cmd := exec.Command(p.ffmpegPath,
		"-i", inputPath,
		"-c:v", "libx264",
		"-c:a", "aac",
		"-strict", "experimental",
		outputPath,
	)
	return cmd.Run()
}

// EnsureDirectory ensures the output directory exists
func (p *Processor) EnsureDirectory(path string) error {
	return os.MkdirAll(filepath.Dir(path), 0755)
}

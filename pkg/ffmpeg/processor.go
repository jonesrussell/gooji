package ffmpeg

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Processor handles video processing operations using FFmpeg
type Processor struct {
	ffmpegPath string
}

// NewProcessor creates a new video processor
func NewProcessor(ffmpegPath string) *Processor {
	return &Processor{
		ffmpegPath: ffmpegPath,
	}
}

// TrimVideo trims a video to the specified duration
func (p *Processor) TrimVideo(inputPath, outputPath string, duration int) error {
	cmd := exec.Command(p.ffmpegPath,
		"-i", inputPath,
		"-t", fmt.Sprintf("%d", duration),
		"-c", "copy",
		outputPath,
	)
	return cmd.Run()
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

// GenerateThumbnail generates a thumbnail from a video
func (p *Processor) GenerateThumbnail(inputPath, outputPath string, timestamp int) error {
	cmd := exec.Command(p.ffmpegPath,
		"-i", inputPath,
		"-ss", fmt.Sprintf("%d", timestamp),
		"-vframes", "1",
		"-q:v", "2",
		outputPath,
	)
	return cmd.Run()
}

// EnsureDirectory ensures the output directory exists
func (p *Processor) EnsureDirectory(path string) error {
	return os.MkdirAll(filepath.Dir(path), 0755)
}

// GetVideoInfo returns basic information about a video file
func (p *Processor) GetVideoInfo(inputPath string) (map[string]string, error) {
	cmd := exec.Command(p.ffmpegPath,
		"-i", inputPath,
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		// FFmpeg returns error when getting info, but we can still parse the output
		if len(output) == 0 {
			return nil, err
		}
	}

	// Parse the output to extract video information
	// This is a basic implementation - you might want to enhance it
	info := make(map[string]string)
	info["raw_output"] = string(output)
	return info, nil
}

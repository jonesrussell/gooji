package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Storage holds storage configuration
type Storage struct {
	BasePath   string `json:"base_path"`
	Uploads    string `json:"uploads"`
	Temp       string `json:"temp"`
	Logs       string `json:"logs"`
	Thumbnails string `json:"thumbnails"`
	Metadata   string `json:"metadata"`
}

// Config holds the application configuration
type Config struct {
	Server struct {
		Port int `json:"port"`
	} `json:"server"`
	Storage Storage `json:"storage"`
	Video   struct {
		MaxSize      int64    `json:"max_size"`
		AllowedTypes []string `json:"allowed_types"`
	} `json:"video"`
	FFmpeg struct {
		Path string `json:"path"`
	} `json:"ffmpeg"`
}

// validatePath ensures a file path is secure
func validatePath(filePath string) error {
	if filePath == "" {
		return fmt.Errorf("file path cannot be empty")
	}

	// Check for path traversal attempts
	if strings.Contains(filePath, "..") {
		return fmt.Errorf("path traversal not allowed: %s", filePath)
	}

	// Check for dangerous characters
	dangerousChars := []string{"|", "&", ";", "`", "$", "(", ")", "{", "}", "[", "]", "*", "?", "\\"}
	for _, char := range dangerousChars {
		if strings.Contains(filePath, char) {
			return fmt.Errorf("dangerous character '%s' not allowed in path: %s", char, filePath)
		}
	}

	// For config files, only allow relative paths or paths within current directory
	if filepath.IsAbs(filePath) {
		currentDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}

		absFilePath, err := filepath.Abs(filePath)
		if err != nil {
			return fmt.Errorf("failed to resolve file path: %w", err)
		}

		if !strings.HasPrefix(absFilePath, currentDir) {
			return fmt.Errorf("config file path %s is outside current directory %s", absFilePath, currentDir)
		}
	}

	return nil
}

// Load reads the configuration from a JSON file and environment variables
func Load(path string) (*Config, error) {
	// Validate config path
	if err := validatePath(path); err != nil {
		return nil, fmt.Errorf("invalid config path: %w", err)
	}

	// Additional security check: ensure path is absolute and no path traversal
	if !filepath.IsAbs(path) {
		return nil, fmt.Errorf("config path must be absolute: %s", path)
	}
	if strings.Contains(path, "..") {
		return nil, fmt.Errorf("path traversal not allowed in config path: %s", path)
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	// Set defaults if not specified
	if config.Server.Port == 0 {
		config.Server.Port = 8080
	}
	if config.Storage.BasePath == "" {
		config.Storage.BasePath = "storage"
	}

	// Override with environment variables if set
	if port := os.Getenv("GOOJI_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			config.Server.Port = p
		}
	}

	// Note: We can't use the logger here as it's not available yet
	// Environment variable debug info will be logged by the main logger
	if config.Storage.Uploads == "" {
		config.Storage.Uploads = "storage/uploads"
	}
	if config.Storage.Temp == "" {
		config.Storage.Temp = "storage/temp"
	}
	if config.Storage.Logs == "" {
		config.Storage.Logs = "storage/logs"
	}
	if config.Storage.Thumbnails == "" {
		config.Storage.Thumbnails = "storage/thumbnails"
	}
	if config.Storage.Metadata == "" {
		config.Storage.Metadata = "storage/metadata"
	}
	if config.Video.MaxSize == 0 {
		config.Video.MaxSize = 100 * 1024 * 1024 // 100MB
	}
	if len(config.Video.AllowedTypes) == 0 {
		config.Video.AllowedTypes = []string{"video/mp4", "video/webm"}
	}
	if config.FFmpeg.Path == "" {
		config.FFmpeg.Path = "ffmpeg"
	}

	return &config, nil
}

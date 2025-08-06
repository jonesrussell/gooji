package config

import (
	"encoding/json"
	"fmt"
	"os"
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
	Video struct {
		MaxSize      int64    `json:"max_size"`
		AllowedTypes []string `json:"allowed_types"`
	} `json:"video"`
	FFmpeg struct {
		Path string `json:"path"`
	} `json:"ffmpeg"`
}

// Load reads the configuration from a JSON file
func Load(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %v", err)
	}
	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %v", err)
	}

	// Set defaults if not specified
	if config.Server.Port == 0 {
		config.Server.Port = 8080
	}
	if config.Storage.BasePath == "" {
		config.Storage.BasePath = "storage"
	}
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

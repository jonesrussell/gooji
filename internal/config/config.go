package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config holds the application configuration
type Config struct {
	Server struct {
		Port int `json:"port"`
	} `json:"server"`
	Video struct {
		StoragePath  string   `json:"storage_path"`
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
	if config.Video.StoragePath == "" {
		config.Video.StoragePath = "videos"
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

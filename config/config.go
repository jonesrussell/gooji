package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config holds the application configuration
type Config struct {
	Port           string   `json:"port"`
	MaxVideoLength int      `json:"maxVideoLength"` // in seconds
	VideoDirectory string   `json:"videoDirectory"`
	AllowedOrigins []string `json:"allowedOrigins"`
	FFmpegPath     string   `json:"ffmpegPath"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Port:           "8080",
		MaxVideoLength: 300, // 5 minutes
		VideoDirectory: "videos",
		AllowedOrigins: []string{"*"},
		FFmpegPath:     "ffmpeg",
	}
}

// Load loads the configuration from a file
func Load(path string) (*Config, error) {
	config := DefaultConfig()

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, err
	}

	// Try to load existing config
	if _, err := os.Stat(path); err == nil {
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		if err := json.NewDecoder(file).Decode(config); err != nil {
			return nil, err
		}
	} else {
		// Create default config file
		file, err := os.Create(path)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		if err := json.NewEncoder(file).Encode(config); err != nil {
			return nil, err
		}
	}

	return config, nil
}

// Save saves the configuration to a file
func (c *Config) Save(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(c)
}

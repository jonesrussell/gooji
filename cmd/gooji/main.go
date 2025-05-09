package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"gooji/config"
	"gooji/internal/video"
	"gooji/pkg/ffmpeg"
)

func main() {
	// Load configuration
	cfg, err := config.Load("config/config.json")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create video processor
	processor := ffmpeg.NewProcessor(cfg.FFmpegPath)

	// Create video handler
	videoHandler, err := video.NewHandler(processor, cfg.VideoDirectory)
	if err != nil {
		log.Fatalf("Failed to create video handler: %v", err)
	}

	// Ensure required directories exist
	dirs := []string{
		"web/static",
		"web/templates",
		cfg.VideoDirectory,
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}

	// Serve static files
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// API routes
	http.HandleFunc("/api/videos/upload", videoHandler.HandleUpload)
	http.HandleFunc("/api/videos", videoHandler.ListVideos)
	http.HandleFunc("/api/videos/", videoHandler.GetVideo)
	http.HandleFunc("/api/videos/thumbnail/", videoHandler.GetThumbnail)

	// Main page handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join("web/templates", "index.html"))
	})

	// Start server
	log.Printf("Starting server on port %s...", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"gooji/internal/handler"
	"gooji/internal/repository/sqlite"
	"gooji/internal/service"
)

func main() {
	// Configuration
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "gooji.db"
	}

	// Database
	db, err := sqlite.NewDB(dbPath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Repositories
	videoRepo := sqlite.NewVideoRepository(db)

	// Services
	videoService := service.NewVideoService(videoRepo)

	// Handlers
	version := "1.0.0" // TODO: Read from build info or manifest
	h := handler.New(videoService, version)

	// Routes
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", h.Home)
	mux.HandleFunc("GET /gallery", h.Gallery)
	mux.HandleFunc("GET /health", h.Health)

	// Static files (Vite build output)
	staticDir := filepath.Join("web", "dist", "assets")
	mux.Handle("GET /assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(staticDir))))

	// Start server
	addr := ":" + port
	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

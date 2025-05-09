package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gooji/internal/config"
	"gooji/internal/logger"
	"gooji/internal/middleware"
	"gooji/internal/video"
)

func main() {
	// Load configuration
	cfg, err := config.Load("config/config.json")
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	log, err := logger.New("logs")
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Close()

	// Create video processor
	processor, err := video.NewProcessor(cfg.Video)
	if err != nil {
		log.Fatal("Failed to create video processor: %v", err)
	}

	// Create video handler
	handler := video.NewHandler(processor, log)

	// Create router
	mux := http.NewServeMux()

	// Serve static files
	fs := http.FileServer(http.Dir("web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// API routes
	mux.HandleFunc("/api/videos", handler.HandleVideos)
	mux.HandleFunc("/api/videos/", handler.HandleVideo)

	// Page routes
	mux.HandleFunc("/", handler.HandleHome)
	mux.HandleFunc("/record", handler.HandleRecord)
	mux.HandleFunc("/edit/", handler.HandleEdit)
	mux.HandleFunc("/gallery", handler.HandleGallery)

	// Create server with middleware
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: middleware.Chain(mux, middleware.Logging(log), middleware.Recovery(log), middleware.CORS()),
	}

	// Start server in a goroutine
	go func() {
		log.Info("Starting server on port %d", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	log.Info("Shutting down server...")
	if err := server.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown: %v", err)
	}

	log.Info("Server stopped")
}

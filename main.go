package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"gooji/internal/config"
	"gooji/internal/handler"
	"gooji/internal/middleware"
	"gooji/internal/repository/filesystem"
	"gooji/internal/repository/sqlite"
	"gooji/internal/service"
	"gooji/pkg/ffmpeg"
)

func main() {
	// Initialize logger
	logger, _ := zap.NewProduction()
	if os.Getenv("APP_ENV") == "development" {
		logger, _ = zap.NewDevelopment()
	}
	defer logger.Sync()

	// Load configuration
	cfg, err := config.Load("config/config.json")
	if err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
	}

	// Initialize database
	db, err := sqlite.NewDB("storage/gooji.db")
	if err != nil {
		logger.Fatal("failed to open database", zap.Error(err))
	}
	defer db.Close()

	if err := sqlite.RunMigrations(db, "migrations"); err != nil {
		logger.Fatal("failed to run migrations", zap.Error(err))
	}

	// Initialize repositories
	videoRepo := sqlite.NewVideoRepository(db)
	fileStorage := filesystem.NewFileStorage(cfg.Storage.Uploads, cfg.Storage.Thumbnails)

	// Initialize FFmpeg processor
	processor := ffmpeg.NewProcessorWithSecurity(cfg.FFmpeg.Path, cfg.Storage.BasePath)

	// Initialize services
	videoService := service.NewVideoService(videoRepo, fileStorage, processor, logger)

	// Initialize Inertia
	isDev := os.Getenv("APP_ENV") == "development" || os.Getenv("VITE_DEV") == "true"
	inertia, err := handler.NewInertia(isDev)
	if err != nil {
		logger.Fatal("failed to create inertia", zap.Error(err))
	}

	// Initialize handlers
	pageHandler := handler.NewPageHandler(inertia, videoService, logger)
	videoHandler := handler.NewVideoHandler(inertia, videoService, logger, cfg.Storage.BasePath)

	// Set up routes
	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("POST /api/videos", videoHandler.Upload)
	mux.HandleFunc("GET /api/videos", videoHandler.List)
	mux.HandleFunc("GET /api/videos/{id}", videoHandler.Get)
	mux.HandleFunc("DELETE /api/videos/{id}", videoHandler.Delete)
	mux.HandleFunc("GET /api/thumbnails", videoHandler.Thumbnail)

	// Page routes
	mux.HandleFunc("GET /", pageHandler.Home)
	mux.HandleFunc("GET /gallery", pageHandler.Gallery)
	mux.HandleFunc("GET /upload", pageHandler.Upload)
	mux.HandleFunc("GET /record", pageHandler.Record)
	mux.HandleFunc("GET /video/{id}", pageHandler.VideoShow)
	mux.HandleFunc("GET /edit/{id}", pageHandler.Edit)

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	})

	// Middleware
	var h http.Handler = mux
	h = middleware.Recovery(logger)(h)
	h = middleware.Logging(logger)(h)
	h = middleware.CORS()(h)
	h = inertia.Middleware(h)

	// Start server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: h,
	}

	go func() {
		logger.Info("starting server", zap.Int("port", cfg.Server.Port))
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			logger.Fatal("server error", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.Info("shutting down server")
	server.Shutdown(ctx)
}

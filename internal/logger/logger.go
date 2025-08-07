package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger handles application logging
type Logger struct {
	zap *zap.Logger
}

// New creates a new logger
func New(logDir string) (*Logger, error) {
	if err := os.MkdirAll(logDir, 0750); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	// Create log file
	logFile := filepath.Join(logDir, fmt.Sprintf("gooji_%s.log", time.Now().Format("2006-01-02")))
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	// Determine log level from environment
	logLevel := zapcore.InfoLevel
	if os.Getenv("APP_DEBUG") == "true" {
		logLevel = zapcore.DebugLevel
	}

	// Configure Zap encoder
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Create core with both file and console output
	core := zapcore.NewTee(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(file),
			logLevel,
		),
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			logLevel,
		),
	)

	// Create logger
	zapLogger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return &Logger{
		zap: zapLogger,
	}, nil
}

// Close closes the logger
func (l *Logger) Close() error {
	return l.zap.Sync()
}

// Debug logs a debug message
func (l *Logger) Debug(format string, args ...interface{}) {
	l.zap.Sugar().Debugf(format, args...)
}

// Info logs an info message
func (l *Logger) Info(format string, args ...interface{}) {
	l.zap.Sugar().Infof(format, args...)
}

// Error logs an error message
func (l *Logger) Error(format string, args ...interface{}) {
	l.zap.Sugar().Errorf(format, args...)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.zap.Sugar().Fatalf(format, args...)
}

#!/bin/bash
set -e

# Download dependencies
go mod download
go mod tidy

# Create required directories
mkdir -p logs
mkdir -p videos

# Check FFmpeg installation
if ! command -v ffmpeg &> /dev/null; then
    echo "FFmpeg is required but not installed. Please install FFmpeg version 6.1.1 or later."
    exit 1
fi

# Check golangci-lint installation
if ! command -v golangci-lint &> /dev/null; then
    echo "Installing golangci-lint..."
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    echo "Please add $(go env GOPATH)/bin to your PATH"
    echo "export PATH=$PATH:$(go env GOPATH)/bin"
    exit 1
fi 
#!/bin/bash
set -e

# Run go vet
go vet ./...

# Check golangci-lint installation
if ! command -v golangci-lint &> /dev/null; then
    echo "golangci-lint not found. Please run 'task setup' first"
    exit 1
fi

# Run golangci-lint
golangci-lint run 

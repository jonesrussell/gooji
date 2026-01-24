#!/bin/bash
set -e

echo "Cleaning and tidying dependencies..."
go mod tidy

echo "Downloading dependencies..."
if ! go mod download -x; then
    echo "Error: Failed to download dependencies"
    exit 1
fi

echo "Verifying dependencies..."
if ! go mod verify; then
    echo "Error: Dependency verification failed"
    exit 1
fi

echo "✅ Dependencies verified and go.sum is up to date" 
#!/bin/bash
set -e

echo "🔍 Running Go linting tools..."

# Run go vet
echo "📋 Running go vet..."
go vet ./...

# Run staticcheck
echo "🔍 Running staticcheck..."
go run honnef.co/go/tools/cmd/staticcheck@latest ./...

# Run govulncheck
echo "🛡️  Running security vulnerability check..."
go run golang.org/x/vuln/cmd/govulncheck@latest ./...

# Run golangci-lint if available (legacy)
if command -v golangci-lint &> /dev/null; then
    echo "🔧 Running golangci-lint (legacy)..."
    golangci-lint run
else
    echo "ℹ️  golangci-lint not found (optional legacy tool)"
fi

echo "✅ Linting complete!" 

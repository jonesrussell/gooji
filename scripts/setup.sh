#!/bin/bash
set -e

echo "🚀 Setting up Gooji development environment..."

# Download and tidy dependencies
echo "📦 Downloading dependencies..."
go mod download
go mod tidy

# Create required directories
echo "📁 Creating directories..."
mkdir -p logs
mkdir -p videos
mkdir -p bin

# Check Go version
echo "🔍 Checking Go version..."
go version

# Check FFmpeg installation
echo "🎬 Checking FFmpeg installation..."
if ! command -v ffmpeg &> /dev/null; then
    echo "❌ FFmpeg is required but not installed."
    echo "   Please install FFmpeg version 6.1.1 or later:"
    echo "   Ubuntu/Debian: sudo apt install ffmpeg"
    echo "   macOS: brew install ffmpeg"
    exit 1
else
    echo "✅ FFmpeg found: $(ffmpeg -version | head -1)"
fi

# Install development tools
echo "🛠️  Installing development tools..."
go install github.com/air-verse/air@latest
go install honnef.co/go/tools/cmd/staticcheck@latest
go install golang.org/x/vuln/cmd/govulncheck@latest

# Create default config if it doesn't exist
if [ ! -f config/config.json ]; then
    echo "⚙️  Creating default configuration..."
    ./scripts/create-config.sh
fi

# Run initial checks
echo "🔍 Running initial checks..."
go vet ./...
go run honnef.co/go/tools/cmd/staticcheck@latest ./... || echo "⚠️  Staticcheck found issues (this is normal for new projects)"

echo "✅ Setup complete! You can now run:"
echo "   task dev    - Start development server with hot reload"
echo "   task run    - Start production server"
echo "   task test   - Run tests"
echo "   task lint   - Run linting tools" 
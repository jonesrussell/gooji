#!/bin/bash
set -e

echo "🔍 Checking Gooji development environment..."

# Check Go installation
echo "📋 Checking Go installation..."
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed"
    exit 1
fi
echo "✅ Go found: $(go version)"

# Check Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
REQUIRED_VERSION="1.24"
if [ "$(printf '%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$REQUIRED_VERSION" ]; then
    echo "⚠️  Go version $GO_VERSION is older than required $REQUIRED_VERSION"
else
    echo "✅ Go version $GO_VERSION meets requirements"
fi

# Check FFmpeg installation
echo "🎬 Checking FFmpeg installation..."
if ! command -v ffmpeg &> /dev/null; then
    echo "❌ FFmpeg is not installed"
    exit 1
fi
echo "✅ FFmpeg found: $(ffmpeg -version | head -1)"

# Check development tools
echo "🛠️  Checking development tools..."

# Check Air
if go list -m github.com/air-verse/air &> /dev/null; then
    echo "✅ Air is available in go.mod"
else
    echo "⚠️  Air not found in go.mod"
fi

# Check if tools can be run
echo "🔧 Testing tool availability..."
if go run github.com/air-verse/air@latest -v &> /dev/null; then
    echo "✅ Air can be executed"
else
    echo "❌ Air cannot be executed"
fi

if go run honnef.co/go/tools/cmd/staticcheck@latest -version &> /dev/null; then
    echo "✅ Staticcheck can be executed"
else
    echo "❌ Staticcheck cannot be executed"
fi

if go run golang.org/x/vuln/cmd/govulncheck@latest -version &> /dev/null; then
    echo "✅ Govulncheck can be executed"
else
    echo "❌ Govulncheck cannot be executed"
fi

# Check required directories
echo "📁 Checking required directories..."
for dir in logs videos bin; do
    if [ -d "$dir" ]; then
        echo "✅ Directory $dir exists"
    else
        echo "⚠️  Directory $dir missing (will be created by setup)"
    fi
done

# Check configuration
echo "⚙️  Checking configuration..."
if [ -f "config/config.json" ]; then
    echo "✅ Configuration file exists"
else
    echo "⚠️  Configuration file missing (run 'task create-config')"
fi

# Check dependencies
echo "📦 Checking dependencies..."
if go mod verify &> /dev/null; then
    echo "✅ Dependencies are valid"
else
    echo "❌ Dependencies are invalid (run 'task verify-deps')"
fi

echo "✅ Environment check complete!" 
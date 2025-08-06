# Gooji Development Guide

## 🚀 **Quick Start**

```bash
# Set up development environment
task setup

# Start development server with hot reload
task dev

# Run tests
task test

# Run linting tools
task lint
```

## 📋 **Available Tasks**

### **Development**
- `task dev` - Start development server with hot reload (Air)
- `task run` - Start production server
- `task build` - Build the application
- `task test` - Run tests with coverage
- `task test-coverage` - Generate coverage report

### **Code Quality**
- `task lint` - Run modern Go linting tools (go vet, staticcheck, govulncheck)
- `task lint-legacy` - Run legacy linters (golangci-lint if available)
- `task format` - Format Go code
- `task security-scan` - Run security vulnerability scan

### **Environment Management**
- `task setup` - Set up development environment
- `task check-env` - Check development environment
- `task check-ffmpeg` - Check FFmpeg installation
- `task install-tools` - Install development tools

### **Dependencies**
- `task verify-deps` - Verify dependencies
- `task update-deps` - Update dependencies
- `task tidy` - Tidy Go modules
- `task check-deps` - Check for outdated dependencies

### **Configuration**
- `task create-config` - Create default configuration
- `task clean` - Clean build artifacts

### **Docker**
- `task docker-build` - Build Docker image
- `task docker-run` - Run Docker container

## 🛠️ **Development Tools**

### **Modern Go Tooling**
The project now uses modern Go tools instead of external dependencies:

- **go vet** - Built-in Go static analysis
- **staticcheck** - Advanced static analysis
- **govulncheck** - Security vulnerability scanning
- **Air** - Hot reload for development

### **Tool Installation**
All tools are installed via `go install` and run with `go run`:

```bash
# Install tools
task install-tools

# Run tools directly
go run github.com/air-verse/air@latest
go run honnef.co/go/tools/cmd/staticcheck@latest ./...
go run golang.org/x/vuln/cmd/govulncheck@latest ./...
```

## 🔧 **Scripts Overview**

### **Updated Scripts**

#### `scripts/setup.sh`
- ✅ Enhanced with better error messages and emojis
- ✅ Installs modern Go tools automatically
- ✅ Creates required directories
- ✅ Runs initial checks
- ✅ Provides helpful output

#### `scripts/lint.sh`
- ✅ Uses modern Go tools (staticcheck, govulncheck)
- ✅ Falls back to golangci-lint if available
- ✅ Better output formatting

#### `scripts/check-env.sh` (New)
- ✅ Comprehensive environment validation
- ✅ Checks Go version, FFmpeg, tools, directories
- ✅ Provides actionable feedback

#### `scripts/check-ffmpeg.sh`
- ✅ Simple FFmpeg availability check
- ✅ Shows version information

#### `scripts/create-config.sh`
- ✅ Creates default configuration
- ✅ Includes all necessary settings

#### `scripts/verify-deps.sh`
- ✅ Verifies Go module dependencies
- ✅ Ensures go.sum is up to date

## 🎯 **Key Improvements**

### **1. Modern Go Tooling**
- Replaced external golangci-lint with built-in Go tools
- Added staticcheck for advanced static analysis
- Added govulncheck for security scanning
- All tools run via `go run` (no global installation needed)

### **2. Better Development Experience**
- `task dev` now properly uses Air for hot reload
- `task run` uses `go run` for faster startup
- Enhanced error messages and feedback
- Comprehensive environment checking

### **3. Improved Taskfile**
- Removed duplicate tasks
- Added new useful tasks (check-env, security-scan, format)
- Better task organization and descriptions
- Fixed syntax issues

### **4. Enhanced Scripts**
- Better error handling and user feedback
- Emoji-based output for better readability
- Comprehensive environment validation
- Automatic tool installation

## 🔍 **Environment Requirements**

### **Required**
- Go 1.24 or later
- FFmpeg 6.1.1 or later
- Task (task runner)

### **Optional**
- golangci-lint (for legacy linting)
- Docker (for containerization)

## 🚨 **Troubleshooting**

### **Air Not Working**
```bash
# Check if Air is available
task check-env

# Reinstall Air
go install github.com/air-verse/air@latest

# Run Air directly
go run github.com/air-verse/air@latest
```

### **FFmpeg Issues**
```bash
# Check FFmpeg installation
task check-ffmpeg

# Install FFmpeg (Ubuntu/Debian)
sudo apt update && sudo apt install ffmpeg

# Install FFmpeg (macOS)
brew install ffmpeg
```

### **Linting Issues**
```bash
# Run modern linting
task lint

# Run legacy linting (if golangci-lint installed)
task lint-legacy

# Install golangci-lint (optional)
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

## 📊 **Best Practices**

### **Development Workflow**
1. **Setup**: Run `task setup` once
2. **Development**: Use `task dev` for hot reload
3. **Testing**: Run `task test` before commits
4. **Linting**: Run `task lint` for code quality
5. **Security**: Run `task security-scan` regularly

### **Code Quality**
- Run `task lint` before committing
- Fix all staticcheck warnings
- Address security vulnerabilities immediately
- Maintain good test coverage

### **Dependencies**
- Use `task update-deps` to update dependencies
- Run `task verify-deps` to ensure consistency
- Check for vulnerabilities with `task security-scan`

## 🔄 **Migration from Old Scripts**

### **What Changed**
- `golangci-lint` → `staticcheck` + `govulncheck`
- External tool installation → `go run` approach
- Basic error messages → Enhanced feedback
- Manual setup → Automated environment checking

### **Benefits**
- ✅ No global tool installation required
- ✅ Better security scanning
- ✅ Faster development startup
- ✅ More reliable tooling
- ✅ Better error messages

---

**Last Updated**: August 6, 2025
**Go Version**: 1.24+
**FFmpeg Version**: 6.1.1+ 
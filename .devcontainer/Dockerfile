ARG VARIANT=1-1-bookworm
FROM mcr.microsoft.com/devcontainers/go:${VARIANT}

# Install OS packages in a single layer to reduce image size
RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends ffmpeg \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Switch to vscode user for Go tool installations
USER vscode

# Install Go tools - using latest versions with error handling
RUN set -e && \
    echo "Installing golangci-lint..." && \
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest && \
    echo "Installing migrate..." && \
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest && \
    echo "Installing templ..." && \
    go install github.com/a-h/templ/cmd/templ@latest && \
    echo "Installing air..." && \
    go install github.com/air-verse/air@latest && \
    echo "All Go tools installed successfully"

# Clean Go module cache to reduce image size
RUN go clean -modcache

# Switch back to root for any remaining system-level operations
USER root

# Uncomment and customize as needed:
# RUN su vscode -c "source /usr/local/share/nvm/nvm.sh && npm install -g <your-package-here>" 2>&1

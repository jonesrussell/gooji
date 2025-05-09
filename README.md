# Gooji - Ojibwe Language & Culture Video Kiosk

Gooji (meaning "video" in Ojibwe) is a modern kiosk application designed to record and edit short videos focused on Ojibwe language and culture. Inspired by MuchMusic's Speaker's Corner, this application provides a user-friendly interface for community members to share their knowledge and stories.

## Features

- Video recording with webcam support
- Basic video editing capabilities
- User-friendly touch interface
- Cultural content organization
- Export and sharing options
- Ojibwe language integration

## Project Structure

```
gooji/
├── cmd/            # Application entry points
├── internal/       # Private application code
├── pkg/           # Public library code
├── web/           # Frontend assets and templates
└── config/        # Configuration files
```

## Requirements

- Go 1.21 or later
- FFmpeg for video processing
- Modern web browser with WebRTC support

## Getting Started

1. Install Go dependencies:
   ```bash
   go mod download
   ```

2. Install FFmpeg:
   ```bash
   # Ubuntu/Debian
   sudo apt-get install ffmpeg
   ```

3. Run the application:
   ```bash
   go run cmd/gooji/main.go
   ```

## Contributing

This project is open to contributions from the community, especially from Ojibwe language speakers and cultural knowledge keepers.

## License

MIT License - See LICENSE file for details 
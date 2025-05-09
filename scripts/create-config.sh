#!/bin/bash
set -e

if [ ! -f config/config.json ]; then
    mkdir -p config
    cat > config/config.json << EOF
{
    "server": {
        "port": 8080,
        "host": "localhost"
    },
    "video": {
        "storage_path": "videos",
        "max_size": 1073741824,
        "allowed_types": ["video/mp4", "video/webm"]
    },
    "ffmpeg": {
        "path": "ffmpeg"
    }
}
EOF
    echo "Created default config/config.json"
else
    echo "config/config.json already exists"
fi 
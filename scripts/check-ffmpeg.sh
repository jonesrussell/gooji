#!/bin/bash
set -e

if ! command -v ffmpeg &> /dev/null; then
    echo "FFmpeg is not installed"
    exit 1
fi

ffmpeg -version 
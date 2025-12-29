#!/bin/bash

# Build script for Password Wordlist Generator

set -e

echo "Building Password Wordlist Generator..."

# Ensure dependencies are downloaded
echo "Downloading dependencies..."
go mod download

# Build for current platform
echo "Building for $(go env GOOS)/$(go env GOARCH)..."
go build -o passgen .

echo "Build complete! Run ./passgen to start the GUI."

# Optional: Build for other platforms
if [ "$1" == "all" ]; then
    echo "Building for all platforms..."
    
    # Linux
    echo "Building for Linux..."
    GOOS=linux GOARCH=amd64 go build -o passgen-linux-amd64 .
    
    # Windows
    echo "Building for Windows..."
    GOOS=windows GOARCH=amd64 go build -o passgen-windows-amd64.exe .
    
    # macOS
    echo "Building for macOS..."
    GOOS=darwin GOARCH=amd64 go build -o passgen-darwin-amd64 .
    GOOS=darwin GOARCH=arm64 go build -o passgen-darwin-arm64 .
    
    echo "All builds complete!"
fi


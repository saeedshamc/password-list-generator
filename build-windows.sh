#!/bin/bash

# Build script for Windows executable with icon

set -e

echo "Building Windows executable with icon..."

# Check if icon exists
if [ ! -f "passgen.png" ]; then
    echo "Error: passgen.png not found!"
    exit 1
fi

# Ensure dependencies are downloaded
echo "Downloading dependencies..."
go mod download

# Generate versioninfo.syso with icon
echo "Generating Windows resource file with icon..."
goversioninfo -64 -o versioninfo.syso versioninfo.json

# Check for mingw-w64
if ! command -v x86_64-w64-mingw32-gcc &> /dev/null; then
    echo "Error: mingw-w64 is not installed!"
    echo "Please install it with: sudo apt-get install -y gcc-mingw-w64-x86-64"
    exit 1
fi

# Build for Windows with CGO
echo "Building for Windows (amd64) with CGO..."
CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o passgen.exe .

# Clean up
echo "Cleaning up..."
rm -f versioninfo.syso

echo "Build complete! Windows executable: passgen.exe"


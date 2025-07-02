#!/bin/bash

# SSO Web Application Startup Script

echo "Starting SSO Web Application..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Go is not installed. Please install Go first."
    exit 1
fi

# Check if .env file exists
if [ ! -f .env ]; then
    echo ".env file not found. Copying from .env.example..."
    cp .env.example .env
    echo "Please edit .env file with your configuration before running the app."
fi

# Install dependencies
echo "Installing dependencies..."
go mod tidy

# Build the application
echo "Building application..."
go build -o bin/sso-app cmd/server/main.go

if [ $? -eq 0 ]; then
    echo "Build successful!"
    echo "Starting server on http://localhost:8080"
    echo "Press Ctrl+C to stop the server"
    echo ""
    ./bin/sso-app
else
    echo "Build failed!"
    exit 1
fi

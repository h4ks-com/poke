#!/bin/bash

echo "Starting Viridian City Bank application..."

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker first."
    exit 1
fi

# Build and run using docker-compose
docker-compose up --build -d

if [ $? -eq 0 ]; then
    echo "✅ Application started successfully!"
    echo "🌐 Frontend available at: http://localhost:8080"
    echo "🔧 API available at: http://localhost:8080/api"
    echo "❤️  Health check: http://localhost:8080/health"
    echo ""
    echo "To view logs: docker-compose logs -f"
    echo "To stop: docker-compose down"
else
    echo "❌ Failed to start application!"
    exit 1
fi

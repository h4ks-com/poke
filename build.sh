#!/bin/bash

echo "Building Viridian City Bank Docker image..."
docker build -t viridian-bank:latest .

if [ $? -eq 0 ]; then
    echo "✅ Build completed successfully!"
    echo "To run the application:"
    echo "  docker run -p 8080:8080 viridian-bank:latest"
    echo "  OR"
    echo "  docker-compose up"
else
    echo "❌ Build failed!"
    exit 1
fi

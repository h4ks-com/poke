# Multi-stage build for Viridian City Bank
FROM golang:1.20 AS backend-builder

# Set working directory for backend
WORKDIR /app/backend

# Install dependencies
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy backend source code
COPY backend/ ./

# Build the Go application with CGO enabled for SQLite
RUN CGO_ENABLED=1 GOOS=linux go build -a -o viridian-bank-backend .

# Final stage
FROM debian:bookworm-slim

# Install ca-certificates and SQLite runtime
RUN apt-get update && apt-get install -y ca-certificates sqlite3 && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy the backend binary
COPY --from=backend-builder /app/backend/viridian-bank-backend ./

# Copy frontend files (HTML, CSS, JS, assets)
COPY index.html ./
COPY debug_test.html ./
COPY css/ ./css/
COPY js/ ./js/
COPY assets/ ./assets/

# Expose port
EXPOSE 8080

# Run the application
CMD ["./viridian-bank-backend"]

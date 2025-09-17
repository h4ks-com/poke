# Multi-stage build for Viridian City Bank
FROM golang:1.20-alpine AS backend-builder

# Set working directory for backend
WORKDIR /app/backend

# Install dependencies
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy backend source code
COPY backend/ ./

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o viridian-bank-backend .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests and serve static files
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the backend binary
COPY --from=backend-builder /app/backend/viridian-bank-backend ./

# Copy frontend files (HTML, CSS, JS, assets)
COPY index.html ./
COPY debug_test.html ./
COPY css/ ./css/
COPY js/ ./js/
COPY assets/ ./assets/

# Copy database and other necessary files
COPY backend/viridian_bank.db ./
COPY backend/.env* ./

# Expose port
EXPOSE 8080

# Run the application
CMD ["./viridian-bank-backend"]

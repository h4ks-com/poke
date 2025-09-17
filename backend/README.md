# Viridian City Bank - Backend

> **Note**: For the complete dockerized setup including frontend, see the main project README and DOCKER_README.md in the root directory.

This directory contains the Go backend for Viridian City Bank. The backend can be run standalone for development or as part of the unified Docker container.

## Development Setup (Standalone Backend)

### Prerequisites
- Go 1.19 or higher
- SQLite (database file included)

### Quick Start
```bash
cd backend
go mod download
go run .
```

The backend will start on port 8080 and serve API endpoints at `/api/*`.

## 🐳 Docker (Recommended)

For production use, use the unified Docker setup from the root directory:
```bash
cd .. # Go to project root
./run.sh
# OR
docker-compose up --build
```

This will run both frontend and backend together in a single container.

## Development vs Docker

### Development Mode (Standalone)
- Backend only on port 8080
- Frontend served separately (if needed)
- Good for backend development and testing

### Docker Mode (Unified)
- Backend + Frontend in one container
- Backend serves static files AND API
- Production-ready setup
- Access everything at http://localhost:8080

## API Endpoints

### Authentication
- `POST /api/register` - Register new user
- `POST /api/login` - User login
- `POST /api/change-password` - Change password (requires auth)

### Banking
- `GET /api/account` - Get account information
- `GET /api/balance` - Get account balance
- `POST /api/transfer` - Transfer money
- `GET /api/transactions` - Get transaction history
- `POST /api/payment-requests` - Create payment request
- `GET /api/payment-requests` - Get payment requests
- `PUT /api/payment-requests/:id` - Approve/reject payment request

### Health Check
- `GET /health` - Server health status

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | Required |
| `PORT` | Server port | 8080 |
| `JWT_SECRET` | JWT signing secret | Required |
| `WEBHOOK_URL` | Webhook endpoint for notifications | Optional |
| `ENVIRONMENT` | Environment (development/production) | development |

## Database Schema

The application uses the following tables:
- `users` - User accounts and authentication
- `transactions` - Transaction records
- `payment_requests` - Payment request records
- `user_sessions` - User session management

## Development

To run in development mode with hot reload:
```bash
# Install air for hot reload (optional)
go install github.com/cosmtrek/air@latest

# Run with hot reload
air
```

## Production Deployment

1. Build the application:
   ```bash
   go build -o viridian-bank .
   ```

2. Set production environment variables
3. Run the binary:
   ```bash
   ./viridian-bank
   ```

## Docker Support

Build and run with Docker:
```bash
# Build image
docker build -t viridian-bank .

# Run container
docker run -p 8080:8080 --env-file .env viridian-bank
```

## Webhook Notifications

The system can send webhook notifications for:
- User registration/login/password changes
- Money transfers
- Payment request creation/approval/rejection

Configure `WEBHOOK_URL` in your environment to receive these notifications.

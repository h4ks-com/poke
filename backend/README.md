# Viridian City Bank - Backend Setup

## Prerequisites

- Go 1.19 or higher
- PostgreSQL 12 or higher
- Git

## Installation

1. **Clone the repository and navigate to backend directory:**
   ```bash
   cd backend
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Set up environment variables:**
   ```bash
   cp .env.example .env
   # Edit .env file with your database credentials and configuration
   ```

4. **Set up PostgreSQL database:**
   ```sql
   -- Connect to PostgreSQL and create database
   CREATE DATABASE viridian_bank;
   
   -- Create a user (optional)
   CREATE USER bank_user WITH PASSWORD 'your_password';
   GRANT ALL PRIVILEGES ON DATABASE viridian_bank TO bank_user;
   ```

5. **Run database migrations:**
   ```bash
   go run scripts/migrate.go
   ```

6. **Start the server:**
   ```bash
   go run .
   ```

The server will start on the port specified in your `.env` file (default: 8080).

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

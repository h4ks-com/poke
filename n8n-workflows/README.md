# Banking Application - n8n Backend Workflows

This directory contains n8n workflow JSON files that implement the backend API for the virtual banking application.

## Overview

The backend consists of 7 main n8n workflows that provide RESTful API endpoints for:
- Database setup and initialization
- User registration with automatic account creation
- Authentication (login/logout) 
- Password management (change password)
- Account management (balance and transactions)
- Money transfers between users
- Payment request management

## Account Creation Process

**New users can register themselves** through the `/api/register` endpoint:
1. **Self-Registration**: Users provide username, password, and confirmation
2. **Automatic Setup**: System creates account with unique account number
3. **Starting Balance**: New accounts receive 1000 rubles to begin with
4. **Security**: Passwords are hashed with bcrypt (12 salt rounds)
5. **Validation**: Username and password strength requirements enforced

## Prerequisites

1. **n8n Installation**: Install n8n either via npm, Docker, or n8n Cloud
2. **PostgreSQL Database**: Set up a PostgreSQL database instance
3. **Database Setup**: Run the provided SQL script to create the required tables

## Setup Instructions

### 1. Database Setup

First, run the database setup script in your PostgreSQL instance:

```bash
psql -U your_username -d your_database -f 01-database-setup.sql
```

This creates:
- `users` table for user accounts
- `user_sessions` table for authentication sessions
- `transactions` table for transaction history
- `payment_requests` table for payment requests
- Demo user accounts for testing (password: "password123" for all demo users)

**Demo Accounts Available:**
- Guild_Master_Alex (25,000 ₽)
- Guild_Member_Sarah (8,750 ₽) 
- Trader_Mike (12,300 ₽)
- Party_Member_Luna (6,500 ₽)
- Equipment_Vendor (18,900 ₽)
- Tournament_Organizer (50,000 ₽)

### 2. n8n Configuration

1. Import each workflow JSON file into your n8n instance
2. Configure the PostgreSQL credentials:
   - Create a new PostgreSQL credential in n8n
   - Name it "Banking Database" 
   - Set the connection details to match your database

### 3. Workflow Import Order

Import the workflows in this order:

1. `02-authentication-workflow.json` - Authentication endpoints
2. `03-account-management-workflow.json` - Account and transaction endpoints  
3. `04-transfer-workflow.json` - Money transfer endpoints
4. `05-payment-request-workflow.json` - Payment request endpoints
5. `06-user-registration-workflow.json` - User registration endpoint
6. `07-password-change-workflow.json` - Password change endpoint

## API Endpoints

### Authentication Workflow

**POST /api/register**
- Creates a new user account
- Request body:
  ```json
  {
    "username": "newuser",
    "password": "password123",
    "confirmPassword": "password123"
  }
  ```
- Returns user details and starts with 1000 rubles balance
- Username must be 3-20 characters (letters, numbers, underscores only)
- Password must be 8+ characters with at least one letter and one number

**POST /api/login**
- Authenticates user credentials
- Returns session token for authenticated requests

**POST /api/logout**  
- Invalidates user session
- Requires Authorization header with Bearer token

**POST /api/change-password**
- Changes user password and invalidates all sessions
- Request body:
  ```json
  {
    "currentPassword": "oldpassword123",
    "newPassword": "newpassword123", 
    "confirmNewPassword": "newpassword123"
  }
  ```
- Requires Authorization header with Bearer token
- User must log in again after password change

### Account Management Workflow

**GET /api/account/balance**
- Returns user's current balance and account number
- Requires Authorization header with Bearer token

**GET /api/account/transactions**
- Returns user's transaction history
- Query parameter: `limit` (optional, defaults to 10)
- Requires Authorization header with Bearer token

### Transfer Workflow

**POST /api/transfer**
- Transfers money between users
- Request body:
  ```json
  {
    "recipient": "username",
    "amount": 100.50,
    "memo": "Payment description"
  }
  ```
- Requires Authorization header with Bearer token

### Payment Request Workflow

**POST /api/payment-request**
- Creates a new payment request
- Request body:
  ```json
  {
    "fromUser": "username",
    "amount": 50.00,
    "description": "Request description"
  }
  ```
- Requires Authorization header with Bearer token

**GET /api/payment-requests**
- Gets payment requests for the authenticated user
- Query parameter: `type` ("sent" or "received", defaults to "received")
- Requires Authorization header with Bearer token

**POST /api/payment-request/{id}/respond**
- Accepts or rejects a payment request
- Request body:
  ```json
  {
    "action": "accept" // or "reject"
  }
  ```
- Requires Authorization header with Bearer token

## Authentication

All endpoints except login require authentication via Bearer token:

```
Authorization: Bearer <session_token>
```

Session tokens are returned by the login endpoint and expire after 24 hours.

## Database Schema

### Users Table
- `id` - Primary key
- `username` - Unique username
- `password_hash` - Bcrypt hashed password
- `account_number` - Unique account identifier
- `balance` - Current account balance (DECIMAL)
- `created_at` - Account creation timestamp

### User Sessions Table
- `session_id` - Primary key
- `user_id` - Foreign key to users table
- `session_token` - Unique session identifier
- `expires_at` - Session expiration timestamp
- `created_at` - Session creation timestamp

### Transactions Table
- `transaction_id` - Primary key
- `from_user_id` - Sender user ID (nullable for deposits)
- `to_user_id` - Recipient user ID (nullable for withdrawals)
- `amount` - Transaction amount
- `transaction_type` - Type of transaction
- `description` - Transaction description
- `memo` - Optional memo/note
- `created_at` - Transaction timestamp

### Payment Requests Table
- `request_id` - Primary key
- `from_user_id` - Requesting user ID
- `to_user_id` - Target user ID
- `amount` - Requested amount
- `description` - Request description
- `status` - Request status (pending/accepted/rejected)
- `created_at` - Request creation timestamp
- `responded_at` - Response timestamp

## Frontend Integration

The frontend JavaScript application should be updated to point to your n8n webhook URLs instead of the mock API endpoints. Update the `API_BASE_URL` in `js/api.js` to match your n8n instance.

Example:
```javascript
const API_BASE_URL = 'https://your-n8n-instance.com';
```

## Security Considerations

1. **HTTPS Only**: Always use HTTPS in production
2. **Session Management**: Sessions expire automatically after 24 hours
3. **Password Security**: 
   - Passwords hashed with bcrypt (12 salt rounds)
   - Minimum 8 characters with at least one letter and one number
   - Password change invalidates all existing sessions
4. **Username Requirements**: 3-20 characters, alphanumeric and underscores only
5. **SQL Injection Protection**: All queries use parameterized statements
6. **Authorization**: All endpoints validate session tokens
7. **Account Security**: Starting balance of 1000 rubles for new accounts

## Testing

You can test the endpoints using curl or any HTTP client:

```bash
# Login
curl -X POST https://your-n8n-instance.com/api/login \
  -H "Content-Type: application/json" \
  -d '{"username": "alice", "password": "password123"}'

# Get balance (replace TOKEN with actual token)
curl -X GET https://your-n8n-instance.com/api/account/balance \
  -H "Authorization: Bearer TOKEN"
```

## Troubleshooting

1. **Database Connection Issues**: Verify PostgreSQL credentials in n8n
2. **Workflow Errors**: Check n8n execution logs for detailed error messages
3. **Authentication Failures**: Ensure session tokens are valid and not expired
4. **Transaction Failures**: Check account balances and user permissions

## Production Deployment

For production use:

1. Use environment variables for database credentials
2. Set up proper SSL/TLS certificates
3. Configure rate limiting and request validation
4. Set up monitoring and logging
5. Use a proper session store (Redis recommended)
6. Implement backup strategies for the database

## Support

If you encounter any issues with the workflows, check:
1. n8n execution logs
2. PostgreSQL logs
3. Network connectivity between n8n and PostgreSQL
4. Credential configuration in n8n

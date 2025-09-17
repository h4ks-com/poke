# Bank - Virtual Banking System

A realistic-looking but clearly fake banking application built with vanilla JavaScript for game purposes.

## Features

- **Secure Login System** (Demo credentials: username: `player1`, password: `password123`)
- **Account Balance Display** with real-time updates
- **Money Transfer** functionality
- **Deposit & Withdrawal** operations
- **Transaction History** with detailed records
- **Responsive Design** for desktop and mobile
- **Game-themed elements** to clearly indicate it's virtual

## Demo Account

- **Username:** `player1`
- **Password:** `password123`
- **Starting Balance:** $15,450.00

## File Structure

```
central_bank/
├── index.html          # Main HTML file
├── css/
│   └── styles.css      # All styling and responsive design
├── js/
│   ├── app.js          # Main application logic
│   └── api.js          # Mock API endpoints for future backend
└── README.md           # This file
```

## Getting Started

1. Open `index.html` in a web browser
2. Use the demo credentials to log in
3. Explore the banking features:
   - View account balance
   - Send money transfers
   - Make deposits and withdrawals
   - Check transaction history

## Game Integration Notes

- The app uses `localStorage` to persist data across sessions
- Clear visual indicators show this is a "Virtual Banking System"
- Game-themed transaction categories (Quest Rewards, Equipment Purchase, etc.)
- All transactions include gaming-related descriptions

## API Structure

The application includes the following API endpoints:

### Authentication Endpoints
- `POST /api/register` - User registration
- `POST /api/login` - User authentication
- `POST /api/change-password` - Change user password

### Banking Endpoints (Authenticated)
- `GET /api/account` - Get account information
- `GET /api/balance` - Get current balance
- `GET /api/transactions` - Get transaction history
- `POST /api/transfer` - Send money transfer
- `POST /api/payment-requests` - Create payment request
- `GET /api/payment-requests` - Get payment requests
- `PUT /api/payment-requests/:id` - Handle payment request
- `GET /api/card` - Get card information
- `POST /api/card/refresh` - Refresh card number

### Administrative Endpoints (Admin Key Required)
- `POST /api/admin/adjust-balance` - Adjust user balance
- `POST /api/admin/merchant-transaction` - Create merchant transaction
- `GET /api/admin/users` - Get all users
- `GET /api/admin/user/:account` - Get user by account number

## Administrative API Usage

### Authentication
All admin endpoints require the `X-Admin-Key` header:
```bash
curl -H "X-Admin-Key: your-admin-secret-key" ...
```

### Adjust User Balance
```bash
POST /api/admin/adjust-balance
Content-Type: application/json
X-Admin-Key: your-admin-secret-key

{
  "user_id": 123,
  "amount": 1000.00,
  "description": "Quest reward payment",
  "merchant_name": "Pokémon Center"
}
```

### Create Merchant Transaction
```bash
POST /api/admin/merchant-transaction
Content-Type: application/json
X-Admin-Key: your-admin-secret-key

{
  "user_id": 123,
  "amount": -50.00,
  "description": "Potion purchase",
  "merchant_name": "Mart"
}
```

### Get All Users
```bash
GET /api/admin/users
X-Admin-Key: your-admin-secret-key
```

## Security Features

- **Admin Key Authentication**: All administrative functions require a secret admin key
- **PokéBank Balance**: The PokéBank system account maintains a fixed balance of 999,999,999.99
- **Virtual Merchants**: Transactions can be created to/from non-existent merchant accounts
- **Transaction Logging**: All administrative actions are logged in the transaction history
- **Webhook Notifications**: Admin actions trigger webhook notifications for external systems

## Customization

To integrate with your actual backend:

1. Update the `baseURL` in `js/api.js`
2. Replace mock API calls with real HTTP requests
3. Update authentication logic
4. Modify transaction categories to match your game

## Features Overview

### Authentication
- Secure login form with validation
- Session management
- Logout functionality

### Account Management
- Real-time balance display
- Account number masking
- Transaction categorization

### Transactions
- Money transfers between accounts
- Deposits from various sources
- Withdrawals for different purposes
- Comprehensive transaction history

### User Experience
- Responsive design for all devices
- Loading states for all operations
- Success/error message notifications
- Professional banking UI with game elements

## Browser Support

- Modern browsers (Chrome, Firefox, Safari, Edge)
- Mobile browsers
- No external dependencies except Google Fonts and Font Awesome

## Security Note

This is a demonstration application only. For production use, implement proper:
- Server-side validation
- Authentication tokens
- HTTPS encryption
- Input sanitization
- Rate limiting

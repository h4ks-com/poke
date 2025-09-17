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

The application includes placeholder API endpoints in `js/api.js` that simulate backend calls:

- `POST /api/auth/login` - User authentication
- `POST /api/auth/logout` - User logout
- `GET /api/account/balance` - Get current balance
- `GET /api/account/transactions` - Get transaction history
- `POST /api/transactions/transfer` - Send money
- `POST /api/transactions/deposit` - Deposit funds
- `POST /api/transactions/withdraw` - Withdraw funds

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

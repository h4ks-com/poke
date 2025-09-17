-- Database Setup for Virtual Banking System
-- Run this SQL to create the required tables

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    account_number VARCHAR(20) UNIQUE NOT NULL,
    balance DECIMAL(15,2) DEFAULT 0.00,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Transactions table
CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    transaction_id VARCHAR(50) UNIQUE NOT NULL,
    from_user_id INTEGER REFERENCES users(id),
    to_user_id INTEGER REFERENCES users(id),
    amount DECIMAL(15,2) NOT NULL,
    transaction_type VARCHAR(20) NOT NULL, -- 'transfer', 'payment_request_fulfillment'
    description TEXT,
    memo TEXT,
    status VARCHAR(20) DEFAULT 'completed', -- 'completed', 'pending', 'failed'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Payment requests table
CREATE TABLE IF NOT EXISTS payment_requests (
    id SERIAL PRIMARY KEY,
    request_id VARCHAR(50) UNIQUE NOT NULL,
    from_user_id INTEGER REFERENCES users(id), -- who is requesting money
    to_user_id INTEGER REFERENCES users(id),   -- who should pay
    amount DECIMAL(15,2) NOT NULL,
    reason TEXT NOT NULL,
    message TEXT,
    status VARCHAR(20) DEFAULT 'pending', -- 'pending', 'accepted', 'rejected', 'cancelled'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- User sessions table (for authentication)
CREATE TABLE IF NOT EXISTS user_sessions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    session_token VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for better performance
CREATE INDEX IF NOT EXISTS idx_transactions_user_id ON transactions(from_user_id, to_user_id);
CREATE INDEX IF NOT EXISTS idx_payment_requests_users ON payment_requests(from_user_id, to_user_id);
CREATE INDEX IF NOT EXISTS idx_payment_requests_status ON payment_requests(status);
CREATE INDEX IF NOT EXISTS idx_user_sessions_token ON user_sessions(session_token);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);

-- Insert demo user
-- Insert first demo user (password: "password123")
INSERT INTO users (username, password_hash, account_number, balance) 
VALUES ('player1', '$2b$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewS9eapGdpqNqSge', 'ACC0123456789', 15450.00)
ON CONFLICT (username) DO NOTHING;

-- Insert demo users for testing (password: "password123" for all accounts)
-- Note: In production, these should be properly hashed with bcrypt
INSERT INTO users (username, password_hash, account_number, balance) VALUES
('Guild_Master_Alex', '$2b$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewS9eapGdpqNqSge', 'ACC4567890123', 25000.00),
('Guild_Member_Sarah', '$2b$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewS9eapGdpqNqSge', 'ACC7890123456', 8750.00),
('Trader_Mike', '$2b$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewS9eapGdpqNqSge', 'ACC1234567890', 12300.00),
('Party_Member_Luna', '$2b$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewS9eapGdpqNqSge', 'ACC9876543210', 6500.00),
('Equipment_Vendor', '$2b$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewS9eapGdpqNqSge', 'ACC2468135790', 18900.00),
('Tournament_Organizer', '$2b$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewS9eapGdpqNqSge', 'ACC1357924680', 50000.00)
ON CONFLICT (username) DO NOTHING;

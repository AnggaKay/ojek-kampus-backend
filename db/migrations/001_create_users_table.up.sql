-- Migration: 001_create_users_table
-- Description: Create unified users table for Passenger and Driver
-- Auth Method: Password-based authentication (Option B)

-- Step 1: Create ENUM types for role and status
CREATE TYPE user_role AS ENUM ('PASSENGER', 'DRIVER');
CREATE TYPE user_status AS ENUM ('ACTIVE', 'SUSPENDED', 'PENDING_VERIFICATION', 'REJECTED');

-- Step 2: Create users table
CREATE TABLE users (
    -- Primary Key: Auto-increment integer
    id SERIAL PRIMARY KEY,

    -- Authentication fields
    phone_number VARCHAR(20) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,  -- Required for all users (password-based auth)
    
    -- Profile fields
    email VARCHAR(255) UNIQUE,  -- Optional, but unique if provided
    full_name VARCHAR(100) NOT NULL,  -- User's full name
    
    -- Role & Status
    role user_role NOT NULL DEFAULT 'PASSENGER',
    status user_status NOT NULL DEFAULT 'ACTIVE',
    
    -- Verification & Security
    phone_verified BOOLEAN DEFAULT FALSE,  -- Track phone verification via OTP
    last_login_at TIMESTAMP,  -- Track last login time for security
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Step 3: Create indexes for performance
CREATE INDEX idx_users_phone ON users(phone_number);
CREATE INDEX idx_users_email ON users(email) WHERE email IS NOT NULL;  -- Partial index for non-null emails
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_role_status ON users(role, status);  -- Composite index for common queries

-- Step 4: Create function for auto-updating updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Step 5: Create trigger for users table
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Step 6: Add table and column comments for documentation
COMMENT ON TABLE users IS 'Unified user table for Passenger and Driver authentication';
COMMENT ON COLUMN users.phone_number IS 'Unique phone number in E.164 format (+6281234567890)';
COMMENT ON COLUMN users.password_hash IS 'Bcrypt hash of user password (cost factor 12)';
COMMENT ON COLUMN users.email IS 'Optional email address, unique if provided';
COMMENT ON COLUMN users.full_name IS 'User full name for display';
COMMENT ON COLUMN users.role IS 'User role: PASSENGER or DRIVER';
COMMENT ON COLUMN users.status IS 'Account status: ACTIVE (default), SUSPENDED, PENDING_VERIFICATION (drivers), REJECTED';
COMMENT ON COLUMN users.phone_verified IS 'Whether phone number has been verified via OTP';
COMMENT ON COLUMN users.last_login_at IS 'Timestamp of last successful login';

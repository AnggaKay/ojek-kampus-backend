-- Migration: 004_create_admin_users
-- Description: Create admin users table for dashboard authentication (separate from users table)

CREATE TYPE admin_role AS ENUM ('SUPER_ADMIN', 'ADMIN', 'VIEWER');

CREATE TABLE admin_users (
    id SERIAL PRIMARY KEY,
    
    -- Authentication
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    
    -- Profile
    full_name VARCHAR(100) NOT NULL,
    
    -- Role & Status
    role admin_role NOT NULL DEFAULT 'ADMIN',
    is_active BOOLEAN DEFAULT TRUE,
    
    -- Tracking
    last_login_at TIMESTAMP,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_admin_username ON admin_users(username);
CREATE INDEX idx_admin_email ON admin_users(email);
CREATE INDEX idx_admin_role ON admin_users(role);
CREATE INDEX idx_admin_active ON admin_users(is_active) WHERE is_active = TRUE;

-- Trigger for auto-update updated_at
CREATE TRIGGER update_admin_users_updated_at
    BEFORE UPDATE ON admin_users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Comments
COMMENT ON TABLE admin_users IS 'Admin users for dashboard access (separate from passenger/driver)';
COMMENT ON COLUMN admin_users.username IS 'Unique username for admin login';
COMMENT ON COLUMN admin_users.role IS 'Admin role: SUPER_ADMIN (full access), ADMIN (moderate), VIEWER (read-only)';
COMMENT ON COLUMN admin_users.is_active IS 'Whether admin account is active (for account suspension)';

-- Add foreign key constraint to driver_profiles (verified_by)
-- Note: This must be done AFTER admin_users table is created
ALTER TABLE driver_profiles
    ADD CONSTRAINT fk_driver_verified_by FOREIGN KEY (verified_by) 
        REFERENCES admin_users(id) ON DELETE SET NULL;

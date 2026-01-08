-- Migration: 006_create_refresh_tokens
-- Description: Create refresh tokens table for JWT token management and device tracking

CREATE TABLE refresh_tokens (
    id SERIAL PRIMARY KEY,
    
    -- User reference
    user_id INT NOT NULL,
    user_type VARCHAR(20) NOT NULL,  -- 'PASSENGER', 'DRIVER', 'ADMIN'
    
    -- Token data
    token_hash TEXT NOT NULL UNIQUE,  -- SHA-256 hash of refresh token
    
    -- Device tracking
    device_info TEXT,  -- User-Agent string
    device_name VARCHAR(100),  -- e.g., "iPhone 13", "Chrome on Windows"
    ip_address VARCHAR(45),  -- IPv4 or IPv6
    
    -- Expiration & revocation
    expires_at TIMESTAMP NOT NULL,
    is_revoked BOOLEAN DEFAULT FALSE,
    revoked_at TIMESTAMP,
    revoke_reason VARCHAR(50),  -- 'LOGOUT', 'PASSWORD_CHANGE', 'SUSPICIOUS_ACTIVITY', 'MANUAL'
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_used_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_refresh_user ON refresh_tokens(user_id, user_type);
CREATE INDEX idx_refresh_token ON refresh_tokens(token_hash);
CREATE INDEX idx_refresh_expires ON refresh_tokens(expires_at);
CREATE INDEX idx_refresh_active ON refresh_tokens(user_id, is_revoked) WHERE is_revoked = FALSE;

-- Comments
COMMENT ON TABLE refresh_tokens IS 'Refresh tokens for JWT-based authentication and multi-device support';
COMMENT ON COLUMN refresh_tokens.user_id IS 'Reference to user ID (users.id or admin_users.id)';
COMMENT ON COLUMN refresh_tokens.user_type IS 'User type: PASSENGER, DRIVER, or ADMIN';
COMMENT ON COLUMN refresh_tokens.token_hash IS 'SHA-256 hash of refresh token for security';
COMMENT ON COLUMN refresh_tokens.device_info IS 'User-Agent string for device identification';
COMMENT ON COLUMN refresh_tokens.is_revoked IS 'Whether token has been revoked (logout, password change, etc.)';
COMMENT ON COLUMN refresh_tokens.expires_at IS 'Token expiration time (typically 7 days)';

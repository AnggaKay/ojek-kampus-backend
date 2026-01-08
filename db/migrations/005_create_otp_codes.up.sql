-- Migration: 005_create_otp_codes
-- Description: Create OTP codes table for phone verification and password reset

CREATE TYPE otp_purpose AS ENUM ('REGISTRATION', 'PASSWORD_RESET', 'PHONE_VERIFICATION');

CREATE TABLE otp_codes (
    id SERIAL PRIMARY KEY,
    
    -- OTP data
    phone_number VARCHAR(20) NOT NULL,
    otp_code VARCHAR(6) NOT NULL,
    purpose otp_purpose NOT NULL,
    
    -- Expiration & usage
    expires_at TIMESTAMP NOT NULL,
    is_used BOOLEAN DEFAULT FALSE,
    used_at TIMESTAMP,
    
    -- Brute force protection
    attempts INT DEFAULT 0,
    
    -- Metadata
    ip_address VARCHAR(45),  -- IPv4 or IPv6
    user_agent TEXT,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_otp_phone ON otp_codes(phone_number);
CREATE INDEX idx_otp_phone_purpose ON otp_codes(phone_number, purpose);
CREATE INDEX idx_otp_expires ON otp_codes(expires_at);
CREATE INDEX idx_otp_used ON otp_codes(is_used, expires_at);

-- Comments
COMMENT ON TABLE otp_codes IS 'OTP codes for phone verification and password reset';
COMMENT ON COLUMN otp_codes.phone_number IS 'Phone number to send OTP (E.164 format)';
COMMENT ON COLUMN otp_codes.otp_code IS '6-digit numeric OTP code';
COMMENT ON COLUMN otp_codes.purpose IS 'Purpose: REGISTRATION, PASSWORD_RESET, or PHONE_VERIFICATION';
COMMENT ON COLUMN otp_codes.expires_at IS 'OTP expiration time (typically 5 minutes from creation)';
COMMENT ON COLUMN otp_codes.is_used IS 'Whether OTP has been used (one-time use only)';
COMMENT ON COLUMN otp_codes.attempts IS 'Number of verification attempts (max 3)';

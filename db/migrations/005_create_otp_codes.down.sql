-- Rollback migration: 005_create_otp_codes

DROP INDEX IF EXISTS idx_otp_used;
DROP INDEX IF EXISTS idx_otp_expires;
DROP INDEX IF EXISTS idx_otp_phone_purpose;
DROP INDEX IF EXISTS idx_otp_phone;
DROP TABLE IF EXISTS otp_codes;
DROP TYPE IF EXISTS otp_purpose;

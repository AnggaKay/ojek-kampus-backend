-- Rollback migration: 004_create_admin_users

-- Remove foreign key from driver_profiles first
ALTER TABLE driver_profiles DROP CONSTRAINT IF EXISTS fk_driver_verified_by;

-- Drop trigger and indexes
DROP TRIGGER IF EXISTS update_admin_users_updated_at ON admin_users;
DROP INDEX IF EXISTS idx_admin_active;
DROP INDEX IF EXISTS idx_admin_role;
DROP INDEX IF EXISTS idx_admin_email;
DROP INDEX IF EXISTS idx_admin_username;

-- Drop table and ENUM
DROP TABLE IF EXISTS admin_users;
DROP TYPE IF EXISTS admin_role;

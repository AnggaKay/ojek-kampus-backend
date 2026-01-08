-- Rollback migration: 001_create_users_table
-- This will undo all changes made in the UP migration

-- Step 1: Drop trigger
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Step 2: Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Step 3: Drop indexes
DROP INDEX IF EXISTS idx_users_role_status;
DROP INDEX IF EXISTS idx_users_status;
DROP INDEX IF EXISTS idx_users_role;
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_phone;

-- Step 4: Drop table
DROP TABLE IF EXISTS users;

-- Step 5: Drop ENUM types
DROP TYPE IF EXISTS user_status;
DROP TYPE IF EXISTS user_role;

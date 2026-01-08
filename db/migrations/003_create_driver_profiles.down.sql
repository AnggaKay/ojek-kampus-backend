-- Rollback migration: 003_create_driver_profiles

DROP TRIGGER IF EXISTS update_driver_profiles_updated_at ON driver_profiles;
DROP INDEX IF EXISTS idx_driver_rating;
DROP INDEX IF EXISTS idx_driver_vehicle_plate;
DROP INDEX IF EXISTS idx_driver_location;
DROP INDEX IF EXISTS idx_driver_active;
DROP INDEX IF EXISTS idx_driver_verified;
DROP INDEX IF EXISTS idx_driver_user;
DROP TABLE IF EXISTS driver_profiles;

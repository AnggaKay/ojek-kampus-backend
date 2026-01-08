-- Rollback migration: 002_create_passenger_profiles

DROP TRIGGER IF EXISTS update_passenger_profiles_updated_at ON passenger_profiles;
DROP INDEX IF EXISTS idx_passenger_rating;
DROP INDEX IF EXISTS idx_passenger_user;
DROP TABLE IF EXISTS passenger_profiles;

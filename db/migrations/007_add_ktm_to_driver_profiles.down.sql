-- Rollback migration: 007_add_ktm_to_driver_profiles

DROP INDEX IF EXISTS idx_driver_ktm;

ALTER TABLE driver_profiles
    DROP COLUMN IF EXISTS ktm_photo,
    DROP COLUMN IF EXISTS ktm_number;

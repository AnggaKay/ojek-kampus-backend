-- Migration: 007_add_ktm_to_driver_profiles
-- Description: Add KTM (Kartu Tanda Mahasiswa) fields for student driver verification

ALTER TABLE driver_profiles
    ADD COLUMN ktm_number VARCHAR(20) UNIQUE,
    ADD COLUMN ktm_photo TEXT;

-- Index for KTM lookup
CREATE INDEX idx_driver_ktm ON driver_profiles(ktm_number) WHERE ktm_number IS NOT NULL;

-- Comments
COMMENT ON COLUMN driver_profiles.ktm_number IS 'Nomor KTM (Kartu Tanda Mahasiswa) - must be unique';
COMMENT ON COLUMN driver_profiles.ktm_photo IS 'URL to KTM photo for verification';

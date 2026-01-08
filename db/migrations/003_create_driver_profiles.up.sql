-- Migration: 003_create_driver_profiles
-- Description: Create driver profile table with vehicle info, documents, and verification status

CREATE TABLE driver_profiles (
    id SERIAL PRIMARY KEY,
    
    -- Foreign key to users table
    user_id INT NOT NULL UNIQUE,
    
    -- Profile data
    profile_picture TEXT,
    fcm_token TEXT,
    
    -- Vehicle information
    vehicle_type VARCHAR(20) NOT NULL DEFAULT 'MOTOR',
    vehicle_plate VARCHAR(15) NOT NULL UNIQUE,
    vehicle_brand VARCHAR(50),
    vehicle_model VARCHAR(50),
    vehicle_color VARCHAR(30),
    
    -- Document verification
    ktp_photo TEXT,  -- URL to KTP (ID card) photo
    sim_photo TEXT,  -- URL to SIM (driver license) photo
    stnk_photo TEXT, -- URL to STNK (vehicle registration) photo
    
    -- Verification status
    is_verified BOOLEAN DEFAULT FALSE,
    verification_notes TEXT,
    verified_by INT,  -- Reference to admin_users.id
    verified_at TIMESTAMP,
    rejection_reason TEXT,
    
    -- Driver status
    is_active BOOLEAN DEFAULT FALSE,  -- Online/offline toggle by driver
    
    -- Location tracking
    current_lat DECIMAL(10, 8),
    current_long DECIMAL(11, 8),
    last_location_update TIMESTAMP,
    
    -- Metrics
    total_completed_orders INT DEFAULT 0,
    total_cancelled_orders INT DEFAULT 0,
    rating_avg DECIMAL(3, 2) DEFAULT 0.00 CHECK (rating_avg >= 0 AND rating_avg <= 5),
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Foreign key constraints
    CONSTRAINT fk_driver_user FOREIGN KEY (user_id) 
        REFERENCES users(id) ON DELETE CASCADE
);

-- Indexes
CREATE INDEX idx_driver_user ON driver_profiles(user_id);
CREATE INDEX idx_driver_verified ON driver_profiles(is_verified);
CREATE INDEX idx_driver_active ON driver_profiles(is_active) WHERE is_active = TRUE;
CREATE INDEX idx_driver_location ON driver_profiles(current_lat, current_long) WHERE is_active = TRUE;
CREATE INDEX idx_driver_vehicle_plate ON driver_profiles(vehicle_plate);
CREATE INDEX idx_driver_rating ON driver_profiles(rating_avg DESC);

-- Trigger for auto-update updated_at
CREATE TRIGGER update_driver_profiles_updated_at
    BEFORE UPDATE ON driver_profiles
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Comments
COMMENT ON TABLE driver_profiles IS 'Driver-specific profile data, vehicle info, and verification status';
COMMENT ON COLUMN driver_profiles.user_id IS 'Reference to users.id (must be role=DRIVER)';
COMMENT ON COLUMN driver_profiles.is_verified IS 'Whether driver has been verified by admin';
COMMENT ON COLUMN driver_profiles.verified_by IS 'Admin user ID who verified this driver';
COMMENT ON COLUMN driver_profiles.is_active IS 'Driver online/offline status (toggle by driver)';
COMMENT ON COLUMN driver_profiles.current_lat IS 'Current latitude for geo-location queries';
COMMENT ON COLUMN driver_profiles.current_long IS 'Current longitude for geo-location queries';

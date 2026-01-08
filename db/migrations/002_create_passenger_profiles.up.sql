-- Migration: 002_create_passenger_profiles
-- Description: Create passenger profile table with metrics and FCM token

CREATE TABLE passenger_profiles (
    id SERIAL PRIMARY KEY,
    
    -- Foreign key to users table
    user_id INT NOT NULL UNIQUE,
    
    -- Profile data
    profile_picture TEXT,  -- URL to profile photo storage
    fcm_token TEXT,  -- Firebase Cloud Messaging token for push notifications
    
    -- Metrics
    total_orders INT DEFAULT 0,
    total_cancellations INT DEFAULT 0,
    rating_avg DECIMAL(3, 2) DEFAULT 0.00 CHECK (rating_avg >= 0 AND rating_avg <= 5),
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Foreign key constraint
    CONSTRAINT fk_passenger_user FOREIGN KEY (user_id) 
        REFERENCES users(id) ON DELETE CASCADE
);

-- Indexes
CREATE INDEX idx_passenger_user ON passenger_profiles(user_id);
CREATE INDEX idx_passenger_rating ON passenger_profiles(rating_avg DESC);

-- Trigger for auto-update updated_at
CREATE TRIGGER update_passenger_profiles_updated_at
    BEFORE UPDATE ON passenger_profiles
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Comments
COMMENT ON TABLE passenger_profiles IS 'Passenger-specific profile data and metrics';
COMMENT ON COLUMN passenger_profiles.user_id IS 'Reference to users.id (must be role=PASSENGER)';
COMMENT ON COLUMN passenger_profiles.fcm_token IS 'Firebase token for push notifications to passenger mobile app';
COMMENT ON COLUMN passenger_profiles.total_orders IS 'Total number of orders placed by passenger';
COMMENT ON COLUMN passenger_profiles.total_cancellations IS 'Total number of orders cancelled by passenger';
COMMENT ON COLUMN passenger_profiles.rating_avg IS 'Average rating given by drivers (0.00 - 5.00)';

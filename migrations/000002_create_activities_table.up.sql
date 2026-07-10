CREATE TABLE IF NOT EXISTS activities (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(150) NOT NULL,
    sport_type VARCHAR(50) NOT NULL DEFAULT 'running',
    distance_km NUMERIC(6,2) NOT NULL,
    duration_seconds INT NOT NULL,
    avg_pace_seconds INT,
    elevation_gain_m INT DEFAULT 0,
    activity_date DATE NOT NULL,
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_activities_user_id ON activities(user_id);
CREATE INDEX IF NOT EXISTS idx_activities_activity_date ON activities(activity_date);
-- Create tables manually (if not running Go server first)
-- Run with: psql -h 51.250.86.178 -U nike -d audioguid -f create_tables.sql

-- Enable PostGIS extension
CREATE EXTENSION IF NOT EXISTS postgis;

-- POI table
CREATE TABLE IF NOT EXISTS places_of_interest (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    epoch VARCHAR(50),
    category VARCHAR(50),
    importance INTEGER DEFAULT 5,
    year_built INTEGER,
    architect VARCHAR(255),
    style VARCHAR(255),
    photos TEXT[],
    wikipedia_url VARCHAR(500),
    metadata JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_places_epoch ON places_of_interest(epoch);
CREATE INDEX IF NOT EXISTS idx_places_category ON places_of_interest(category);

-- Routes table
CREATE TABLE IF NOT EXISTS routes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255),
    description TEXT,
    total_distance DOUBLE PRECISION,
    estimated_duration INTEGER,
    epochs TEXT[],
    categories TEXT[],
    created_at TIMESTAMP DEFAULT NOW()
);

-- Waypoints table
CREATE TABLE IF NOT EXISTS waypoints (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    route_id UUID NOT NULL REFERENCES routes(id) ON DELETE CASCADE,
    poi_id UUID NOT NULL REFERENCES places_of_interest(id) ON DELETE CASCADE,
    "order" INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_waypoints_route_id ON waypoints(route_id);

-- Contents table
CREATE TABLE IF NOT EXISTS contents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    waypoint_id UUID NOT NULL UNIQUE REFERENCES waypoints(id) ON DELETE CASCADE,
    text TEXT NOT NULL,
    audio_url VARCHAR(500),
    audio_path VARCHAR(500),
    duration INTEGER,
    photos TEXT[],
    generated BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_contents_waypoint_id ON contents(waypoint_id);

-- Verify tables created
SELECT table_name FROM information_schema.tables 
WHERE table_schema = 'public' 
ORDER BY table_name;

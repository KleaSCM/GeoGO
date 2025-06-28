-- Unified Datasets Schema for GeoGO
-- This schema supports multiple dataset types in a single table

-- Drop existing tables if they exist
DROP TABLE IF EXISTS datasets CASCADE;
DROP TABLE IF EXISTS locations CASCADE;

-- Create the unified datasets table
CREATE TABLE datasets (
    id SERIAL PRIMARY KEY,
    dataset_type VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    lat DOUBLE PRECISION NOT NULL,
    lon DOUBLE PRECISION NOT NULL,
    value DOUBLE PRECISION,
    unit VARCHAR(20),
    timestamp TIMESTAMP,
    metadata JSONB,
    
    -- Legacy fields for backward compatibility
    recclass VARCHAR(100),
    mass DOUBLE PRECISION,
    year INTEGER,
    nametype VARCHAR(50),
    fall VARCHAR(50),
    
    -- Create spatial index
    geom GEOMETRY(POINT, 4326)
);

-- Create spatial index for fast geospatial queries
CREATE INDEX idx_datasets_geom ON datasets USING GIST (geom);
CREATE INDEX idx_datasets_type ON datasets (dataset_type);
CREATE INDEX idx_datasets_name ON datasets (name);

-- Create a function to automatically update the geometry column
CREATE OR REPLACE FUNCTION update_dataset_geometry()
RETURNS TRIGGER AS $$
BEGIN
    NEW.geom = ST_SetSRID(ST_MakePoint(NEW.lon, NEW.lat), 4326);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger to automatically update geometry
CREATE TRIGGER trigger_update_dataset_geometry
    BEFORE INSERT OR UPDATE ON datasets
    FOR EACH ROW
    EXECUTE FUNCTION update_dataset_geometry();

-- Create view for meteorites (backward compatibility)
CREATE VIEW meteorites AS
SELECT 
    id,
    name,
    recclass,
    mass,
    year,
    lat,
    lon,
    nametype,
    fall,
    geom
FROM datasets 
WHERE dataset_type = 'meteorite';

-- Create view for climate data
CREATE VIEW climate_stations AS
SELECT 
    id,
    name,
    lat,
    lon,
    value,
    unit,
    metadata,
    geom
FROM datasets 
WHERE dataset_type = 'climate';

-- Create view for wind data
CREATE VIEW wind_observations AS
SELECT 
    id,
    name,
    lat,
    lon,
    value,
    unit,
    timestamp,
    metadata,
    geom
FROM datasets 
WHERE dataset_type = 'wind';

-- Create view for vegetation data
CREATE VIEW vegetation_zones AS
SELECT 
    id,
    name,
    lat,
    lon,
    value,
    unit,
    metadata,
    geom
FROM datasets 
WHERE dataset_type = 'vegetation';

-- Insert sample data for testing
INSERT INTO datasets (dataset_type, name, lat, lon, value, unit, metadata) VALUES
('meteorite', 'Sample Meteorite', -37.8136, 144.9631, 1000.0, 'g', '{"recclass": "L6", "year": 2020, "fall": "Found"}'),
('climate', 'Melbourne Climate Station', -37.8136, 144.9631, 20.5, '°C', '{"station_name": "Melbourne", "model": "CSIRO"}'),
('wind', 'Melbourne Wind Station', -37.8136, 144.9631, 15.2, 'm/s', '{"location_description": "Melbourne CBD"}');

-- Grant permissions
GRANT ALL PRIVILEGES ON TABLE datasets TO postgres;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO postgres;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO postgres; 
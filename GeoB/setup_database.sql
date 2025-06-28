-- GeoGO Multi-Dataset Database Setup
-- Run this script to set up the unified database schema

-- Create the unified schema
\i utils/SQL/create_unified_schema.sql

-- Load processed data (run these after processing datasets)
-- \i utils/SQL/meteorites.sql
-- \i utils/SQL/climate_avg_temperature.sql
-- \i utils/SQL/climate_max_temperature.sql
-- \i utils/SQL/climate_min_temperature.sql
-- \i utils/SQL/climate_humidity.sql
-- \i utils/SQL/climate_evaporation.sql
-- \i utils/SQL/wind_observations.sql
-- \i utils/SQL/vegetation_zones.sql

-- Verify the setup
SELECT 
    dataset_type,
    COUNT(*) as record_count,
    MIN(value) as min_value,
    MAX(value) as max_value
FROM datasets 
GROUP BY dataset_type
ORDER BY dataset_type; 
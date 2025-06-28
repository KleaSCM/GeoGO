# 🌍 GeoGO Multi-Dataset Extension Setup Guide

## Overview
GeoGO has been extended to support multiple geospatial datasets beyond just meteorites! The app now includes:

- **☄️ Meteorites** - Global meteorite impact data (45,716 records)
- **🌡️ Climate Data** - Australian climate stations (43,720 records)
  - Average temperature, max/min temperature, humidity, evaporation
- **💨 Wind Observations** - Wind speed and direction data (60,416 records)
- **🌿 Vegetation Zones** - Vegetation classification data (7 records)
- **🏗️ Infrastructure** - Infrastructure and utility data
- **🔥 Fire Data** - Fire projection and risk data

## 🚀 Quick Setup

### 1. Database Setup
```bash
# Connect to your PostgreSQL database
psql -U postgres -d GeoGo

# Run the setup script
\i setup_database.sql
```

### 2. Process Your Datasets
```bash
cd GeoB
python utils/process_datasets.py
```

### 3. Load Data into Database
```bash
# In psql, run the generated SQL files:
\i utils/SQL/meteorites.sql
\i utils/SQL/climate_avg_temperature.sql
\i utils/SQL/climate_max_temperature.sql
\i utils/SQL/climate_min_temperature.sql
\i utils/SQL/climate_humidity.sql
\i utils/SQL/climate_evaporation.sql
\i utils/SQL/wind_observations.sql
\i utils/SQL/vegetation_zones.sql
```

### 4. Start the Backend
```bash
cd GeoB
go run main.go
```

### 5. Start the Frontend
```bash
cd geofe
npm run dev
```

## 🗄️ Database Schema

### Unified Datasets Table
```sql
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
    
    -- Spatial index
    geom GEOMETRY(POINT, 4326)
);
```

## 📡 New API Endpoints

### Unified Dataset Endpoints
- `GET /datasets` - Get all datasets with filtering
- `GET /datasets/types` - Get available dataset types and metadata
- `GET /datasets/stats/:type` - Get statistics for a dataset type
- `GET /datasets/:type` - Get datasets of a specific type

### Query Parameters
- `type` - Dataset type (meteorite, climate, wind, vegetation, etc.)
- `value_min` / `value_max` - Filter by value range
- `location` - Search near location (forward geocoding)
- `radius` - Search radius in meters
- `limit` / `offset` - Pagination

### Legacy Endpoints (Still Available)
- `GET /meteorites` - Original meteorite endpoints
- `GET /meteorites/largest`
- `GET /meteorites/nearby`
- `GET /meteorites/location`

## 🎨 Frontend Features

### Dataset Selector
- Visual cards for each dataset type
- Record counts and value ranges
- Easy switching between datasets

### Enhanced Search
- Value-based filtering (instead of just mass)
- Dataset-specific parameter mapping
- Unified search interface

### Dynamic Display
- Dataset-specific icons and titles
- Appropriate value display (mass, temperature, wind speed, etc.)
- Metadata parsing and display

## 📊 Dataset Statistics

| Dataset Type | Records | Description |
|--------------|---------|-------------|
| Meteorites | 45,716 | Global meteorite impacts |
| Climate (Avg Temp) | 11,648 | Average temperature stations |
| Climate (Max Temp) | 11,200 | Maximum temperature stations |
| Climate (Min Temp) | 11,200 | Minimum temperature stations |
| Climate (Humidity) | 4,160 | Humidity stations |
| Climate (Evaporation) | 5,512 | Evaporation stations |
| Wind | 60,416 | Wind observations |
| Vegetation | 7 | Vegetation zones |

## 🔧 Adding New Datasets

### 1. Add Dataset Type
```go
// In models/dataset.go
const (
    DatasetTypeYourNewType DatasetType = "your_new_type"
)
```

### 2. Add Processing Logic
```python
# In utils/process_datasets.py
def process_your_new_data(self):
    # Add your data processing logic
    pass
```

### 3. Update Frontend
```typescript
// In components/DatasetSelector.tsx
function getDatasetIcon(datasetType: string): string {
    case "your_new_type":
        return "🎯";
}
```

## 🚨 Troubleshooting

### Common Issues

1. **Database Connection**
   - Check PostgreSQL is running
   - Verify connection string in `db/database.go`
   - Ensure PostGIS extension is installed

2. **Data Processing**
   - Check file paths in `utils/process_datasets.py`
   - Ensure pandas is installed: `pip install pandas`
   - Verify CSV/Excel files are readable

3. **Frontend Issues**
   - Check API URL in environment variables
   - Ensure all dependencies are installed: `npm install`
   - Check browser console for errors

### Verification Commands
```bash
# Check database connection
psql -U postgres -d GeoGo -c "SELECT COUNT(*) FROM datasets;"

# Check dataset types
curl http://localhost:8080/datasets/types

# Test specific dataset
curl "http://localhost:8080/datasets?type=climate&limit=5"
```

## 🎯 Next Steps

1. **Load your data** into the unified schema
2. **Test the API** with different dataset types
3. **Explore the frontend** with dataset switching
4. **Add custom visualizations** for specific dataset types
5. **Implement advanced filtering** for dataset-specific fields

## 🌟 Features to Add

- **Advanced Visualizations** - Charts, heatmaps, time series
- **Dataset Comparison** - Side-by-side analysis
- **Export Functionality** - CSV, GeoJSON, KML export
- **Real-time Updates** - Live data feeds
- **User Preferences** - Saved searches and favorites

---

**🎉 Congratulations!** Your GeoGO app now supports multiple geospatial datasets with a unified interface! 
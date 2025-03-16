# GeoGO: Meteorite Impact Data API 🌍☄️

## Overview
GeoGO is a geospatial system that allows users to query meteorite impact locations worldwide. It provides a RESTful interface for retrieving meteorite data, searching for nearby impacts, and analyzing their distribution based on real-world scientific datasets. It also includes a **Next.js** frontend that displays these meteorites on an interactive map, showcasing the power of Go + geospatial technologies.

**Key Features:**
- Geospatial Meteorite Search – Query meteorite impact data by year, mass, location, or proximity
- Interactive Next.js Frontend – Displays results on searchable meteorite cards + an interactive map
- Real-World Dataset – Over 47,000+ recorded meteorite landings
- Find Nearby Meteorites – Uses PostGIS to compute distance-based queries
- Reverse Geocoding – Convert lat/lon → location names for better readability
- Optimized Geospatial Performance – Powered by PostgreSQL + PostGIS for fast spatial indexing
- High-Performance Go Backend – Concurrency (goroutines + channels) for parallel queries
- Caching with Redis – Speeds up location-based lookups and geocoding responses
- Automated Data Ingestion (ETL) – Optional Python-based automation for bulk data updates
- Leaflet.js Map – Visualize meteorites dynamically with real-world coordinates

---

## 🛠️ Tech Stack

| **Technology**             | **Purpose**                                   |
|----------------------------|-----------------------------------------------|
| **Go (Gin Framework)**     | High-performance backend & concurrency        |
| **PostgreSQL + PostGIS**   | Geospatial database & indexing                |
| **Redis**                  | Caching for reverse geocoding & queries       |
| **Python (Pandas)**        | Data processing & ingestion                   |
| **Leaflet + Next.js**      | Interactive map & UI for meteorite data       |
| **Docker**                 | Containerization for production deployments   |

---

## Installation & Setup

1. **Clone the repository**:
   ```bash
   git clone https://github.com/KleaSCM/GeoGO.git
   cd GeoGO
``
###  Install Dependencies:
Go dependencies
``go mod tidy``

Make sure you have PostgreSQL with PostGIS installed.
Update DSN (host/user/password/dbname) in db.InitDB() (e.g., dsn := "host=localhost port=5432 user=postgres password=secret dbname=GeoGO sslmode=disable")


Python dependencies (optional, for data processing) 
``pip install pandas``
### Prepare the Database

Optional: Run Python ETL scripts if needed (e.g., python utils/ParseData.py).
``python utils/ParseData.py``

### Seed the Database
``go run cmd/seed/main.go``
### Run the API Server
``go run main.go``

The server should now be running on http://localhost:8080.

### 📡API Endpoints📡

| **Endpoint**               | **Description**                                       |
|----------------------------|-------------------------------------------------------|
| `GET /meteorites`          | Fetch all meteorites with filters (year, mass, etc.)  |
| `GET /meteorites/largest`  | Fetch the top 10 largest meteorites                   |
| `GET /meteorites/nearby`   | Find meteorites near a given lat/lon + radius         |
| `GET /meteorites/location` | Reverse geocode a given lat/lon                       |


id – Unique meteorite ID
name – Official meteorite name
nametype – Valid/Relict classification
recclass – Type of meteorite
mass – Mass (grams)
fall – Whether it was found or fell
year – Year of impact/discovery
latitude & longitude – Impact coordinates

PostGIS functions like ST_X(geom) AS lon and ST_Y(geom) AS lat are used to retrieve numeric coordinates for the frontend.

### 🛠️ Future Improvements

Enhanced Caching for common queries (Redis expansions)
Further concurrency optimization & worker pools in Go
Frontend: Additional map layers, advanced filters, and UI enhancements


### 💡 Contributing

Want to improve GeoGO?
Fork the repo, make changes, and submit a pull request! :)

## 🌍 Explore the Universe, One Rock at a Time! ☄️

# GeoGO: Meteorite Impact Data API 🌍☄️

## Overview
GeoGO is a geospatial system that allows users to query meteorite impact locations worldwide. It provides a RESTful interface for retrieving meteorite data, searching for nearby impacts, and analyzing their distribution based on real-world scientific datasets. It also includes a **Next.js** frontend that displays these meteorites on an interactive map, showcasing the power of Go + geospatial technologies.

**Key Features:**
- RESTful API for meteorite impact data  
- Find Nearby Meteorites using precise geospatial queries  
- Real-World Dataset (47,000+ meteorite landings)  
- Optimized with **PostgreSQL + PostGIS** for geospatial indexing  
- **Concurrency** in Go (goroutines + channels) for parallel queries  
- **Redis caching** for reverse geocoding and performance boosts  
- Automated Data Ingestion with Python (optional ETL)  
- A **Next.js** frontend that displays meteorite cards and an interactive Leaflet map

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


### 🛠️ Future Improvements

🔹 Caching for performance optimization
🔹 Spatial Indexing with PostGIS or RTree
🔹 Data Visualization (Mapping impacts on a frontend)
🔹 Machine Learning to predict impact zones


### 💡 Contributing

Want to improve GeoGO?
Fork the repo, make changes, and submit a pull request! :)

## 🌍 Explore the Universe, One Rock at a Time! ☄️

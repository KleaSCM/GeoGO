# GeoGO: Meteorite Impact Data API 🌍☄️

## Overview

GeoGO is a geospatial API that allows users to query meteorite impact locations worldwide. It provides a RESTful interface for retrieving meteorite landings, searching for nearby impacts, and analyzing their distribution based on real-world scientific data.

This project integrates
    - Public Meteorite Data (NASA/Global Databases)
    - SQLite for fast lookups
    - Go (Gin framework) for API services
    - Haversine Formula for distance-based queries
    - Python for ETL (Extract, Transform, Load)
    
Features

- RESTful API for meteorite impact data
- Find Nearby Meteorites using precise geospatial calculations
- Real-World Dataset (over 47,000+ meteorite landings)
- Optimized with SQLite for fast queries
- Automated Data Ingestion with Python
- Custom Database Schema for structured analysis

## 🛠️ Tech Stack

| **Technology**          | **Purpose**                               |
|------------------------|-------------------------------------------|
| **Go (Gin Framework)** | High-performance API backend             |
| **SQLite**            | Lightweight relational database           |
| **Python (Pandas)**   | Data processing & ingestion               |
| **Haversine Formula** | Geospatial distance calculations         |
| **cURL/Postman**      | API Testing                               |


## Setup & Installation
``git clone https://github.com/yourusername/GeoGO.git
cd GeoGO``
###  Install Dependencies:
Go dependencies
``go mod tidy``
Python dependencies (optional, for data processing)
``pip install pandas``
### Prepare the Database
``python utils/ParseData.py``
### Seed the Database
``go run cmd/seed/main.go``
### Run the API Server
``go run main.go``

### 📡API Endpoints📡

## 🌍 Get All Meteorite Landings
``GET /locations``

###  Database Schema
``CREATE TABLE locations (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  nametype TEXT,
  recclass TEXT,
  mass REAL,
  fall TEXT,
  year INTEGER,
  latitude REAL NOT NULL,
  longitude REAL NOT NULL
);``

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

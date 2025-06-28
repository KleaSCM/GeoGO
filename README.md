# 🌍 GeoGO: Global Geospatial Data Explorer

> **Explore the Universe, One Dataset at a Time** ☄️🌡️💨🌿🏗️🔥

## Overview

GeoGO is a modern, full-stack geospatial data exploration platform that allows users to query and visualize multiple types of global geospatial datasets. From meteorite impacts to climate data, wind observations to vegetation zones, GeoGO provides a unified interface for exploring Earth's diverse geospatial information.

**Key Features:**
- 🌍 **Multi-Dataset Support** – Meteorites, Climate, Wind, Vegetation, Infrastructure, Fire data
- 🎨 **Beautiful Dark Mode UI** – Modern, responsive Next.js frontend with stunning visual design
- 🗺️ **Interactive Maps** – Leaflet.js integration for dynamic geospatial visualization
- 🔍 **Advanced Search** – Dataset-specific search forms with location-based queries
- ⚡ **High Performance** – Go backend with PostgreSQL/PostGIS for lightning-fast spatial queries
- 🎯 **Real-World Data** – 100,000+ records across multiple scientific datasets
- 🔄 **Live Statistics** – Real-time dataset counts and metadata
- 📱 **Responsive Design** – Works perfectly on desktop, tablet, and mobile

---

## 🛠️ Tech Stack

| **Technology**             | **Purpose**                                   |
|----------------------------|-----------------------------------------------|
| **Go (Gin Framework)**     | High-performance backend API                  |
| **PostgreSQL + PostGIS**   | Geospatial database with spatial indexing    |
| **Redis**                  | Caching for geocoding & query optimization   |
| **Next.js 14**             | Modern React frontend with App Router        |
| **TypeScript**             | Type-safe frontend development               |
| **Leaflet.js**             | Interactive maps & geospatial visualization  |
| **Python (Pandas)**        | Data processing & ETL automation             |
| **SCSS Modules**           | Styled components with dark mode theming     |

---

## 📊 Supported Datasets

| **Dataset** | **Icon** | **Records** | **Description** |
|-------------|----------|-------------|-----------------|
| **Meteorites** | ☄️ | 45,000+ | Global meteorite impact and discovery data |
| **Climate** | 🌡️ | 43,000+ | Australian climate station data (temperature, humidity, evaporation) |
| **Wind** | 💨 | 60,000+ | Wind speed and direction observations |
| **Vegetation** | 🌿 | 7+ | Vegetation classification and area data |
| **Infrastructure** | 🏗️ | 1,000+ | Infrastructure and utility data |
| **Fire** | 🔥 | Coming Soon | Fire risk and projection data |

---

## 🚀 Quick Start

### Prerequisites
- **Go** 1.21+ 
- **Node.js** 18+
- **PostgreSQL** 14+ with **PostGIS** extension
- **Redis** (optional, for caching)

### 1. Clone & Setup
```bash
git clone https://github.com/yourusername/GeoGO.git
cd GeoGO
```

### 2. Backend Setup (GeoB)
```bash
cd GeoB

# Install Go dependencies
go mod tidy

# Setup PostgreSQL with PostGIS
# Create database: geogo
# Enable PostGIS extension

# Configure database connection in db/database.go
# Update DSN: host=localhost port=5432 user=postgres password=your_password dbname=geogo

# Create database schema
psql -h localhost -U postgres -d geogo -f utils/SQL/create_unified_schema.sql

# Process and load datasets
python utils/process_datasets.py
psql -h localhost -U postgres -d geogo -f utils/SQL/meteorites.sql
psql -h localhost -U postgres -d geogo -f utils/SQL/climate_avg_temperature.sql
# ... load other datasets as needed

# Run the API server
go run main.go
```

### 3. Frontend Setup (geofe)
```bash
cd geofe

# Install dependencies
npm install

# Start development server
npm run dev
```

### 4. Access the Application
- **Frontend**: http://localhost:3000
- **API**: http://localhost:8080

---

## 📡 API Endpoints

### Core Endpoints
| **Endpoint** | **Method** | **Description** |
|--------------|------------|-----------------|
| `/datasets/types` | GET | Get all dataset types with metadata |
| `/datasets` | GET | Query datasets with filters |
| `/datasets/stats/:type` | GET | Get statistics for specific dataset type |

### Query Parameters
- `type` - Dataset type (meteorite, climate, wind, etc.)
- `value_min` / `value_max` - Value range filtering
- `location` - Location-based search with radius
- `limit` / `offset` - Pagination

### Example Queries
```bash
# Get all meteorites
curl "http://localhost:8080/datasets?type=meteorite&limit=10"

# Get climate data near Melbourne
curl "http://localhost:8080/datasets?type=climate&location=Melbourne&radius=100000"

# Get dataset statistics
curl "http://localhost:8080/datasets/stats/meteorite"
```

---

## 🎨 Frontend Features

### Beautiful Dark Mode Design
- **Modern UI** with glassmorphism effects
- **Gradient backgrounds** and smooth animations
- **Responsive grid layout** for dataset cards
- **Interactive hover effects** and transitions

### Dataset-Specific Pages
- **Dynamic search forms** tailored to each dataset type
- **Interactive maps** with Leaflet.js integration
- **Real-time statistics** and metadata display
- **Geocoding integration** for location names

### User Experience
- **Fast navigation** between dataset categories
- **Live data updates** from the API
- **Error handling** with user-friendly messages
- **Loading states** and progress indicators

---

## 🗄️ Database Schema

### Unified Dataset Table
```sql
CREATE TABLE datasets (
    id SERIAL PRIMARY KEY,
    dataset_type VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    lat DOUBLE PRECISION,
    lon DOUBLE PRECISION,
    value DOUBLE PRECISION,
    unit VARCHAR(50),
    metadata JSONB,
    geom GEOMETRY(POINT, 4326)
);
```

### Spatial Indexing
- **PostGIS spatial indexes** for fast location queries
- **GIST indexes** on geometry columns
- **B-tree indexes** on dataset_type and value columns

---

## 🔧 Development

### Project Structure
```
GeoGO/
├── GeoB/                 # Go backend
│   ├── api/             # API handlers
│   ├── db/              # Database layer
│   ├── models/          # Data models
│   ├── utils/           # Data processing scripts
│   └── main.go          # Server entry point
├── geofe/               # Next.js frontend
│   ├── app/             # App Router pages
│   ├── components/      # React components
│   └── public/          # Static assets
└── README.md           # This file
```

### Adding New Datasets
1. **Process data** in `GeoB/utils/process_datasets.py`
2. **Load SQL** into database
3. **Update frontend** dataset cards in `geofe/app/page.tsx`
4. **Add search fields** in `geofe/components/SearchForm.tsx`

---

## 🚀 Deployment

### Docker (Recommended)
```bash
# Build and run with Docker Compose
docker-compose up -d
```

### Manual Deployment
1. **Build frontend**: `npm run build`
2. **Build backend**: `go build -o geogo main.go`
3. **Setup production database** with PostGIS
4. **Configure environment variables**
5. **Deploy with your preferred hosting**

---

## 🤝 Contributing

We welcome contributions! Here's how to get started:

1. **Fork** the repository
2. **Create** a feature branch: `git checkout -b feature/amazing-feature`
3. **Commit** your changes: `git commit -m 'Add amazing feature'`
4. **Push** to the branch: `git push origin feature/amazing-feature`
5. **Open** a Pull Request

### Development Guidelines
- Follow **Go** and **TypeScript** best practices
- Add **tests** for new features
- Update **documentation** for API changes
- Use **conventional commits** for commit messages

---

## 📈 Roadmap

### 🎯 Upcoming Features
- [ ] **Real-time data streaming** with WebSockets
- [ ] **Advanced analytics** and data visualization
- [ ] **Machine learning** integration for predictions
- [ ] **Mobile app** with React Native
- [ ] **API rate limiting** and authentication
- [ ] **Data export** functionality (CSV, GeoJSON)

### 🔧 Technical Improvements
- [ ] **GraphQL API** for more flexible queries
- [ ] **Microservices architecture** for scalability
- [ ] **Kubernetes deployment** configurations
- [ ] **Performance monitoring** and metrics
- [ ] **Automated testing** pipeline

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🙏 Acknowledgments

- **NASA** for meteorite impact data
- **Australian Bureau of Meteorology** for climate data
- **OpenStreetMap** for mapping data
- **PostGIS** community for spatial database tools
- **Next.js** team for the amazing React framework

---

## 🌟 Star the Repository

If you find GeoGO useful, please give it a ⭐ on GitHub!


package api

import (
	"GeoGO/db"
	"GeoGO/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetAllMeteorites - Handles general search requests with filters and pagination
func GetAllMeteorites(c *gin.Context) {
	limit := 50
	offset := 0

	if c.Query("limit") != "" {
		limit, _ = strconv.Atoi(c.Query("limit"))
	}
	if c.Query("offset") != "" {
		offset, _ = strconv.Atoi(c.Query("offset"))
	}

	log.Printf("📡 Fetching meteorites with limit=%d, offset=%d", limit, offset)

	query := `
		SELECT id, name, recclass, mass, year, ST_AsText(geom) AS location
		FROM locations`

	meteorites, err := FetchMeteorites(c, query, limit, offset)
	if err != nil {
		log.Printf("❌ Failed to fetch meteorites: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data", "details": err.Error()})
		return
	}

	log.Printf("✅ Returning %d meteorites", len(meteorites))
	c.JSON(http.StatusOK, meteorites)
}

// GetLargestMeteorites - Returns the 10 largest meteorites
func GetLargestMeteorites(c *gin.Context) {
	log.Println("📡 Fetching the 10 largest meteorites...")

	query := `
		SELECT name, recclass, mass, year, ST_AsText(geom) AS location
		FROM locations
		ORDER BY mass DESC
		LIMIT 10;
	`

	var meteorites []models.Meteorite
	err := db.DB.Select(&meteorites, query)
	if err != nil {
		log.Printf("❌ Query error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data"})
		return
	}

	log.Printf("✅ Returning %d largest meteorites", len(meteorites))
	c.JSON(http.StatusOK, meteorites)
}

// GetNearbyMeteorites - Returns meteorites near a given location
func GetNearbyMeteorites(c *gin.Context) {
	lat, err1 := strconv.ParseFloat(c.Query("lat"), 64)
	lon, err2 := strconv.ParseFloat(c.Query("lon"), 64)
	radius, err3 := strconv.ParseFloat(c.Query("radius"), 64)

	// Filters for year and mass
	yearStart, err4 := strconv.Atoi(c.DefaultQuery("year_start", "0"))
	yearEnd, err5 := strconv.Atoi(c.DefaultQuery("year_end", "9999"))
	massMin, err6 := strconv.ParseFloat(c.DefaultQuery("mass_min", "0"), 64)
	massMax, err7 := strconv.ParseFloat(c.DefaultQuery("mass_max", "10000000"), 64)

	// Handling invalid parameters
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil || err6 != nil || err7 != nil {
		log.Printf("❌ Invalid query parameters: lat=%f, lon=%f, radius=%f, year_start=%d, year_end=%d, mass_min=%f, mass_max=%f",
			lat, lon, radius, yearStart, yearEnd, massMin, massMax)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
		return
	}

	log.Printf("📡 Fetching meteorites near lat=%f, lon=%f, radius=%f", lat, lon, radius)

	// Query for meteorites near a location with filters
	query := `
		SELECT name, recclass, mass, year, ST_AsText(geom) AS location,
		       ST_Distance(geom::geography, ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography) AS distance_km
		FROM locations
		WHERE ST_DWithin(geom::geography, ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography, $3)
		AND year BETWEEN $4 AND $5
		AND mass BETWEEN $6 AND $7
		ORDER BY distance_km;
	`

	var meteorites []models.Meteorite
	err := db.DB.Select(&meteorites, query, lon, lat, radius, yearStart, yearEnd, massMin, massMax)

	if err != nil {
		log.Printf("❌ Query error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data"})
		return
	}

	log.Printf("✅ Found %d meteorites near given location", len(meteorites))
	c.JSON(http.StatusOK, meteorites)
}

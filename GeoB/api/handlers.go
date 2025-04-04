// handlers.go
//
// HTTP endpoint layer for GeoGO
// Exposes meteorite data search and geocoding endpoints.
// Compliance Level: Moderate
//
// - Handles query parameter parsing and input validation
// - Interfaces with database and geocoding modules
// - Uses structured logging and error propagation
//
// TODO: Implement rate limiting middleware for API endpoints
// TODO: Add request ID tracking for distributed tracing
// TODO: Consider implementing API versioning strategy
// TODO: Add health check and metrics endpoints
// TODO: Implement circuit breaker for geocoding service calls
//
// NOTE: Current pagination implementation may not scale well with large result sets
// NOTE: Consider implementing cursor-based pagination for better performance
//
// Cache Policy:
// - No direct caching at handler level
// - Relies on geocoding service for location caching
// - Database query results are not cached (consider adding Redis layer)

package api

import (
	"GeoGO/api/geocoding"
	"GeoGO/db"
	"GeoGO/models"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetMeteoriteLocation handles reverse geocoding requests to convert coordinates to a human-readable location.
// It validates input coordinates and delegates to the geocoding service for location lookup.
//
// TODO: Add rate limiting middleware to prevent abuse
// TODO: Consider adding request tracing for monitoring
// TODO: Add input validation for coordinate ranges (-90 to 90 for lat, -180 to 180 for lon)
func GetMeteoriteLocation(c *gin.Context) {
	geocoding.GetMeteoriteLocation(c)
}

// GetCoordinatesFromLocation handles forward geocoding requests to convert location names to coordinates.
// It validates the location input and uses the geocoding service to resolve coordinates.
//
// TODO: Add input sanitization for location names
// TODO: Implement location name normalization
// TODO: Add support for fuzzy matching of location names
func GetCoordinatesFromLocation(c *gin.Context) {
	geocoding.GetCoordinatesFromLocation(c)
}

// GetAllMeteorites provides a flexible search endpoint for meteorite data with multiple filter options.
// It supports filtering by year range, mass range, and location proximity, with pagination.
//
// TODO: Implement cursor-based pagination for better performance with large datasets
// TODO: Add support for sorting by multiple fields
// TODO: Consider implementing query result caching for common filter combinations
// TODO: Add support for bulk export of search results
//
// NOTE: Current offset-based pagination may become inefficient with large datasets
// NOTE: Consider implementing materialized views for common filter combinations
func GetAllMeteorites(c *gin.Context) {
	limit := 50
	offset := 0
	yearStart, yearEnd := 0, 9999
	massMin, massMax := 0.0, 10000000.0
	var lat, lon, radius float64
	location := c.Query("location")

	// Parse query parameters
	if v := c.Query("limit"); v != "" {
		limit, _ = strconv.Atoi(v)
	}
	if v := c.Query("offset"); v != "" {
		offset, _ = strconv.Atoi(v)
	}
	if v := c.Query("year_start"); v != "" {
		yearStart, _ = strconv.Atoi(v)
	}
	if v := c.Query("year_end"); v != "" {
		yearEnd, _ = strconv.Atoi(v)
	}
	if v := c.Query("mass_min"); v != "" {
		massMin, _ = strconv.ParseFloat(v, 64)
	}
	if v := c.Query("mass_max"); v != "" {
		massMax, _ = strconv.ParseFloat(v, 64)
	}
	if v := c.Query("radius"); v != "" {
		radius, _ = strconv.ParseFloat(v, 64)
	}

	// Convert location to coords
	var locationFilter string
	if location != "" {
		coords, err := geocoding.ForwardGeocode(location)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get coordinates for location"})
			return
		}
		lat, lon = coords.Lat, coords.Lon
		log.Printf("🌍 Location search: '%s' -> [lat: %.6f, lon: %.6f]", location, lat, lon)
		locationFilter = " AND ST_DWithin(geom::geography, ST_SetSRID(ST_MakePoint($5, $6), 4326)::geography, $7)"
	}

	log.Printf("📡 Fetching meteorites: limit=%d, offset=%d, year=[%d-%d], mass=[%.2f-%.2f], location=%s",
		limit, offset, yearStart, yearEnd, massMin, massMax, location)

	// Construct SQL Query Dynamically
	query := `
		SELECT id, name, recclass, mass, year, ST_X(geom) AS lon, ST_Y(geom) AS lat
		FROM locations
		WHERE year BETWEEN $1 AND $2
		AND mass BETWEEN $3 AND $4
	`
	query += locationFilter //append location filter
	query += " ORDER BY year DESC LIMIT $%d OFFSET $%d"

	// Adjust SQL parameter positions dynamically
	if location != "" {
		query = formatQuery(query, 8, 9) // location adds $5, $6, $7
	} else {
		query = formatQuery(query, 5, 6) // normal $5, $6 for limit/offset
	}

	// Prep query arguments dynamically
	args := []interface{}{yearStart, yearEnd, massMin, massMax}
	if location != "" {
		args = append(args, lon, lat, radius)
	}
	args = append(args, limit, offset)

	// Execute query
	meteorites, err := FetchMeteoritesRaw(c, query, args...)
	if err != nil {
		log.Printf("❌ Failed to fetch meteorites: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data", "details": err.Error()})
		return
	}

	log.Printf("✅ Returning %d meteorites", len(meteorites))
	c.JSON(http.StatusOK, meteorites)
}

// FetchMeteoritesRaw executes a parameterized SQL query and returns the results as meteorite data.
// It handles database operations and error propagation while ensuring proper resource management.
//
// TODO: Add query timeout context
// TODO: Implement query retry logic for transient failures
// TODO: Add support for query cancellation
// TODO: Consider implementing connection pooling metrics
func FetchMeteoritesRaw(c *gin.Context, query string, args ...interface{}) ([]models.Meteorite, error) {
	var meteorites []models.Meteorite
	err := db.DB.Select(&meteorites, query, args...)
	if err != nil {
		return nil, err
	}
	return meteorites, nil
}

// GetLargestMeteorites retrieves the 10 largest meteorites by mass from the database.
// The function implements a simple, optimized query for this specific use case.
//
// TODO: Make the limit configurable via query parameter
// TODO: Add support for filtering by meteorite class
// TODO: Consider implementing result caching for this frequently accessed endpoint
func GetLargestMeteorites(c *gin.Context) {
	log.Println("📡 Fetching the 10 largest meteorites...")
	query := `
		SELECT id, name, recclass, mass, year, ST_X(geom) AS lon, ST_Y(geom) AS lat
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

// GetNearbyMeteorites searches for meteorites within a specified radius of given coordinates.
// It supports additional filters for year range and mass, with proper input validation.
//
// TODO: Add support for different distance units (km, miles)
// TODO: Implement spatial indexing for better performance
// TODO: Add support for complex shapes (polygons, etc.)
// TODO: Consider implementing result caching for common location queries
//
// NOTE: Current implementation uses simple distance calculation
// NOTE: Consider using PostGIS spatial functions for more complex queries
func GetNearbyMeteorites(c *gin.Context) {
	lat, err1 := strconv.ParseFloat(c.Query("lat"), 64)
	lon, err2 := strconv.ParseFloat(c.Query("lon"), 64)
	radius, err3 := strconv.ParseFloat(c.Query("radius"), 64)

	yearStart, err4 := strconv.Atoi(c.DefaultQuery("year_start", "0"))
	yearEnd, err5 := strconv.Atoi(c.DefaultQuery("year_end", "9999"))
	massMin, err6 := strconv.ParseFloat(c.DefaultQuery("mass_min", "0"), 64)
	massMax, err7 := strconv.ParseFloat(c.DefaultQuery("mass_max", "10000000"), 64)

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil || err6 != nil || err7 != nil {
		log.Printf("❌ Invalid query parameters")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
		return
	}

	log.Printf("📡 Fetching meteorites near lat=%.6f, lon=%.6f, radius=%.2f km", lat, lon, radius)

	query := `
		SELECT id, name, recclass, mass, year, ST_X(geom) AS lon, ST_Y(geom) AS lat
		FROM locations
		WHERE ST_DWithin(geom::geography, ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography, $3)
		AND year BETWEEN $4 AND $5
		AND mass BETWEEN $6 AND $7
	`
	args := []interface{}{lon, lat, radius, yearStart, yearEnd, massMin, massMax}

	meteorites, err := FetchMeteoritesRaw(c, query, args...)
	if err != nil {
		log.Printf("❌ Failed to fetch meteorites: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data", "details": err.Error()})
		return
	}

	log.Printf("✅ Found %d meteorites near given location", len(meteorites))
	c.JSON(http.StatusOK, meteorites)
}

// formatQuery updates SQL query placeholders for pagination parameters.
// It ensures proper parameter indexing when combining multiple filters.
//
// TODO: Add validation for parameter positions
// TODO: Consider using named parameters instead of positional
// TODO: Add support for query template caching
func formatQuery(query string, limitPos, offsetPos int) string {
	return fmt.Sprintf(query, limitPos, offsetPos)
}

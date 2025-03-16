package api

import (
	"GeoGO/api/geocoding"
	"GeoGO/db"
	"GeoGO/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// reverse Geocode (Coords → Place)
func GetMeteoriteLocation(c *gin.Context) {
	geocoding.GetMeteoriteLocation(c)
}

// forward Geocode (Place → Coords)
func GetCoordinatesFromLocation(c *gin.Context) {
	geocoding.GetCoordinatesFromLocation(c)
}

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

	// convert place name to coordinates
	if location != "" {
		coords, err := geocoding.ForwardGeocode(location)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get coordinates for location"})
			return
		}
		lat, lon = coords.Lat, coords.Lon
		log.Printf("🌍 Location search: '%s' -> [lat: %.6f, lon: %.6f]", location, lat, lon)
	}

	log.Printf("📡 Fetching meteorites: limit=%d, offset=%d, year=[%d-%d], mass=[%.2f-%.2f], location=%s",
		limit, offset, yearStart, yearEnd, massMin, massMax, location)

	// Base Query
	query := `
		SELECT id, name, recclass, mass, year, ST_X(geom) AS lon, ST_Y(geom) AS lat
		FROM locations
		WHERE year BETWEEN $1 AND $2
		AND mass BETWEEN $3 AND $4
	`

	var args []interface{}
	args = append(args, yearStart, yearEnd, massMin, massMax)

	if location != "" {
		query += ` AND ST_DWithin(geom::geography, ST_SetSRID(ST_MakePoint($5, $6), 4326)::geography, $7)`
		args = append(args, lon, lat, radius)
	}

	query += " ORDER BY year DESC LIMIT $8 OFFSET $9"
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

func FetchMeteoritesRaw(c *gin.Context, query string, args ...interface{}) ([]models.Meteorite, error) {
	var meteorites []models.Meteorite
	err := db.DB.Select(&meteorites, query, args...)
	if err != nil {
		return nil, err
	}
	return meteorites, nil
}

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

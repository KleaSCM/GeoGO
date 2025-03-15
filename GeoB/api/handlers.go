package api

import (
	"GeoGO/db"
	"GeoGO/models"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

func GetAllMeteorites(c *gin.Context) {
	limit := 50
	offset := 0
	yearStart, yearEnd := 0, 9999
	massMin, massMax := 0.0, 10000000.0

	// Parse Query Parameters
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

	log.Printf("📡 Fetching meteorites with limit=%d, offset=%d, year=[%d-%d], mass=[%.2f-%.2f]", limit, offset, yearStart, yearEnd, massMin, massMax)

	query := `
		SELECT id, name, recclass, mass, year, ST_AsText(geom) AS location
		FROM locations
		WHERE year BETWEEN $1 AND $2
		AND mass BETWEEN $3 AND $4
		ORDER BY year DESC
		LIMIT $5 OFFSET $6;
	`

	meteorites, err := FetchMeteoritesRaw(c, query, yearStart, yearEnd, massMin, massMax, limit, offset)
	if err != nil {
		log.Printf("❌ Failed to fetch meteorites: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data", "details": err.Error()})
		return
	}
	log.Printf("✅ Returning %d meteorites", len(meteorites))
	c.JSON(http.StatusOK, meteorites)
}

// Return the 10 largest meteorites
func GetLargestMeteorites(c *gin.Context) {
	log.Println("📡 Fetching the 10 largest meteorites...")

	query := `
		SELECT id, name, recclass, mass, year, ST_AsText(geom) AS location
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
	// Filters for year and mass
	yearStart, err4 := strconv.Atoi(c.DefaultQuery("year_start", "0"))
	yearEnd, err5 := strconv.Atoi(c.DefaultQuery("year_end", "9999"))
	massMin, err6 := strconv.ParseFloat(c.DefaultQuery("mass_min", "0"), 64)
	massMax, err7 := strconv.ParseFloat(c.DefaultQuery("mass_max", "10000000"), 64)
	// Handle invalid query parameters
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil || err6 != nil || err7 != nil {
		log.Printf("❌ Invalid query parameters: lat=%f, lon=%f, radius=%f, year_start=%d, year_end=%d, mass_min=%f, mass_max=%f",
			lat, lon, radius, yearStart, yearEnd, massMin, massMax)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
		return
	}
	log.Printf("📡 Fetching meteorites near lat=%f, lon=%f, radius=%f", lat, lon, radius)
	var wg sync.WaitGroup
	resultsChan := make(chan []models.Meteorite, 3)
	errorsChan := make(chan error, 3)

	queries := []struct {
		query string
		args  []interface{}
	}{
		{
			query: "SELECT id, name, recclass, mass, year, ST_AsText(geom) AS location FROM locations WHERE ST_DWithin(geom::geography, ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography, $3)",
			args:  []interface{}{lon, lat, radius},
		},
		{
			query: "SELECT id, name, recclass, mass, year, ST_AsText(geom) AS location FROM locations WHERE year BETWEEN $1 AND $2",
			args:  []interface{}{yearStart, yearEnd},
		},
		{
			query: "SELECT id, name, recclass, mass, year, ST_AsText(geom) AS location FROM locations WHERE mass BETWEEN $1 AND $2",
			args:  []interface{}{massMin, massMax},
		},
	}
	for _, q := range queries {
		wg.Add(1)
		go func(query string, args []interface{}) {
			defer wg.Done()
			data, err := FetchMeteoritesRaw(c, query, args...)
			if err != nil {
				errorsChan <- err
				return
			}
			resultsChan <- data
		}(q.query, q.args)
	}
	wg.Wait()
	close(resultsChan)
	close(errorsChan)

	// Merge results
	var mergedResults []models.Meteorite
	for r := range resultsChan {
		mergedResults = append(mergedResults, r...)
	}
	for err := range errorsChan {
		log.Printf("❌ Error fetching data: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data"})
		return
	}
	log.Printf("✅ Found %d meteorites near given location", len(mergedResults))
	c.JSON(http.StatusOK, mergedResults)
}
func FetchMeteoritesRaw(c *gin.Context, query string, args ...interface{}) ([]models.Meteorite, error) {
	var meteorites []models.Meteorite
	err := db.DB.Select(&meteorites, query, args...)
	if err != nil {
		return nil, err
	}
	return meteorites, nil
}

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

// GetDatasets provides a unified endpoint for all dataset types
func GetDatasets(c *gin.Context) {
	datasetType := c.Query("type")
	limit := 50
	offset := 0
	valueMin, valueMax := 0.0, 10000000.0
	var lat, lon, radius float64
	location := c.Query("location")

	// Parse query parameters
	if v := c.Query("limit"); v != "" {
		limit, _ = strconv.Atoi(v)
	}
	if v := c.Query("offset"); v != "" {
		offset, _ = strconv.Atoi(v)
	}
	if v := c.Query("value_min"); v != "" {
		valueMin, _ = strconv.ParseFloat(v, 64)
	}
	if v := c.Query("value_max"); v != "" {
		valueMax, _ = strconv.ParseFloat(v, 64)
	}
	if v := c.Query("radius"); v != "" {
		radius, _ = strconv.ParseFloat(v, 64)
	}

	// Build query
	query := `
		SELECT id, dataset_type, name, lat, lon, value, unit, metadata, 
		       recclass, mass, year, nametype, fall
		FROM datasets
		WHERE value BETWEEN $1 AND $2
	`
	args := []interface{}{valueMin, valueMax}
	paramCount := 2

	// Add dataset type filter
	if datasetType != "" {
		paramCount++
		query += fmt.Sprintf(" AND dataset_type = $%d", paramCount)
		args = append(args, datasetType)
	}

	// Add location filter
	if location != "" {
		coords, err := geocoding.ForwardGeocode(location)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get coordinates for location"})
			return
		}
		lat, lon = coords.Lat, coords.Lon
		paramCount++
		query += fmt.Sprintf(" AND ST_DWithin(geom::geography, ST_SetSRID(ST_MakePoint($%d, $%d), 4326)::geography, $%d)",
			paramCount+1, paramCount+2, paramCount+3)
		args = append(args, lon, lat, radius)
		paramCount += 3
	}

	query += fmt.Sprintf(" ORDER BY id DESC LIMIT $%d OFFSET $%d", paramCount+1, paramCount+2)
	args = append(args, limit, offset)

	// Execute query
	datasets := make([]models.Dataset, 0)
	err := db.DB.Select(&datasets, query, args...)
	if err != nil {
		log.Printf("❌ Failed to fetch datasets: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data"})
		return
	}

	log.Printf("✅ Returning %d datasets", len(datasets))
	c.JSON(http.StatusOK, datasets)
}

// GetDatasetTypes returns information about available dataset types
func GetDatasetTypes(c *gin.Context) {
	query := `
		SELECT 
			dataset_type,
			COUNT(*) as count,
			MIN(value) as min_value,
			MAX(value) as max_value,
			MIN(timestamp) as min_date,
			MAX(timestamp) as max_date
		FROM datasets 
		GROUP BY dataset_type
	`

	var results []struct {
		Type     string   `db:"dataset_type"`
		Count    int      `db:"count"`
		MinValue *float64 `db:"min_value"`
		MaxValue *float64 `db:"max_value"`
		MinDate  *string  `db:"min_date"`
		MaxDate  *string  `db:"max_date"`
	}

	err := db.DB.Select(&results, query)
	if err != nil {
		log.Printf("❌ Failed to fetch dataset types: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch dataset types"})
		return
	}

	var datasetInfos []models.DatasetInfo
	for _, result := range results {
		info := models.DatasetInfo{
			Type:  models.DatasetType(result.Type),
			Count: result.Count,
		}

		// Set description based on type
		switch info.Type {
		case models.DatasetTypeMeteorite:
			info.Name = "Meteorite Landings"
			info.Description = "Global meteorite impact and discovery data"
		case models.DatasetTypeClimate:
			info.Name = "Climate Stations"
			info.Description = "Australian climate station data with temperature, humidity, and evaporation"
		case models.DatasetTypeWind:
			info.Name = "Wind Observations"
			info.Description = "Wind speed and direction observations"
		case models.DatasetTypeVegetation:
			info.Name = "Vegetation Zones"
			info.Description = "Vegetation classification and area data"
		case models.DatasetTypeInfrastructure:
			info.Name = "Infrastructure"
			info.Description = "Infrastructure and utility data"
		case models.DatasetTypeFire:
			info.Name = "Fire Projections"
			info.Description = "Fire risk and projection data"
		default:
			info.Name = string(info.Type)
			info.Description = "Geospatial dataset"
		}

		// Format value range
		if result.MinValue != nil && result.MaxValue != nil {
			info.ValueRange = fmt.Sprintf("%.2f - %.2f", *result.MinValue, *result.MaxValue)
		}

		// Format date range
		if result.MinDate != nil && result.MaxDate != nil {
			info.DateRange = fmt.Sprintf("%s to %s", *result.MinDate, *result.MaxDate)
		}

		datasetInfos = append(datasetInfos, info)
	}

	log.Printf("✅ Returning %d dataset types", len(datasetInfos))
	c.JSON(http.StatusOK, datasetInfos)
}

// GetDatasetStats returns statistics for a specific dataset type
func GetDatasetStats(c *gin.Context) {
	datasetType := c.Param("type")
	if datasetType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dataset type is required"})
		return
	}

	query := `
		SELECT 
			COUNT(*) as total_count,
			AVG(value) as avg_value,
			MIN(value) as min_value,
			MAX(value) as max_value,
			ST_Extent(geom) as spatial_extent
		FROM datasets 
		WHERE dataset_type = $1
	`

	var stats struct {
		TotalCount    int     `db:"total_count" json:"total_count"`
		AvgValue      float64 `db:"avg_value" json:"avg_value"`
		MinValue      float64 `db:"min_value" json:"min_value"`
		MaxValue      float64 `db:"max_value" json:"max_value"`
		SpatialExtent string  `db:"spatial_extent" json:"spatial_extent"`
	}

	err := db.DB.Get(&stats, query, datasetType)
	if err != nil {
		log.Printf("❌ Failed to fetch dataset stats: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch dataset stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetDatasetsByType returns datasets of a specific type with filtering
func GetDatasetsByType(c *gin.Context) {
	datasetType := c.Param("type")
	if datasetType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dataset type is required"})
		return
	}

	// Add type filter to the main query
	c.Set("dataset_type", datasetType)
	GetDatasets(c)
}

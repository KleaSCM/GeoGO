package api

import (
	"GeoGO/db"
	"GeoGO/models"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

func buildQueryFilters(c *gin.Context) (string, []interface{}, error) {
	var filters []string
	var args []interface{}
	paramIndex := 1
	// Extract query parameters
	yearStart, err1 := strconv.Atoi(c.DefaultQuery("year_start", "0"))
	yearEnd, err2 := strconv.Atoi(c.DefaultQuery("year_end", "9999"))
	recclass := c.DefaultQuery("recclass", "")
	massMin, err3 := strconv.ParseFloat(c.DefaultQuery("mass_min", "0"), 64)
	massMax, err4 := strconv.ParseFloat(c.DefaultQuery("mass_max", "10000000"), 64)
	location := c.DefaultQuery("location", "")
	// Validate parameters
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		return "", nil, fmt.Errorf("invalid query parameters")
	}
	// Year filter
	filters = append(filters, fmt.Sprintf("year BETWEEN $%d AND $%d", paramIndex, paramIndex+1))
	args = append(args, yearStart, yearEnd)
	paramIndex += 2
	// Recclass filter
	if recclass != "" {
		filters = append(filters, fmt.Sprintf("recclass = $%d", paramIndex))
		args = append(args, recclass)
		paramIndex++
	}
	// Mass filter
	filters = append(filters, fmt.Sprintf("mass BETWEEN $%d AND $%d", paramIndex, paramIndex+1))
	args = append(args, massMin, massMax)
	paramIndex += 2
	// Location filter
	if location != "" {
		coords := strings.Split(location, ",")
		if len(coords) == 2 {
			lat, errLat := strconv.ParseFloat(strings.TrimSpace(coords[0]), 64)
			lon, errLon := strconv.ParseFloat(strings.TrimSpace(coords[1]), 64)
			if errLat == nil && errLon == nil {
				filters = append(filters, fmt.Sprintf("ST_DWithin(geom::geography, ST_SetSRID(ST_MakePoint($%d, $%d), 4326)::geography, $%d)", paramIndex, paramIndex+1, paramIndex+2))
				args = append(args, lon, lat, 50000) // 50km default radius
				paramIndex += 3
			} else {
				return "", nil, fmt.Errorf("invalid location format, expected lat,lon")
			}
		} else {
			return "", nil, fmt.Errorf("location query is invalid")
		}
	}

	whereClause := ""
	if len(filters) > 0 {
		whereClause = " WHERE " + strings.Join(filters, " AND ")
	}

	log.Printf("🔍 Generated SQL WHERE clause: %s", whereClause)
	return whereClause, args, nil
}
func FetchMeteorites(c *gin.Context, query string, limit, offset int) ([]models.Meteorite, error) {
	whereClause, args, err := buildQueryFilters(c)
	if err != nil {
		return nil, fmt.Errorf("error building filters: %v", err)
	}

	queryWithFilters := query + whereClause + fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	args = append(args, limit, offset)

	log.Printf("📡 Executing Query: %s", queryWithFilters)

	var meteorites []models.Meteorite
	err = db.DB.Select(&meteorites, queryWithFilters, args...)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}

	log.Printf("✅ Query successful, found %d results", len(meteorites))
	return meteorites, nil
}

// Concurrent Query Execution
func FetchParallel(c *gin.Context, queries []string) [][]models.Meteorite {
	var wg sync.WaitGroup
	results := make([][]models.Meteorite, len(queries))
	ch := make(chan struct {
		index int
		data  []models.Meteorite
		err   error
	}, len(queries))
	// Execute each query in parallel
	for i, query := range queries {
		wg.Add(1)
		go func(i int, query string) {
			defer wg.Done()
			data, err := FetchMeteorites(c, query, 50, 0)
			ch <- struct {
				index int
				data  []models.Meteorite
				err   error
			}{i, data, err}
		}(i, query)
	}
	wg.Wait()
	close(ch)
	for res := range ch {
		results[res.index] = res.data
	}

	return results
}

// package api

// import (
// 	"GeoGO/db"
// 	"GeoGO/models"
// 	"fmt"
// 	"log"
// 	"strconv"
// 	"strings"

// 	"github.com/gin-gonic/gin"
// )

// // buildQueryFilters - Constructs SQL WHERE clause with dynamic parameters
// func buildQueryFilters(c *gin.Context) (string, []interface{}, error) {
// 	var filters []string
// 	var args []interface{}
// 	paramIndex := 1 // SQL placeholder counter

// 	// Extract query parameters
// 	yearStart, err1 := strconv.Atoi(c.DefaultQuery("year_start", "0"))
// 	yearEnd, err2 := strconv.Atoi(c.DefaultQuery("year_end", "9999"))
// 	recclass := c.DefaultQuery("recclass", "")
// 	massMin, err3 := strconv.ParseFloat(c.DefaultQuery("mass_min", "0"), 64)
// 	massMax, err4 := strconv.ParseFloat(c.DefaultQuery("mass_max", "10000000"), 64)
// 	location := c.DefaultQuery("location", "")

// 	// Check for any invalid parameters
// 	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
// 		return "", nil, fmt.Errorf("invalid query parameters")
// 	}

// 	// Year filter
// 	filters = append(filters, fmt.Sprintf("year BETWEEN $%d AND $%d", paramIndex, paramIndex+1))
// 	args = append(args, yearStart, yearEnd)
// 	paramIndex += 2

// 	// Recclass filter
// 	if recclass != "" {
// 		filters = append(filters, fmt.Sprintf("recclass = $%d", paramIndex))
// 		args = append(args, recclass)
// 		paramIndex++
// 	}

// 	// Mass filter
// 	filters = append(filters, fmt.Sprintf("mass BETWEEN $%d AND $%d", paramIndex, paramIndex+1))
// 	args = append(args, massMin, massMax)
// 	paramIndex += 2

// 	// Location filter (if provided)
// 	if location != "" {
// 		coords := strings.Split(location, ",")
// 		if len(coords) == 2 {
// 			lat, errLat := strconv.ParseFloat(strings.TrimSpace(coords[0]), 64)
// 			lon, errLon := strconv.ParseFloat(strings.TrimSpace(coords[1]), 64)
// 			if errLat == nil && errLon == nil {
// 				filters = append(filters, fmt.Sprintf("ST_DWithin(geom::geography, ST_SetSRID(ST_MakePoint($%d, $%d), 4326)::geography, $%d)", paramIndex, paramIndex+1, paramIndex+2))
// 				args = append(args, lon, lat, 50000) // Default radius = 50km
// 				paramIndex += 3
// 			} else {
// 				return "", nil, fmt.Errorf("invalid location format, expected lat,lon")
// 			}
// 		} else {
// 			return "", nil, fmt.Errorf("location query is invalid")
// 		}
// 	}

// 	// Build WHERE clause
// 	whereClause := ""
// 	if len(filters) > 0 {
// 		whereClause = " WHERE " + strings.Join(filters, " AND ")
// 	}

// 	log.Printf("🔍 Generated SQL WHERE clause: %s", whereClause)
// 	return whereClause, args, nil
// }

// // FetchMeteorites retrieves meteorites based on user filters
// func FetchMeteorites(c *gin.Context, query string, limit, offset int) ([]models.Meteorite, error) {
// 	whereClause, args, err := buildQueryFilters(c)
// 	if err != nil {
// 		return nil, fmt.Errorf("error building filters: %v", err)
// 	}

// 	queryWithFilters := query + whereClause + fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
// 	args = append(args, limit, offset)

// 	log.Printf("📡 Executing Query: %s", queryWithFilters)

// 	var meteorites []models.Meteorite
// 	err = db.DB.Select(&meteorites, queryWithFilters, args...)
// 	if err != nil {
// 		return nil, fmt.Errorf("query error: %v", err)
// 	}

// 	log.Printf("✅ Query successful, found %d results", len(meteorites))
// 	return meteorites, nil
// }

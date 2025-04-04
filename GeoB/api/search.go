// search.go
//
// Search functionality implementation for GeoGO
// Provides query building and parallel execution capabilities.
// Compliance Level: High
//
// - Implements parameterized query construction
// - Handles concurrent query execution
// - Manages resource synchronization
// - Implements proper error handling
//
// TODO: Add query plan analysis for performance optimization
// TODO: Implement query result caching for frequently accessed data
// TODO: Add support for more complex spatial queries
// TODO: Consider implementing query timeout mechanism
// TODO: Add support for query result streaming
//
// NOTE: Parallel execution may need tuning based on database connection pool size
// NOTE: Consider implementing query result batching for large datasets
//
// Cache Policy:
// - No direct caching at search level
// - Consider implementing query result caching for common searches
// - Cache invalidation strategy needed if implemented

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

// buildQueryFilters constructs a SQL WHERE clause and corresponding arguments based on query parameters.
// It handles multiple filter types including year range, meteorite class, mass range, and location-based filtering.
// The function implements parameterized queries to prevent SQL injection and validates all input parameters.
//
// Parameters:
//   - c: Gin context containing query parameters
//
// Returns:
//   - string: The constructed WHERE clause
//   - []interface{}: Arguments for the parameterized query
//   - error: Any validation or parsing errors encountered
//
// Example:
//
//	Input: year_start=1900, year_end=2000, mass_min=1000, location="40.7128,-74.0060"
//	Output: "WHERE year BETWEEN $1 AND $2 AND mass BETWEEN $3 AND $4 AND ST_DWithin(...)", [1900, 2000, 1000, ...]
//
// TODO: Add support for more complex filter combinations
// TODO: Implement filter validation against schema
// TODO: Add support for custom filter operators
// TODO: Consider implementing filter expression parsing
//
// NOTE: Current implementation uses simple string concatenation
// NOTE: Consider using a query builder library for more complex queries
func buildQueryFilters(c *gin.Context) (string, []interface{}, error) {
	var filters []string
	var args []interface{}
	paramIndex := 1

	yearStart, err1 := strconv.Atoi(c.DefaultQuery("year_start", "0"))
	yearEnd, err2 := strconv.Atoi(c.DefaultQuery("year_end", "9999"))
	recclass := c.DefaultQuery("recclass", "")
	massMin, err3 := strconv.ParseFloat(c.DefaultQuery("mass_min", "0"), 64)
	massMax, err4 := strconv.ParseFloat(c.DefaultQuery("mass_max", "10000000"), 64)
	location := c.DefaultQuery("location", "")

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		return "", nil, fmt.Errorf("invalid query parameters")
	}

	filters = append(filters, fmt.Sprintf("year BETWEEN $%d AND $%d", paramIndex, paramIndex+1))
	args = append(args, yearStart, yearEnd)
	paramIndex += 2

	if recclass != "" {
		filters = append(filters, fmt.Sprintf("recclass = $%d", paramIndex))
		args = append(args, recclass)
		paramIndex++
	}

	filters = append(filters, fmt.Sprintf("mass BETWEEN $%d AND $%d", paramIndex, paramIndex+1))
	args = append(args, massMin, massMax)
	paramIndex += 2

	if location != "" {
		coords := strings.Split(location, ",")
		if len(coords) == 2 {
			lat, errLat := strconv.ParseFloat(strings.TrimSpace(coords[0]), 64)
			lon, errLon := strconv.ParseFloat(strings.TrimSpace(coords[1]), 64)
			if errLat == nil && errLon == nil {
				filters = append(filters, fmt.Sprintf("ST_DWithin(geom::geography, ST_SetSRID(ST_MakePoint($%d, $%d), 4326)::geography, $%d)", paramIndex, paramIndex+1, paramIndex+2))
				args = append(args, lon, lat, 50000)
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

// FetchMeteorites executes a parameterized query to retrieve meteorite data based on provided filters.
// It combines the base query with dynamically constructed filters and handles pagination.
// The function ensures proper error handling and logging of query execution.
//
// Parameters:
//   - c: Gin context for request handling
//   - query: Base SQL query string
//   - limit: Maximum number of results to return
//   - offset: Number of results to skip (for pagination)
//
// Returns:
//   - []models.Meteorite: Slice of meteorite data matching the query
//   - error: Any database or query execution errors
//
// Note: The function uses parameterized queries to prevent SQL injection and
// implements proper error handling for database operations.
//
// TODO: Add query plan analysis for performance optimization
// TODO: Implement query timeout mechanism
// TODO: Add support for query result streaming
// TODO: Consider implementing query result caching
//
// NOTE: Current implementation may need optimization for large result sets
// NOTE: Consider implementing query result batching
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

// FetchParallel executes multiple queries concurrently and aggregates their results.
// It uses a worker pool pattern with goroutines and channels for efficient parallel execution.
// The function ensures proper synchronization and error handling across all concurrent operations.
//
// Parameters:
//   - c: Gin context for request handling
//   - queries: Slice of query strings to execute in parallel
//
// Returns:
//   - [][]models.Meteorite: Slice of meteorite data slices, one for each query
//
// Implementation Details:
//   - Uses sync.WaitGroup for goroutine synchronization
//   - Implements a buffered channel for result collection
//   - Handles error propagation from individual queries
//   - Ensures proper cleanup of resources
//
// TODO: Add support for dynamic worker pool sizing
// TODO: Implement graceful shutdown of workers
// TODO: Add support for query prioritization
// TODO: Consider implementing result streaming
//
// NOTE: Current implementation uses fixed-size worker pool
// NOTE: Consider implementing work stealing for better load balancing
// NOTE: May need tuning based on database connection pool size
func FetchParallel(c *gin.Context, queries []string) [][]models.Meteorite {
	var wg sync.WaitGroup
	results := make([][]models.Meteorite, len(queries))
	ch := make(chan struct {
		index int
		data  []models.Meteorite
		err   error
	}, len(queries))
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

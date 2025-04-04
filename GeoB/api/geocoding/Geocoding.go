package geocoding

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// Redis client configuration for caching geocoding results.
// Uses a local Redis instance with default settings for development.
// In production, consider using environment variables for configuration.
var redisClient = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})

// ReverseGeocodeResponse represents the structure of the response from the Nominatim reverse geocoding API.
// It contains the human-readable location name for given coordinates.
type ReverseGeocodeResponse struct {
	DisplayName string `json:"display_name"`
}

// ForwardGeocodeResponse represents the structure of the response from the Nominatim forward geocoding API.
// It contains the coordinates (latitude and longitude) for a given location name.
type ForwardGeocodeResponse struct {
	Lat float64 `json:"lat,string"`
	Lon float64 `json:"lon,string"`
}

// GetMeteoriteLocation handles HTTP requests for reverse geocoding.
// It validates input coordinates and returns a human-readable location name.
//
// Query Parameters:
//   - lat: Latitude coordinate
//   - lon: Longitude coordinate
//
// Response:
//   - 200 OK: Location information in JSON format
//   - 400 Bad Request: Invalid or missing coordinates
//   - 500 Internal Server Error: Geocoding service failure
func GetMeteoriteLocation(c *gin.Context) {
	latStr := c.Query("lat")
	lonStr := c.Query("lon")
	if latStr == "" || lonStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing latitude or longitude"})
		return
	}
	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid latitude"})
		return
	}
	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid longitude"})
		return
	}
	location, err := ReverseGeocode(lat, lon)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Reverse geocoding failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"lat":      lat,
		"lon":      lon,
		"location": location,
	})
}

// ReverseGeocode converts coordinates to a human-readable location name using the Nominatim API.
// It implements caching using Redis to reduce API calls and improve response times.
//
// Parameters:
//   - lat: Latitude coordinate
//   - lon: Longitude coordinate
//
// Returns:
//   - string: Human-readable location name
//   - error: Any geocoding or caching errors
//
// Implementation Details:
//   - Uses Redis for caching with a 24-hour TTL
//   - Implements a 10-second timeout for API requests
//   - Provides fallback to coordinate string on API failure
func ReverseGeocode(lat, lon float64) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	redisKey := fmt.Sprintf("geo:%f,%f", lat, lon)
	cached, err := redisClient.Get(ctx, redisKey).Result()
	if err == nil && cached != "" {
		log.Printf("🗺️ Cache Hit: %s -> %s", redisKey, cached)
		return cached, nil
	} else if err != redis.Nil {
		log.Printf("⚠️ Redis error: %v", err)
	}
	url := fmt.Sprintf("https://nominatim.openstreetmap.org/reverse?format=json&lat=%f&lon=%f", lat, lon)
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		log.Println("❌ Reverse geocoding API request failed:", err)
		return fmt.Sprintf("Coordinates: %.4f, %.4f", lat, lon), nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("❌ Nominatim API returned status %d", resp.StatusCode)
		return fmt.Sprintf("Coordinates: %.4f, %.4f", lat, lon), nil
	}
	var result ReverseGeocodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("❌ Error decoding API response:", err)
		return fmt.Sprintf("Coordinates: %.4f, %.4f", lat, lon), nil
	}
	if result.DisplayName == "" {
		result.DisplayName = fmt.Sprintf("Coordinates: %.4f, %.4f", lat, lon)
	}
	if err := redisClient.Set(ctx, redisKey, result.DisplayName, 24*time.Hour).Err(); err != nil {
		log.Printf("⚠️ Redis cache save failed: %v", err)
	} else {
		log.Printf("🌍 Fetched from API & cached: %s -> %s", redisKey, result.DisplayName)
	}
	return result.DisplayName, nil
}

// GetCoordinatesFromLocation handles HTTP requests for forward geocoding.
// It converts a location name to coordinates using the Nominatim API.
//
// Query Parameters:
//   - location: Human-readable location name
//
// Response:
//   - 200 OK: Coordinates in JSON format
//   - 400 Bad Request: Missing location parameter
//   - 500 Internal Server Error: Geocoding service failure
func GetCoordinatesFromLocation(c *gin.Context) {
	location := c.Query("location")
	if location == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing location name"})
		return
	}
	coords, err := ForwardGeocode(location)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Forward geocoding failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"location": location,
		"lat":      coords.Lat,
		"lon":      coords.Lon,
	})
}

// ForwardGeocode converts a location name to coordinates using the Nominatim API.
// It implements caching using Redis to reduce API calls and improve response times.
//
// Parameters:
//   - location: Human-readable location name
//
// Returns:
//   - *ForwardGeocodeResponse: Coordinates (latitude and longitude)
//   - error: Any geocoding or caching errors
//
// Implementation Details:
//   - Uses Redis for caching with a 24-hour TTL
//   - Implements a 10-second timeout for API requests
//   - Handles multiple results by returning the first match
func ForwardGeocode(location string) (*ForwardGeocodeResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	redisKey := fmt.Sprintf("geo:%s", strings.ToLower(location))
	cached, err := redisClient.Get(ctx, redisKey).Result()
	if err == nil && cached != "" {
		log.Printf("🗺️ Cache Hit: %s -> %s", redisKey, cached)
		parts := strings.Split(cached, ",")
		if len(parts) == 2 {
			lat, _ := strconv.ParseFloat(parts[0], 64)
			lon, _ := strconv.ParseFloat(parts[1], 64)
			return &ForwardGeocodeResponse{Lat: lat, Lon: lon}, nil
		}
	} else if err != redis.Nil {
		log.Printf("⚠️ Redis error: %v", err)
	}
	url := fmt.Sprintf("https://nominatim.openstreetmap.org/search?format=json&q=%s", url.QueryEscape(location))
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		log.Println("❌ Forward geocoding API request failed:", err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("❌ Nominatim API returned status %d", resp.StatusCode)
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}
	var results []ForwardGeocodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil || len(results) == 0 {
		log.Println("❌ Error decoding API response:", err)
		return nil, fmt.Errorf("Failed to parse API response")
	}
	redisValue := fmt.Sprintf("%f,%f", results[0].Lat, results[0].Lon)
	if err := redisClient.Set(ctx, redisKey, redisValue, 24*time.Hour).Err(); err != nil {
		log.Printf("⚠️ Redis cache save failed: %v", err)
	} else {
		log.Printf("🌍 Fetched from API & cached: %s -> %s", redisKey, redisValue)
	}
	return &results[0], nil
}

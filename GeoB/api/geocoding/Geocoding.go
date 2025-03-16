package geocoding

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var redisClient = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})

type ReverseGeocodeResponse struct {
	DisplayName string `json:"display_name"`
}

// GetMeteoriteLocation handles reverse geocoding requests and sends back the location name.
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

// ReverseGeocode fetches a human-readable location for given coordinates, using caching.
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

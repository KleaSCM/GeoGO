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

func ReverseGeocode(lat, lon float64) (string, error) {
	ctx := context.Background()
	redisKey := fmt.Sprintf("geo:%f,%f", lat, lon)
	cached, err := redisClient.Get(ctx, redisKey).Result()
	if err == nil && cached != "" {
		log.Printf("🗺️ Cache Hit: %s -> %s", redisKey, cached)
		return cached, nil
	} else if err != redis.Nil {
		log.Printf("⚠️ Redis error: %v", err)
	}
	url := fmt.Sprintf("https://nominatim.openstreetmap.org/reverse?format=json&lat=%f&lon=%f", lat, lon)
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		log.Println("❌ Reverse geocoding API request failed:", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("❌ Nominatim API returned status %d", resp.StatusCode)
		return "", fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var result struct {
		DisplayName string `json:"display_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("❌ Error decoding API response:", err)
		return "", err
	}
	if result.DisplayName != "" {
		err := redisClient.Set(ctx, redisKey, result.DisplayName, 24*time.Hour).Err()
		if err != nil {
			log.Printf("⚠️ Redis cache save failed: %v", err)
		} else {
			log.Printf("🌍 Fetched from API & cached: %s -> %s", redisKey, result.DisplayName)
		}
	} else {
		log.Println("⚠️ Nominatim API returned an empty location name")
	}

	return result.DisplayName, nil
}
func GetMeteoriteLocation(c *gin.Context) {
	lat, err1 := strconv.ParseFloat(c.Query("lat"), 64)
	lon, err2 := strconv.ParseFloat(c.Query("lon"), 64)

	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid latitude/longitude"})
		return
	}

	location, err := ReverseGeocode(lat, lon)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Geocoding failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"lat": lat, "lon": lon, "location": location})
}

package main

import (
	"GeoGO/api"
	"GeoGO/db"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()
	r := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://127.0.0.1:3000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(config))

	// Legacy meteorite endpoints (backward compatibility)
	r.GET("/meteorites", api.GetAllMeteorites)
	r.GET("/meteorites/largest", api.GetLargestMeteorites)
	r.GET("/meteorites/nearby", api.GetNearbyMeteorites)
	r.GET("/meteorites/location", api.GetMeteoriteLocation)

	// New unified dataset endpoints
	r.GET("/datasets", api.GetDatasets)
	r.GET("/datasets/types", api.GetDatasetTypes)
	r.GET("/datasets/stats/:type", api.GetDatasetStats)
	r.GET("/datasets/:type", api.GetDatasetsByType)

	log.Println("🚀 Server running on http://localhost:8080")
	r.Run(":8080")
}

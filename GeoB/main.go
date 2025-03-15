package main

import (
	"GeoGO/api"
	"GeoGO/db"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()
	r := gin.Default()

	r.GET("/meteorites", api.GetAllMeteorites)
	r.GET("/meteorites/largest", api.GetLargestMeteorites)
	r.GET("/meteorites/nearby", api.GetNearbyMeteorites)
	log.Println("🚀 Server running on http://localhost:8080")
	r.Run(":8080")
}

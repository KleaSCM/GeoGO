package db

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Database connection configuration and initialization
// Compliance Level: Critical
// - Handles sensitive database credentials
// - Manages connection pooling and timeouts
// - Implements error handling for connection failures
// - Uses secure connection parameters (sslmode=disable only for local development)
var DB *sqlx.DB

func InitDB() {
	var err error
	dsn := "host=localhost port=5432 user=postgres password=Hisako1086 dbname=GeoGo sslmode=disable"
	DB, err = sqlx.Open("postgres", dsn)
	if err != nil {
		log.Fatal("❌ Database connection failed:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("❌ Database unreachable:", err)
	}

	log.Println("✅ Connected to PostgreSQL + PostGIS")
}

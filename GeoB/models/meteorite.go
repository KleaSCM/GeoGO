package models

// Meteorite data model definition
// Compliance Level: Critical
// - Defines data structure for database operations
// - Implements JSON serialization
// - Manages data validation
// - Handles data type conversions
type Meteorite struct {
	ID       int     `db:"id" json:"id"`
	Name     string  `db:"name" json:"name"`
	Recclass string  `db:"recclass" json:"recclass"`
	Mass     float64 `db:"mass" json:"mass"`
	Year     int     `db:"year" json:"year"`
	Lat      float64 `db:"lat" json:"lat"`
	Lon      float64 `db:"lon" json:"lon"`
	Nametype string  `db:"nametype" json:"nametype"`
	Fall     string  `db:"fall" json:"fall"`
}

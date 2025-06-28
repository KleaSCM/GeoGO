package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

// DatasetType represents different types of geospatial datasets
type DatasetType string

const (
	DatasetTypeMeteorite      DatasetType = "meteorite"
	DatasetTypeClimate        DatasetType = "climate"
	DatasetTypeWind           DatasetType = "wind"
	DatasetTypeVegetation     DatasetType = "vegetation"
	DatasetTypeInfrastructure DatasetType = "infrastructure"
	DatasetTypeFire           DatasetType = "fire"
)

// Dataset represents a unified geospatial data point
type Dataset struct {
	ID          int         `db:"id" json:"id"`
	DatasetType DatasetType `db:"dataset_type" json:"dataset_type"`
	Name        string      `db:"name" json:"name"`
	Lat         float64     `db:"lat" json:"lat"`
	Lon         float64     `db:"lon" json:"lon"`

	// Common fields
	Value     sql.NullFloat64 `db:"value" json:"-"`
	Unit      sql.NullString  `db:"unit" json:"-"`
	Timestamp *time.Time      `db:"timestamp" json:"timestamp,omitempty"`

	// Metadata fields (JSON)
	Metadata sql.NullString `db:"metadata" json:"-"`

	// Original fields for backward compatibility
	Recclass sql.NullString  `db:"recclass" json:"-"`
	Mass     sql.NullFloat64 `db:"mass" json:"-"`
	Year     sql.NullInt64   `db:"year" json:"-"`
	Nametype sql.NullString  `db:"nametype" json:"-"`
	Fall     sql.NullString  `db:"fall" json:"-"`
}

// MarshalJSON implements custom JSON marshaling for Dataset
func (d Dataset) MarshalJSON() ([]byte, error) {
	type Alias Dataset // Avoid recursive marshaling

	// Create a map for the JSON output
	output := map[string]interface{}{
		"id":           d.ID,
		"dataset_type": d.DatasetType,
		"name":         d.Name,
		"lat":          d.Lat,
		"lon":          d.Lon,
	}

	// Add nullable fields only if they are valid
	if d.Value.Valid {
		output["value"] = d.Value.Float64
	}
	if d.Unit.Valid {
		output["unit"] = d.Unit.String
	}
	if d.Timestamp != nil {
		output["timestamp"] = d.Timestamp
	}
	if d.Metadata.Valid {
		output["metadata"] = d.Metadata.String
	}
	if d.Recclass.Valid {
		output["recclass"] = d.Recclass.String
	}
	if d.Mass.Valid {
		output["mass"] = d.Mass.Float64
	}
	if d.Year.Valid {
		output["year"] = d.Year.Int64
	}
	if d.Nametype.Valid {
		output["nametype"] = d.Nametype.String
	}
	if d.Fall.Valid {
		output["fall"] = d.Fall.String
	}

	return json.Marshal(output)
}

// DatasetInfo represents metadata about available datasets
type DatasetInfo struct {
	Type        DatasetType `json:"type"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Count       int         `json:"count"`
	DateRange   string      `json:"date_range,omitempty"`
	ValueRange  string      `json:"value_range,omitempty"`
}

// ClimateData represents climate station data
type ClimateData struct {
	StationName string  `db:"station_name" json:"station_name"`
	StationID   string  `db:"station_id" json:"station_id"`
	Lat         float64 `db:"lat" json:"lat"`
	Lon         float64 `db:"lon" json:"lon"`
	Annual      float64 `db:"annual" json:"annual"`
	DJF         float64 `db:"djf" json:"djf"` // Summer
	MAM         float64 `db:"mam" json:"mam"` // Autumn
	JJA         float64 `db:"jja" json:"jja"` // Winter
	SON         float64 `db:"son" json:"son"` // Spring
	Model       string  `db:"model" json:"model"`
	RCP         string  `db:"rcp" json:"rcp"`
}

// WindData represents wind observation data
type WindData struct {
	DateTime              time.Time `db:"date_time" json:"date_time"`
	LocationDescription   string    `db:"location_description" json:"location_description"`
	Lat                   float64   `db:"lat" json:"lat"`
	Lon                   float64   `db:"lon" json:"lon"`
	AverageWindSpeed      float64   `db:"average_wind_speed" json:"average_wind_speed"`
	GustSpeed             float64   `db:"gust_speed" json:"gust_speed"`
	WindDirection         float64   `db:"wind_direction" json:"wind_direction"`
	WindDirectionCardinal string    `db:"wind_direction_cardinal" json:"wind_direction_cardinal"`
}

// VegetationData represents vegetation zone data
type VegetationData struct {
	Zone        string  `db:"zone" json:"zone"`
	Type        string  `db:"type" json:"type"`
	ShapeArea   float64 `db:"shape_area" json:"shape_area"`
	ShapeLength float64 `db:"shape_length" json:"shape_length"`
	Link        string  `db:"link" json:"link"`
}

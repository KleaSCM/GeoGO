package models

type Meteorite struct {
	ID         int     `db:"id" json:"id"`
	Name       string  `db:"name" json:"name"`
	Recclass   string  `db:"recclass" json:"recclass"`
	Mass       float64 `db:"mass" json:"mass"`
	Year       int     `db:"year" json:"year"`
	Location   string  `db:"location" json:"location"`
	DistanceKm float64 `db:"distance_km" json:"distance_km,omitempty"`
}

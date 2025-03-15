package models

type Meteorite struct {
	ID       int     `db:"id" json:"id"`
	Name     string  `db:"name" json:"name"`
	Recclass string  `db:"recclass" json:"recclass"`
	Mass     float64 `db:"mass" json:"mass"`
	Year     int     `db:"year" json:"year"`
	Location string  `db:"location" json:"location"`
	Nametype string  `db:"nametype" json:"nametype"`
	Fall     string  `db:"fall" json:"fall"`
}

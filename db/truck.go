package db

import (
	"database/sql"
)

var TruckSQL = sqlObject{
	`TABLE`, `truck`, `(
	id INTEGER PRIMARY KEY,
	last_pos coordinate NOT NULL,
	warehouse_id INTEGER -- wh heading to; NULL if idle
)`}

// CreateTruck creates a new truck.
// It should only be called when connecting to a new world.
func CreateTruck(tx *sql.Tx, id int32, lastPos Coord) error {
	sql := `INSERT INTO truck(id, last_pos) VALUES($1,$2)`
	_, err := tx.Exec(sql, id, lastPos)
	return err
}

// UpdateTruckPos update the last-known position of a truck.
func UpdateTruckPos(tx *sql.Tx, id int32, pos Coord) error {
	sql := `UPDATE truck SET last_pos = $1 WHERE id = $2`
	_, err := tx.Exec(sql, pos, id)
	return err
}

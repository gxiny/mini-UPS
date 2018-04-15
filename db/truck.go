package db

import (
	"database/sql"
)

const TruckSQL = `
CREATE TABLE truck (
	id INTEGER PRIMARY KEY,
	last_pos coordinate NOT NULL,
	warehouse_id INTEGER -- wh heading to; NULL if idle
);`

// CreateTruck creates a new truck.
// It should only be called when connecting to a new world.
func CreateTruck(tx *sql.Tx, id int64, lastPos Coord) error {
	sql := `INSERT INTO truck(id, last_pos) VALUES($1,$2)`
	_, err := tx.Exec(sql, id, lastPos)
	return err
}

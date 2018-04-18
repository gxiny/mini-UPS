package db

import (
	"database/sql"
)

var TruckSQL = sqlObject{
	`TABLE`, `truck`, `(
	id INTEGER PRIMARY KEY,
	last_pos coordinate NOT NULL,
	warehouse_id INTEGER, -- wh heading to; NULL if idle
	delivering BOOLEAN NOT NULL DEFAULT FALSE
)`}

type Truck int32

// Create creates a new truck with the given ID.
// It should only be called when connecting to a new world.
func (id Truck) Create(tx *sql.Tx, pos Coord) error {
	sql := `INSERT INTO truck(id, last_pos) VALUES($1,$2)`
	_, err := tx.Exec(sql, id, pos)
	return err
}

// UpdatePos update the last-known position of a truck.
func (id Truck) UpdatePos(tx *sql.Tx, pos Coord) error {
	sql := `UPDATE truck SET last_pos = $1 WHERE id = $2`
	_, err := tx.Exec(sql, pos, id)
	return err
}

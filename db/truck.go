package db

import (
	"database/sql"
)

// TruckStatus indicates the status of a truck.
type TruckStatus string

const (
	IDLE         TruckStatus = "idle"         // not doing anything
	TO_WAREHOUSE TruckStatus = "to_warehouse" // on the way to warehouse
	AT_WAREHOUSE TruckStatus = "at_warehouse" // staying at warehouse
	DELIVERING   TruckStatus = "delivering"   // delivering packages
)

var truckStatus = sqlObject{
	`TYPE`, `truck_status`,
	`AS ENUM ('idle', 'to_warehouse', 'at_warehouse', 'delivering')`,
}

var TruckSQL = sqlObject{
	`TABLE`, `truck`, `(
	id INTEGER PRIMARY KEY,
	last_pos coordinate NOT NULL,
	warehouse_id INTEGER, -- wh heading to; NULL if idle
	status truck_status NOT NULL DEFAULT 'idle'
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

// SendToWarehouse sets warehouse_id of the truck and
// sets its status to TO_WAREHOUSE.
func (id Truck) SendToWarehouse(tx *sql.Tx, warehouseId int32) error {
	sql := `UPDATE truck SET warehouse_id = $1, status = $2 WHERE id = $3`
	_, err := tx.Exec(sql, warehouseId, TO_WAREHOUSE, id)
	return err
}

// UpdateStatus modifies a truck's status.
func (id Truck) UpdateStatus(tx *sql.Tx, status TruckStatus) error {
	sql := `UPDATE truck SET status = $1 WHERE id = $2`
	_, err := tx.Exec(sql, status, id)
	return err
}

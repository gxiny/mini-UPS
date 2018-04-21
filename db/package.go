package db

import (
	"database/sql"
	"time"
)

var PackageTable = sqlObject{
	`TABLE`, `package`, `(
	id BIGSERIAL PRIMARY KEY,
	detail TEXT,
	destination coordinate NOT NULL,
	user_id BIGINT REFERENCES "user"(id),
	warehouse_id INTEGER NOT NULL,
	truck_id INTEGER REFERENCES truck(id),
	create_time BIGINT NOT NULL,
	deliver_time BIGINT
)`}

type Package int64

// Create creates a new package, setting its creation time to current time.
// The receiver is modified to the ID of the new package.
func (id *Package) Create(tx *sql.Tx, detail string, destination Coord, warehouseId int32) error {
	const querySQL = `INSERT INTO package(detail, destination, warehouse_id, create_time) VALUES($1,$2,$3,$4) RETURNING id`
	now := time.Now()
	return tx.QueryRow(querySQL, detail, destination, warehouseId, now.Unix()).Scan(id)
}

// SetDelivered sets delivery time of the package to current time.
func (id *Package) SetDelivered(tx *sql.Tx) (err error) {
	const querySQL = `UPDATE package SET deliver_time = $1 WHERE id = $2`
	now := time.Now()
	result, err := tx.Exec(querySQL, now.Unix(), id)
	if err != nil {
		return
	}
	n, err := result.RowsAffected()
	if err != nil {
		return
	}
	if n == 0 {
		err = sql.ErrNoRows
	}
	return
}

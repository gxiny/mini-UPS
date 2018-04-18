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

// Create creates a new package.
// The receiver is modified to the ID of the new package.
func (id *Package) Create(tx *sql.Tx, detail string, destination Coord, warehouseId int64) error {
	const sql = `INSERT INTO package(detail, destination, warehouse_id, create_time) VALUES($1,$2,$3,$4) RETURNING id`
	now := time.Now()
	return tx.QueryRow(sql, detail, destination, warehouseId, now.Unix()).Scan(id)
}

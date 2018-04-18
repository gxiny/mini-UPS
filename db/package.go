package db

import (
	"database/sql"
)

var PackageSQL = sqlObject{
	`TABLE`, `package`, `(
	id BIGSERIAL PRIMARY KEY,
	detail TEXT,
	destination coordinate NOT NULL,
	user_id BIGINT REFERENCES "user"(id),
	warehouse_id INTEGER NOT NULL,
	truck_id INTEGER REFERENCES truck(id)
)`}

type Package int64

// Create creates a new package.
// The receiver is modified to the ID of the new package.
func (id *Package) Create(tx *sql.Tx, detail string, destination Coord, warehouseId int64) error {
	sql := `INSERT INTO package(detail, destination, warehouse_id) VALUES($1,$2,$3) RETURNING id`
	return tx.QueryRow(sql, detail, destination, warehouseId).Scan(id)
}

package db

import (
	"database/sql"
)

const PackageSQL = `
CREATE TABLE package (
	id BIGSERIAL PRIMARY KEY,
	detail TEXT,
	destination coordinate NOT NULL,
	user_id BIGINT REFERENCES "user"(id),
	warehouse_id INTEGER NOT NULL,
	truck_id INTEGER REFERENCES truck(id)
);

CREATE INDEX package_idx_warehouse_id ON package(warehouse_id)
WHERE truck_id IS NULL;
`

// CreatePackage returns the ID of a newly created package.
func CreatePackage(tx *sql.Tx, detail string, destination Coord, warehouseId int64) (id int64, err error) {
	sql := `INSERT INTO package(detail, destination, warehouse_id) VALUES($1,$2,$3) RETURNING id`
	err = tx.QueryRow(sql, detail, destination, warehouseId).Scan(&id)
	return
}

package db

import (
	"database/sql"
	"database/sql/driver"
	"time"

	"github.com/golang/protobuf/proto"
)

var PackageTable = sqlObject{
	`TABLE`, `package`, `(
	id BIGSERIAL PRIMARY KEY,
	items BYTEA NOT NULL, -- serialized PackageItems
	destination coordinate NOT NULL,
	user_id BIGINT REFERENCES "user"(id),
	warehouse_id INTEGER NOT NULL,
	truck_id INTEGER REFERENCES truck(id),
	create_time BIGINT NOT NULL,
	deliver_time BIGINT
)`}

type Package int64

//go:generate protoc --go_out=. package_items.proto

// Scan implements sql.Scanner.
func (c *PackageItems) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return ErrTypeUnsupported
	}
	return proto.Unmarshal(b, c)
}

// Value implements driver.Valuer.
func (c *PackageItems) Value() (value driver.Value, err error) {
	value, err = proto.Marshal(c)
	return
}

// Create creates a new package, setting its creation time to current time.
// The receiver is modified to the ID of the new package.
func (id *Package) Create(tx *sql.Tx, items *PackageItems, destination Coord, userId sql.NullInt64, warehouseId int32) error {
	const querySQL = `INSERT INTO package(items, destination, user_id, warehouse_id, create_time)`+
		`VALUES($1,$2,$3,$4,$5) RETURNING id`
	now := time.Now()
	return tx.QueryRow(querySQL, items, destination, userId, warehouseId, now.Unix()).Scan(id)
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

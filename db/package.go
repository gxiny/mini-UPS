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
	load_time BIGINT,
	deliver_time BIGINT
)`}

var PackageView = sqlObject{
	`VIEW`, `package_view`, `AS
	SELECT package.id, package.items, package.destination,
		package.user_id,
	        package.create_time,
		package.load_time,
		package.deliver_time,
		truck.status AS truck_status
	FROM package LEFT JOIN truck
	ON package.truck_id = truck.id`,
}

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
	const querySQL = `INSERT INTO package(items, destination, user_id, warehouse_id, create_time)` +
		`VALUES($1,$2,$3,$4,$5) RETURNING id`
	now := time.Now()
	return tx.QueryRow(querySQL, items, destination, userId, warehouseId, now.Unix()).Scan(id)
}

func checkUpdate(result sql.Result, err error) error {
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err == nil && n == 0 {
		err = sql.ErrNoRows
	}
	return err
}

// SetLoaded sets loaded time of the package to current time.
func (id Package) SetLoaded(tx *sql.Tx, truck Truck) error {
	const query = `UPDATE package SET load_time = $1, truck_id = $2 WHERE id = $3`
	return checkUpdate(tx.Exec(query, time.Now().Unix(), truck, id))
}

// SetDelivered sets delivery time of the package to current time.
func (id Package) SetDelivered(tx *sql.Tx) error {
	const query = `UPDATE package SET deliver_time = $1 WHERE id = $2`
	return checkUpdate(tx.Exec(query, time.Now().Unix(), id))
}

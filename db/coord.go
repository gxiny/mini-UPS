package db

import (
	"bytes"
	"database/sql/driver"
	"errors"
	"fmt"
)

var CoordSQL = sqlObject{
	`TYPE`, `coordinate`, `AS (
	x INTEGER,
	y INTEGER
)`}

// Coord represents a coordinate in the world.
type Coord struct {
	X, Y int32
}

var ErrTypeUnsupported = errors.New("Unsupported type")

// Scan implements sql.Scanner.
func (c *Coord) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return ErrTypeUnsupported
	}
	r := bytes.NewReader(b)
	_, err := fmt.Fscanf(r, "(%d,%d)", &c.X, &c.Y)
	return err
}

// Value implements driver.Valuer.
func (c Coord) Value() (value driver.Value, err error) {
	w := bytes.NewBuffer(nil)
	_, err = fmt.Fprintf(w, "(%d,%d)", c.X, c.Y)
	if err != nil {
		return
	}
	value = w.Bytes()
	return
}

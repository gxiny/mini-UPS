package db

import (
	"testing"

	_ "github.com/lib/pq"
)

func TestCoordSQL(t *testing.T) {
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	s := Coord{1, -23}
	err = tx.QueryRow(`SELECT $1::coordinate`, s).Scan(&s)
	if err != nil {
		t.Error(err)
	}
	if (s != Coord{1, -23}) {
		t.Errorf("s = %v", s)
	}
}

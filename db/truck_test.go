package db

import (
	"testing"
)

func TestCreateTruck(t *testing.T) {
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	err = CreateTruck(tx, 1, Coord{3, 4})
	if err != nil {
		t.Error(err)
	}
}

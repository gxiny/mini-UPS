package db

import (
	"testing"
)

func TestTruckCreate(t *testing.T) {
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	truck := Truck(1)
	err = truck.Create(tx, Coord{3, 4})
	if err != nil {
		t.Error(err)
	}
}

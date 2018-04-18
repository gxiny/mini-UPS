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

func TestTruckPos(t *testing.T) {
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

	err = truck.UpdatePos(tx, Coord{5, 6})
	if err != nil {
		t.Error(err)
	}
	var pos Coord
	err = tx.QueryRow(`SELECT last_pos FROM truck WHERE id = $1`, truck).
		Scan(&pos)
	if err != nil {
		t.Error(err)
	}
	if (pos != Coord{5, 6}) {
		t.Error("last_pos != (5,6)")
	}
}

func TestTruckStatus(t *testing.T) {
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

	err = truck.SendToWarehouse(tx, 1001)
	if err != nil {
		t.Error(err)
	}
	var (
		whId   int32
		status TruckStatus
	)
	err = tx.QueryRow(`SELECT warehouse_id, status FROM truck WHERE id = $1`, truck).
		Scan(&whId, &status)
	if err != nil {
		t.Error(err)
	}
	if whId != 1001 {
		t.Error("warehouse_id != 1001")
	}
	if status != TO_WAREHOUSE {
		t.Error("status != 'to_warehouse'")
	}

	err = truck.UpdateStatus(tx, AT_WAREHOUSE)
	if err != nil {
		t.Error(err)
	}
	err = tx.QueryRow(`SELECT status FROM truck WHERE id = $1`, truck).
		Scan(&status)
	if err != nil {
		t.Error(err)
	}
	if status != AT_WAREHOUSE {
		t.Error("status != 'at_warehouse'")
	}
}

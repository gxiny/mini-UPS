package server

import (
	"database/sql"
	"errors"
	"log"

	"github.com/golang/protobuf/proto"
	"gitlab.oit.duke.edu/rz78/ups/db"
	"gitlab.oit.duke.edu/rz78/ups/pb/ups"
)

var (
	errNoIdleTruck = errors.New("no idle truck")
	errNoWarehouse = errors.New("no warehouse in need")
)

func getWarehouseInNeed(tx *sql.Tx) (warehouseId int32, err error) {
	// get a warehouse which has not-loaded-yet packages
	err = tx.QueryRow(`SELECT warehouse_id FROM package ` +
		`WHERE truck_id IS NULL ORDER BY create_time ASC ` +
		`LIMIT 1 FOR UPDATE`).Scan(&warehouseId)
	if err == sql.ErrNoRows {
		err = errNoWarehouse
	}
	return
}

func getTruckForWarehouse(tx *sql.Tx, warehouseId int32) (truck db.Truck, status db.TruckStatus, err error) {
	err = tx.QueryRow(`SELECT id, status FROM truck `+
		`WHERE status IN ($1,$2) AND warehouse_id = $3 LIMIT 1 FOR UPDATE`,
		db.TO_WAREHOUSE, db.AT_WAREHOUSE, warehouseId).Scan(&truck, &status)
	switch err {
	case nil: // found a truck
		return
	case sql.ErrNoRows:
		// fine, we'll look for idle trucks
	default: // other error
		return
	}

	// this one can be optimized, e.g., pick the nearest truck
	err = tx.QueryRow(`SELECT id FROM truck WHERE status = '` + string(db.IDLE) + `' LIMIT 1 FOR UPDATE`).
		Scan(&truck)
	switch err {
	case nil:
		status = db.IDLE
	case sql.ErrNoRows:
		err = errNoIdleTruck
	}
	return
}

func (s *Server) schedTruck(truck db.Truck) error {
	// when a truck is idle, look for where it should go
	return db.WithTx(s.db, func(tx *sql.Tx) (err error) {
		var status db.TruckStatus
		err = tx.QueryRow(`SELECT status FROM truck WHERE id=$1`, truck).Scan(&status)
		if err != nil {
			return
		}
		if status != db.IDLE { // lost a race?
			err = errNoIdleTruck
			return
		}
		warehouseId, err := getWarehouseInNeed(tx)
		if err != nil {
			return
		}
		err = s.schedPickup(tx, truck, status, warehouseId)
		return
	})
}

func (s *Server) schedWarehouse(warehouseId int32) error {
	return db.WithTx(s.db, func(tx *sql.Tx) (err error) {
		truck, status, err := getTruckForWarehouse(tx, warehouseId)
		if err != nil {
			return
		}
		err = s.schedPickup(tx, truck, status, warehouseId)
		return
	})
}

func (s *Server) schedAny() error {
	return db.WithTx(s.db, func(tx *sql.Tx) (err error) {
		warehouseId, err := getWarehouseInNeed(tx)
		if err != nil {
			return
		}
		truck, status, err := getTruckForWarehouse(tx, warehouseId)
		if err != nil {
			return
		}
		err = s.schedPickup(tx, truck, status, warehouseId)
		return
	})
}

var errNoPkgPickup = errors.New("no packages to pickup")

func (s *Server) schedPickup(tx *sql.Tx, truck db.Truck, status db.TruckStatus, warehouseId int32) (err error) {
	result, err := tx.Exec(`UPDATE package SET truck_id = $1 `+
		`WHERE warehouse_id = $2 AND truck_id IS NULL`,
		truck, warehouseId)
	if err != nil {
		return
	}
	n, err := result.RowsAffected()
	if err != nil {
		return
	}
	if n == 0 {
		err = errNoPkgPickup
		return
	}
	log.Printf("sched: %d (%s) -> %d (%d)", truck, status, warehouseId, n)
	if status == db.IDLE {
		err = truck.SendToWarehouse(tx, warehouseId)
		if err != nil {
			return
		}
		// tx succeeds; tell the world
		// what if world fails? (don't have much to do; maybe rollback)
		err = s.WriteWorld(&ups.Commands{
			Pickups: []*ups.GoPickup{
				{
					TruckId:     proto.Int32(int32(truck)),
					WarehouseId: proto.Int32(warehouseId),
				},
			},
		})
	}
	return
}

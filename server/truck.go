package server

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang/protobuf/proto"
	"gitlab.oit.duke.edu/rz78/ups/db"
	"gitlab.oit.duke.edu/rz78/ups/pb/bridge"
	"gitlab.oit.duke.edu/rz78/ups/pb/ups"
)

// initTrucks will wait until info of n trucks are received
func (s *Server) initTrucks(tx *sql.Tx, n int32) (err error) {
	for n > 0 {
		resp := new(ups.Responses)
		err = s.sim.ReadProto(resp)
		if err != nil {
			return
		}
		for _, v := range resp.GetCompletions() {
			truck := db.Truck(v.GetTruckId())
			pos := db.CoordXY(v)
			err = truck.UpdatePos(tx, pos)
			if err != nil {
				return
			}
			n--
			log.Println("created truck", truck, "at", pos)
		}
	}
	return
}

func (s *Server) onTruckFinish(truck db.Truck, pos db.Coord) (err error) {
	var (
		status      db.TruckStatus
		warehouseId int32
	)
	err = db.WithTx(s.db, func(tx *sql.Tx) (err error) {
		const sql = `SELECT status, warehouse_id FROM truck WHERE id = $1 FOR UPDATE`
		err = tx.QueryRow(sql, truck).Scan(&status, &warehouseId)
		if err != nil {
			return
		}
		switch status {
		case db.TO_WAREHOUSE:
			status = db.AT_WAREHOUSE // arrived
		case db.DELIVERING:
			status = db.IDLE // all done
		default: // can't be in other states
			panic("lost track of trucks")
		}
		err = truck.UpdateStatus(tx, status)
		if err != nil {
			return
		}
		err = truck.UpdatePos(tx, pos)
		return
	})
	if err != nil {
		return
	}
	switch status {
	case db.IDLE:
		err = s.schedTruck(truck)
	case db.AT_WAREHOUSE:
		err = s.TellAmz(&bridge.UCommands{
			Arrival: &bridge.TruckArrival{
				TruckId:     proto.Int32(int32(truck)),
				WarehouseId: &warehouseId,
			},
		})
	}
	return
}

func (s *Server) TruckReq(warehouseId int32) error {
	return s.schedWarehouse(warehouseId)
}

var errNoPkgLoaded = errors.New("no package is loaded")

func (s *Server) onTruckLoaded(loaded *bridge.PackagesLoaded) (err error) {
	truck := db.Truck(loaded.GetTruckId())
	packages := loaded.GetPackageIds()
	return db.WithTx(s.db, func(tx *sql.Tx) (err error) {
		var (
			warehouseId int32
			status      db.TruckStatus
		)
		err = tx.QueryRow(`SELECT warehouse_id, status FROM truck WHERE id = $1 FOR UPDATE`,
			truck).Scan(&warehouseId, &status)
		switch err {
		case nil:
			// fine
		case sql.ErrNoRows:
			err = fmt.Errorf("truck %d does not exist", truck)
			return
		default: // other error
			return
		}
		if status != db.AT_WAREHOUSE || warehouseId != loaded.GetWarehouseId() {
			// how can you load things if the truck is not at warehouse?
			err = fmt.Errorf("truck %d is not at warehouse %d", truck, loaded.GetWarehouseId())
			return
		}
		stmt, err := tx.Prepare(`SELECT warehouse_id, destination, load_time `+
			`FROM package WHERE id = $1 FOR UPDATE`)
		if err != nil {
			return
		}
		defer stmt.Close()
		dLocs := []*ups.DeliveryLocation{}
		for _, pkgId := range packages {
			var (
				whId int32
				dest db.Coord
				lTime sql.NullInt64
			)
			pkg := db.Package(pkgId)
			err = stmt.QueryRow(pkg).Scan(&whId, &dest, &lTime)
			if err == sql.ErrNoRows {
				err = fmt.Errorf("package %d does not exist", pkg)
			}
			if err != nil {
				return
			}
			if whId != warehouseId {
				err = fmt.Errorf("package %d is not at warehouse %d",
					pkg, warehouseId)
				return
			}
			if lTime.Valid {
				err = fmt.Errorf("package %d was loaded at %v",
					pkg, time.Unix(lTime.Int64, 0).UTC())
				return
			}
			err = pkg.SetLoaded(tx, truck)
			if err != nil {
				return
			}
			dLocs = append(dLocs, &ups.DeliveryLocation{
				PackageId: proto.Int64(pkgId),
				X:         &dest.X,
				Y:         &dest.Y,
			})
		}
		err = truck.UpdateStatus(tx, db.DELIVERING)
		if err != nil {
			return
		}
		_, err = tx.Exec(`UPDATE package SET truck_id = NULL `+
			`WHERE warehouse_id = $1 AND load_time IS NULL`,
			warehouseId)
		if err != nil {
			return
		}
		// tx succeeds; tell the world
		// don't bother with other warehouses; just go deliver
		err = s.WriteWorld(&ups.Commands{
			Deliveries: []*ups.GoDeliver{
				{
					TruckId:  proto.Int32(int32(truck)),
					Packages: dLocs,
				},
			},
		})
		return
	})
}

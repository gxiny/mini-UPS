package server

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang/protobuf/proto"
	"gitlab.oit.duke.edu/rz78/ups/db"
	"gitlab.oit.duke.edu/rz78/ups/pb"
	"gitlab.oit.duke.edu/rz78/ups/pb/bridge"
	"gitlab.oit.duke.edu/rz78/ups/pb/ups"
)

func (s *Server) createTrucks(n int32) error {
	// enclose all creations inside one transaction
	return db.WithTx(s.db, func(tx *sql.Tx) (err error) {
		r := new(ups.Responses)
		for i := int32(0); i < n; i++ {
			_, err = pb.ReadProto(s.world, r)
			if err != nil {
				return
			}
			t := r.GetCompletions()[0]
			truck := db.Truck(t.GetTruckId())
			coord := db.CoordXY(t)
			err = truck.Create(tx, coord)
			if err != nil {
				return
			}
			log.Println("Created truck", truck, "at", coord)
		}
		return
	})
}

func (s *Server) onTruckFinish(truck db.Truck, pos db.Coord) (err error) {
	return db.WithTx(s.db, func(tx *sql.Tx) (err error) {
		// there isn't a concurrent-access issue; FOR UPDATE may not be necessary
		const sql = `SELECT status FROM truck WHERE id = $1 FOR UPDATE`
		var status db.TruckStatus
		err = tx.QueryRow(sql, truck).Scan(&status)
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
}

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
		if status != db.AT_WAREHOUSE {
			// how can you load things if the truck is not at warehouse?
			err = fmt.Errorf("truck %d is not at warehouse", truck)
			return
		}
		stmt, err := tx.Prepare(`SELECT warehouse_id, destination FROM package WHERE id = $1`)
		if err != nil {
			return
		}
		defer stmt.Close()
		dLocs := []*ups.DeliveryLocation{}
		for _, pkg := range packages {
			var (
				whId int32
				dest db.Coord
			)
			err = stmt.QueryRow(pkg).Scan(&whId, &dest)
			if err == sql.ErrNoRows {
				err = fmt.Errorf("package %d does not exist", pkg)
			}
			if err != nil {
				return
			}
			if whId != warehouseId {
				err = fmt.Errorf("package %d is at warehouse %d, but truck %d is at warehouse %d",
					pkg, whId, truck, warehouseId)
				return
			}
			// BUG(rz78): The sanity check in onTruckLoaded is not enough.
			// If Amazon lies about a package, e.g., it says a package is loaded,
			// but that package actually was loaded before or is already delivered,
			// I might be tricked to make a delivery again (and get an error from world).
			dLocs = append(dLocs, &ups.DeliveryLocation{
				PackageId: proto.Int64(pkg),
				X:         &dest.X,
				Y:         &dest.Y,
			})
		}
		err = truck.UpdateStatus(tx, db.DELIVERING)
		if err != nil {
			return
		}
		// tx succeeds; tell the world
		// don't bother with other warehouses; just go deliver
		err = s.TellWorld(&ups.Commands{
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

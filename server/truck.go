package server

import (
	"database/sql"
	"log"

	"gitlab.oit.duke.edu/rz78/ups/db"
	"gitlab.oit.duke.edu/rz78/ups/pb"
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
			coord := db.Coord{t.GetX(), t.GetY()}
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

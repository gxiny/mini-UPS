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


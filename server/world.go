package server

import (
	"database/sql"
	"errors"
	"log"

	"github.com/golang/protobuf/proto"
	"gitlab.oit.duke.edu/rz78/ups/db"
	"gitlab.oit.duke.edu/rz78/ups/pb"
	"gitlab.oit.duke.edu/rz78/ups/pb/ups"
)

// NewWorld talks to the world simulator to create a new world.
func (s *Server) NewWorld(initTrucks int32) (err error) {
	_, err = s.connectWorld(&ups.Connect{
		NumTrucksInit: &initTrucks,
	})
	if err != nil {
		return
	}
	err = s.createTrucks(initTrucks)
	return
}

func (s *Server) connectWorld(c *ups.Connect) (r *ups.Connected, err error) {
	_, err = pb.WriteProto(s.world, c)
	if err != nil {
		return
	}
	r = new(ups.Connected)
	_, err = pb.ReadProto(s.world, r)
	if r.Error != nil {
		err = errors.New(*r.Error)
	}
	log.Println("Connected to world", r.GetWorldId())
	return
}

func (s *Server) createTrucks(n int32) error {
	// enclose all creations inside one transaction
	return db.WithTx(s.db, func(tx *sql.Tx) (err error) {
		r := new(ups.Responses)
		for i := int32(0); i < n; i++ {
			_, err = pb.ReadProto(s.world, r)
			if err != nil {
				return
			}
			completions := r.GetCompletions()
			truck := completions[0]
			id, coord := truck.GetTruckId(), db.Coord{truck.GetX(), truck.GetY()}
			err = db.CreateTruck(tx, id, coord)
			if err != nil {
				return
			}
			log.Println("Created truck", id, "at", coord)
		}
		return
	})
}

// ReconnectWorld reconnects to the world specified by worldId
func (s *Server) ReconnectWorld(worldId int64) (err error) {
	r, err := s.connectWorld(&ups.Connect{
		ReconnectId: &worldId,
	})
	if r.GetWorldId() != worldId {
		panic("world_id != reconnect_id")
	}
	return
}

// DisconnectWorld sends a disconnect request to the world.
// The connection is not closed.
func (s *Server) DisconnectWorld() {
	c := ups.Commands{
		Disconnect: proto.Bool(true),
	}
	_, err := pb.WriteProto(s.world, &c)
	if err != nil {
		log.Println(err)
	}
	var r ups.Responses
	_, err = pb.ReadProto(s.world, &r)
	if err != nil {
		log.Println(err)
	}
}

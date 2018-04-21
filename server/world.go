package server

import (
	"database/sql"
	"log"
	"strconv"

	"github.com/golang/protobuf/proto"
	"gitlab.oit.duke.edu/rz78/ups/db"
	"gitlab.oit.duke.edu/rz78/ups/pb/ups"
)

func (s *Server) GetWorldId() (worldId int64, err error) {
	err = db.WithTx(s.db, func(tx *sql.Tx) (err error) {
		value, err := db.GetMeta(tx, "world_id")
		if err == nil {
			worldId, err = strconv.ParseInt(value, 10, 64)
		}
		return
	})
	return
}

// NewWorld talks to the world simulator to create a new world.
func (s *Server) NewWorld(addr string, initTrucks int32) (err error) {
	worldId, err := s.sim.NewWorld(addr, initTrucks)
	if err != nil {
		return
	}
	err = db.WithTx(s.db, func(tx *sql.Tx) error {
		return db.SetMeta(tx, "world_id", strconv.FormatInt(worldId, 10))
	})
	if err != nil {
		return
	}
	log.Println("created world", worldId)
	err = s.initTrucks(initTrucks)
	return
}

// ReconnectWorld reconnects to the world specified by worldId
func (s *Server) ReconnectWorld(addr string, worldId int64) (err error) {
	err = s.sim.ReconnectWorld(addr, worldId)
	if err != nil {
		return
	}
	log.Println("reconnected to world", worldId)
	return
}

// WriteWorld writes Commands to world.
// It use a mutex to prevent concurrent writes.
func (s *Server) WriteWorld(c *ups.Commands) (err error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	err = s.sim.WriteProto(c)
	if err != nil {
		log.Println(err)
	}
	return
}

func (s *Server) acceptWorldMessages() {
	s.wg.Add(1)
	defer s.wg.Done()
	defer s.sim.Close()
	ch := s.sim.ReadChan(proto.MessageType("ups.Responses"))
	for msg := range ch {
		log.Println("world:", msg)
		if resp, ok := msg.(*ups.Responses); ok {
			s.ProcessWorldEvent(resp)
		}
	}
}

func (s *Server) ProcessWorldEvent(r *ups.Responses) {
	if completions := r.GetCompletions(); completions != nil {
		for _, finished := range completions {
			truck := db.Truck(finished.GetTruckId())
			pos := db.CoordXY(finished)
			err := s.onTruckFinish(truck, pos)
			if err != nil {
				log.Println(err)
			}
		}
	}
	if delivered := r.GetDelivered(); delivered != nil {
		for _, delivery := range delivered {
			pkg := db.Package(delivery.GetPackageId())
			err := s.onPackageDelivered(pkg)
			if err != nil {
				log.Println(err)
			}
		}
	}
	if r.Error != nil {
		log.Println(*r.Error)
	}
}

package server

import (
	"database/sql"
	"log"
	"strconv"

	"github.com/golang/protobuf/proto"
	"gitlab.oit.duke.edu/rz78/ups/db"
	"gitlab.oit.duke.edu/rz78/ups/pb/ups"
	"gitlab.oit.duke.edu/rz78/ups/world"
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
	connect := &ups.Connect{
		NumTrucksInit: &initTrucks,
	}
	connected := new(ups.Connected)
	s.sim, err = world.Connect(addr, connect, connected)
	if err != nil {
		return
	}
	worldId := connected.GetWorldId()
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
	connect := &ups.Connect{
		ReconnectId: &worldId,
	}
	connected := new(ups.Connected)
	s.sim, err = world.Connect(addr, connect, connected)
	if err != nil {
		return
	}
	if connected.GetWorldId() != worldId {
		panic("world_id != reconnect_id")
	}
	log.Println("reconnected to world", worldId)
	return
}

func (s *Server) DisconnectWorld() {
	c := &ups.Commands{
		Disconnect: proto.Bool(true),
	}
	err := s.sim.WriteProto(c)
	if err != nil {
		log.Println(err)
	}
}

func (s *Server) TellWorld(c *ups.Commands) (err error) {
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
		resp, ok := msg.(*ups.Responses)
		if ok {
			s.ProcessWorldEvent(resp)
		} else {
			log.Printf("received %T from world", msg)
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

package server

import (
	"database/sql"
	"errors"
	"log"
	"strconv"

	"github.com/golang/protobuf/proto"
	"gitlab.oit.duke.edu/rz78/ups/db"
	"gitlab.oit.duke.edu/rz78/ups/pb"
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
func (s *Server) NewWorld(initTrucks int32) (err error) {
	r, err := s.connectWorld(&ups.Connect{
		NumTrucksInit: &initTrucks,
	})
	if err != nil {
		return
	}
	worldId := r.GetWorldId()
	err = db.WithTx(s.db, func(tx *sql.Tx) error {
		return db.SetMeta(tx, "world_id", strconv.FormatInt(worldId, 10))
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
}

func (s *Server) TellWorld(c *ups.Commands) (err error) {
	_, err = pb.WriteProto(s.world, c)
	if err != nil {
		log.Println(err)
	}
	return
}

var errWorldDisconnect = errors.New("world disconnected")

func (s *Server) ListenWorld() (err error) {
	r := new(ups.Responses)
	_, err = pb.ReadProto(s.world, r)
	if err != nil {
		return
	}
	if r.GetFinished() {
		err = errWorldDisconnect
		return
	}
	s.ProcessWorldEvent(r)
	return
}

func (s *Server) ProcessWorldEvent(r *ups.Responses) {
	if completions := r.GetCompletions(); completions != nil {
		for _, finished := range completions {
			truck := db.Truck(finished.GetTruckId())
			pos := db.Coord{finished.GetX(), finished.GetY()}
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
}

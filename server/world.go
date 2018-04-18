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
func (s *Server) NewWorld(initTrucks int32) (worldId int64, err error) {
	r, err := s.connectWorld(&ups.Connect{
		NumTrucksInit: &initTrucks,
	})
	if err != nil {
		return
	}
	worldId = r.GetWorldId()
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
		log.Println(err)
		return
	}
	if r.GetFinished() {
		err = errWorldDisconnect
		return
	}
	// TODO process the message
	log.Println(r)
	return
}

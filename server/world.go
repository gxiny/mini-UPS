package server

import (
	"log"

	"github.com/golang/protobuf/proto"
	"gitlab.oit.duke.edu/rz78/ups/pb"
	"gitlab.oit.duke.edu/rz78/ups/pb/ups"
)

func (s *Server) ConnectWorld(worldId int64, initTrucks int32) (err error) {
	c := ups.Connect{
		ReconnectId: &worldId,
		NumTrucksInit: &initTrucks,
	}
	_, err = pb.WriteProto(s.world, &c)
	if err != nil {
		return
	}
	var r ups.Connected
	_, err = pb.ReadProto(s.world, &r)
	return
}

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

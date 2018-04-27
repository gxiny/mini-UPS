package server

import (
	"bufio"
	"fmt"
	"log"
	"net"

	"github.com/golang/protobuf/proto"
	"gitlab.oit.duke.edu/rz78/ups/pb"
	"gitlab.oit.duke.edu/rz78/ups/pb/bridge"
)

func (s *Server) acceptConnections() {
	s.wg.Add(1)
	defer s.wg.Done()
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			log.Println(err)
			return
		}
		go s.HandleConnection(conn)
	}
}

func (s *Server) HandleConnection(conn net.Conn) {
	s.wg.Add(1)
	defer s.wg.Done()
	defer conn.Close()

	reader := bufio.NewReader(conn)
	c := new(bridge.ACommands)
	_, err := pb.ReadProto(reader, c)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(conn.RemoteAddr(), c)
	r := s.HandleCommand(c)
	_, err = pb.WriteProto(conn, r)
	if err != nil {
		log.Println(err)
	}
}

func (s *Server) HandleCommand(c *bridge.ACommands) (resp *bridge.UResponses) {
	resp = new(bridge.UResponses)
	if req := c.GetPackageIdReq(); req != nil {
		resp.PackageId = s.PackageIdReq(req)
	} else if req := c.GetTruckReq(); req != nil {
		err := s.TruckReq(req.GetWarehouseId())
		resp.Ack = errorToAck(err)
	} else if req := c.GetLoaded(); req != nil {
		err := s.onTruckLoaded(req)
		resp.Ack = errorToAck(err)
	}
	return
}

func errorToAck(err error) *bridge.Acknowledgement {
	var ack bridge.Acknowledgement
	if err == nil {
		ack.Success = proto.Bool(true)
	} else {
		ack.Success = proto.Bool(false)
		ack.Error = proto.String(err.Error())
	}
	return &ack
}

func (s *Server) TellAmz(c *bridge.UCommands) (err error) {
	conn, err := net.Dial("tcp", s.amz)
	if err != nil {
		return
	}
	defer conn.Close()
	_, err = pb.WriteProto(conn, c)
	if err != nil {
		return
	}
	r := new(bridge.AResponses)
	_, err = pb.ReadProto(bufio.NewReader(conn), r)
	if err == nil && !r.GetAck().GetSuccess() {
		err = fmt.Errorf("amz: %s", r.GetAck().GetError())
	}
	return
}

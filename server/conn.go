package server

import (
	"bufio"
	"fmt"
	"log"
	"net"

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
	c := new(bridge.UCommands)
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

func (s *Server) HandleCommand(c *bridge.UCommands) (resp *bridge.UResponses) {
	resp = new(bridge.UResponses)
	var err error
	defer func() {
		if err != nil {
			s := err.Error()
			log.Println(s)
			resp.Error = &s
		}
	}()
	if req := c.GetPackageIdReq(); req != nil {
		resp.PackageIds, err = s.PackageIdReqs(req)
		if err != nil {
			return
		}
	}
	if req := c.GetTruckReq(); req != nil {
		err = s.TruckReq(req.GetWarehouseId()) // this one has no response
		if err != nil {
			return
		}
	}
	if req := c.GetLoaded(); req != nil {
		err = s.onTruckLoaded(req)
		if err != nil {
			return
		}
	}
	return
}

func (s *Server) TellAmz(c *bridge.ACommands) (err error) {
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
	if err == nil && r.Error != nil {
		err = fmt.Errorf("amz: %s", *r.Error)
	}
	return
}

package server

import (
	"io"
	"log"

	"gitlab.oit.duke.edu/rz78/ups/pb"
	"gitlab.oit.duke.edu/rz78/ups/pb/bridge"
)

func (s *Server) HandleConnection(conn io.ReadWriteCloser) {
	s.wg.Add(1)
	defer s.wg.Done()
	defer conn.Close()

	rw := newBufRW(conn)
	var c bridge.UCommands
	_, err := pb.ReadProto(rw, &c)
	if err != nil {
		log.Println(err)
		return
	}
	r := s.HandleCommand(&c)
	pb.WriteProto(rw, r)
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
	/* TODO(rz78)
	if req := c.GetTruckReq(); req != nil {
		err = s.TruckReq(req) // this one has no response
		if err != nil {
			return
		}
	}*/
	if req := c.GetLoaded(); req != nil {
		err = s.onTruckLoaded(req)
		if err != nil {
			return
		}
	}
	return
}

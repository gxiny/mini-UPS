package server

import (
	"bufio"
	"database/sql"
	"io"
	"log"
	"net"
	"sync"
)

type bufRW struct {
	*bufio.Reader
	io.WriteCloser
}

func newBufRW(rw io.ReadWriteCloser) *bufRW {
	return &bufRW{
		Reader:      bufio.NewReader(rw),
		WriteCloser: rw,
	}
}

type Server struct {
	db     *sql.DB
	ln     net.Listener
	wg     sync.WaitGroup
	world  *bufRW
}

// New returns a new server.
// Caller need to provide connections to the database and the world
// simulator.  These are not closed when the server is shut down.
// world is typically a TCP connection; the server will do its
// own buffering for correctly reading varint-prefixed Protocol Buffer
// objects from the stream.
func New(db *sql.DB, world io.ReadWriteCloser) *Server {
	return &Server{
		db:    db,
		world: newBufRW(world),
	}
}

// Start make the server start listening and accepting connections.
func (s *Server) Start(listenAddr string) (err error) {
	s.ln, err = net.Listen("tcp", listenAddr)
	if err != nil {
		return
	}
	log.Println("Server started")
	go s.acceptConnections()
	return
}

func (s *Server) acceptConnections() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			log.Println(err)
			return
		}
		go s.HandleConnection(conn)
	}
}

// Stop stops the server from accepting connections and waits for
// all pending connections.
func (s *Server) Stop() {
	s.ln.Close()
	s.wg.Wait()
	log.Println("Server stopped")
}

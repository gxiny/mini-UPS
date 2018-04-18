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
// Caller need to provide connections to the database.
// It is not closed when the server is shut down.
// The server will make a connection to the world.
func New(db *sql.DB, worldAddr string) (s *Server, err error) {
	worldConn, err := net.Dial("tcp", worldAddr)
	if err != nil {
		return
	}
	s = &Server{
		db:    db,
		world: newBufRW(worldConn),
	}
	return
}

// Start make the server start listening and accepting connections.
func (s *Server) Start(listenAddr string) (err error) {
	s.ln, err = net.Listen("tcp", listenAddr)
	if err != nil {
		return
	}
	log.Println("Server started")
	go s.acceptWorldMessages()
	go s.acceptConnections()
	return
}

func (s *Server) acceptWorldMessages() {
	s.wg.Add(1)
	defer s.wg.Done()
	defer s.world.Close()
	for {
		err := s.ListenWorld()
		if err == errWorldDisconnect {
			log.Println("Disconnected from world")
			return
		}
		if err != nil {
			log.Println(err)
		}
	}
}

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

// Stop stops the server from accepting connections and waits for
// all pending connections.
func (s *Server) Stop() {
	s.ln.Close()
	s.DisconnectWorld()
	s.wg.Wait()
	log.Println("Server stopped")
}

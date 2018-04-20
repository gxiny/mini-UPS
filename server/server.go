package server

import (
	"database/sql"
	"log"
	"net"
	"sync"

	"gitlab.oit.duke.edu/rz78/ups/world"
)

type Server struct {
	db  *sql.DB
	ln  net.Listener
	wg  sync.WaitGroup
	sim world.SimU
	mtx sync.Mutex
}

// New returns a new server.
// Caller need to provide connections to the database.
// It is not closed when the server is shut down.
func New(db *sql.DB) *Server {
	return &Server{
		db: db,
	}
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

// Stop stops the server from accepting connections and waits for
// all pending connections.
func (s *Server) Stop() {
	s.ln.Close()
	s.sim.DisconnectWorld()
	s.wg.Wait()
	log.Println("Server stopped")
}

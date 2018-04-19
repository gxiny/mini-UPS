// Command ups_server runs the UPS server.
// The server communicates with Amazon server and the world simulator.
package main

import (
	"database/sql"
	"flag"
	"log"
	"os"
	"os/signal"

	_ "github.com/lib/pq"
	"gitlab.oit.duke.edu/rz78/ups/server"
)

var (
	dbOptions  = flag.String("db", "dbname=ups_server user=postgres password=passw0rd", "database options")
	listenAddr = flag.String("l", ":23333", "listen address for communication with Amazon")
	worldAddr  = flag.String("sim", ":12345", "world simulator address")
	initTrucks = flag.Int("trucks", 10, "number of trucks if connecting to a new world")
)

func main() {
	flag.Parse()

	database, err := sql.Open("postgres", *dbOptions)
	if err != nil {
		log.Println(err)
		return
	}
	defer database.Close()

	s, err := server.New(database, *worldAddr)
	if err != nil {
		s.DisconnectWorld()
		log.Println(err)
		return
	}
	worldId, err := s.GetWorldId()
	if err != nil {
		err = s.NewWorld(int32(*initTrucks))
	} else {
		err = s.ReconnectWorld(worldId)
	}
	if err != nil {
		s.DisconnectWorld()
		log.Println(err)
		return
	}

	s.Start(*listenAddr)
	defer s.Stop()

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)
	<-ch

	s.Stop()
}

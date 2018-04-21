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
	listenAddr = flag.String("l", ":23333", "listen address (for receving from amz)")
	worldAddr  = flag.String("sim", ":12345", "world simulator address")
	amzAddr    = flag.String("amz", ":2333", "amz server address (for sending)")
	initTrucks = flag.Int("trucks", 10, "number of trucks if connecting to a new world")
	forceInit  = flag.Bool("init", false, "always create a new world")
)

func main() {
	flag.Parse()

	database, err := sql.Open("postgres", *dbOptions)
	if err != nil {
		log.Println(err)
		return
	}
	defer database.Close()

	s := server.New(database, *amzAddr)
	if err != nil {
		log.Println(err)
		return
	}
	worldId, err := s.GetWorldId()
	if err != nil || *forceInit {
		err = s.NewWorld(*worldAddr, int32(*initTrucks))
	} else {
		err = s.ReconnectWorld(*worldAddr, worldId)
	}
	if err != nil {
		log.Println(err)
		return
	}

	s.Start(*listenAddr)

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)
	<-ch

	s.Stop()
}

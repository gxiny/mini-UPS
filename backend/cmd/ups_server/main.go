// Command ups_server runs the UPS server.
// The server communicates with Amazon server and the world simulator.
package main

import (
	"database/sql"
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/golang/protobuf/proto"
	_ "github.com/lib/pq"
	"gitlab.oit.duke.edu/rz78/ups/pb/ups"
	"gitlab.oit.duke.edu/rz78/ups/server"
)

var (
	dbOptions  = flag.String("db", "dbname=ups_server user=postgres password=passw0rd", "database options")
	listenAddr = flag.String("l", ":23333", "listen address (for receving from amz)")
	worldAddr  = flag.String("sim", ":12345", "world simulator address")
	amzAddr    = flag.String("amz", ":2333", "amz server address (for sending)")
	initTrucks = flag.Int("trucks", 10, "number of trucks if connecting to a new world")
	forceInit  = flag.Bool("init", false, "always create a new world")
	simSpeed   = flag.Uint("speed", 0, "simulation speed (0 = default)")
)

func main() {
	flag.Parse()

	database, err := sql.Open("postgres", *dbOptions)
	if err != nil {
		log.Println(err)
		return
	}
	defer database.Close()

	waitForDB(database)

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
	if *simSpeed != 0 {
		err = s.WriteWorld(&ups.Commands{
			SimSpeed: proto.Uint32(uint32(*simSpeed)),
		})
		if err != nil {
			log.Println(err)
			return
		}
	}

	s.Start(*listenAddr)

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)
	<-ch

	s.Stop()
}

func waitForDB(database *sql.DB) {
	for {
		err := database.Ping()
		if err == nil {
			return
		}
		log.Println("wait for database:", err)
		time.Sleep(3 * time.Second)
	}
}

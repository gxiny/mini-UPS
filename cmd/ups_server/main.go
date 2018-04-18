// Command ups_server runs the UPS server.
// The server communicates with Amazon server and the world simulator.
package main

import (
	"database/sql"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"

	_ "github.com/lib/pq"
	"gitlab.oit.duke.edu/rz78/ups/db"
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

	var worldId int64
	err = db.WithTx(database, func(tx *sql.Tx) (err error) {
		value, err := db.GetMeta(tx, "world_id")
		if err == nil {
			worldId, err = strconv.ParseInt(value, 10, 64)
		}
		return
	})
	createWorld := false
	if err == sql.ErrNoRows {
		createWorld = true
	} else if err != nil {
		log.Println(err)
		return
	}

	worldConn, err := net.Dial("tcp", *worldAddr)
	if err != nil {
		log.Println(err)
		return
	}
	defer worldConn.Close()

	s := server.New(database, worldConn)
	defer s.DisconnectWorld()

	if createWorld {
		worldId, err = s.NewWorld(int32(*initTrucks))
		if err == nil {
			err = db.WithTx(database, func(tx *sql.Tx) error {
				return db.SetMeta(tx, "world_id", strconv.FormatInt(worldId, 10))
			})
		}
	} else {
		err = s.ReconnectWorld(worldId)
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

package server

import (
	"database/sql"
	"flag"
	"net"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"gitlab.oit.duke.edu/rz78/ups/db"
)

var (
	database *sql.DB
	server   *Server
)

var dbOptions = flag.String("db", "dbname=test user=postgres password=passw0rd", "database options")

const worldSimAddr = ":12345"

func TestMain(m *testing.M) {
	var err error
	database, err = sql.Open("postgres", *dbOptions)
	if err != nil {
		panic(err)
	}
	if err != nil {
		panic(err)
	}
	err = db.WithTx(database, func(tx *sql.Tx) error {
		return db.InitSchema(tx)
	})
	if err != nil {
		panic(err)
	}
	defer func() { // no matter what, destroy the schema
		err := db.WithTx(database, func(tx *sql.Tx) error {
			return db.DestroySchema(tx)
		})
		if err != nil {
			panic(err)
		}
	}()
	world, err := net.Dial("tcp", worldSimAddr)
	if err != nil {
		panic(err)
	}
	server = New(database, world)
	if err != nil {
		panic(err)
	}
	server.NewWorld(10)
	defer server.DisconnectWorld()
	m.Run()
}

func TestStartStop(t *testing.T) {
	err := server.Start(":23333") // hopefully this port is not in use
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second / 10)
	server.Stop()
}

package server

import (
	"database/sql"
	"flag"
	"testing"

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
		db.DestroySchema(tx)
		return db.InitSchema(tx)
	})
	if err != nil {
		panic(err)
	}
	server, err = New(database, worldSimAddr)
	if err != nil {
		panic(err)
	}
	_, err = server.NewWorld(10)
	if err != nil {
		panic(err)
	}
	err = server.Start(":23333")
	if err != nil {
		panic(err)
	}
	defer server.Stop()

	m.Run()
}

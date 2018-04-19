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

func initTestEnv() (err error) {
	database, err = sql.Open("postgres", *dbOptions)
	if err != nil {
		return
	}
	err = db.WithTx(database, func(tx *sql.Tx) error {
		db.DestroySchema(tx)
		return db.InitSchema(tx)
	})
	if err != nil {
		return
	}
	server = New(database)
	err = server.NewWorld(worldSimAddr, 10)
	if err != nil {
		return
	}
	err = server.Start(":23333")
	return
}

func TestMain(m *testing.M) {
	err := initTestEnv()
	if err != nil {
		panic(err)
	}
	defer database.Close()
	defer server.Stop()

	m.Run()
}

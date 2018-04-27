package server

import (
	"database/sql"
	"flag"
	"testing"

	"github.com/golang/protobuf/proto"
	_ "github.com/lib/pq"
	"gitlab.oit.duke.edu/rz78/ups/db"
	"gitlab.oit.duke.edu/rz78/ups/pb/amz"
	"gitlab.oit.duke.edu/rz78/ups/world"
)

var (
	database *sql.DB
	server   *Server
)

var dbOptions = flag.String("db", "dbname=test user=postgres password=passw0rd", "database options")

const (
	worldAddr    = ":12345"
	amzWorldAddr = ":23456"
	amzAddr      = ":2333"
)

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
	server = New(database, amzAddr)
	err = server.NewWorld(worldAddr, 10)
	if err != nil {
		return
	}
	worldId, err := server.GetWorldId()
	if err != nil {
		return
	}
	err = initWarehouses(worldId)
	if err != nil {
		return
	}
	err = server.Start(":23333")
	return
}

func initWarehouses(worldId int64) (err error) {
	connect := &amz.Connect{
		WorldId: &worldId,
		InitWarehouses: []*amz.InitWarehouse{
			{X: proto.Int32(1), Y: proto.Int32(2)},
			{X: proto.Int32(3), Y: proto.Int32(4)},
			{X: proto.Int32(5), Y: proto.Int32(6)},
		},
	}
	connected := new(amz.Connected)
	var w world.Sim
	err = w.Connect(amzWorldAddr, connect, connected)
	if err != nil {
		return
	}
	defer w.Close()
	err = w.WriteProto(&amz.Commands{
		Disconnect: proto.Bool(true),
	})
	if err != nil {
		return
	}
	err = w.ReadProto(&amz.Responses{})
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

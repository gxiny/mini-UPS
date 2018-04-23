package server

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"gitlab.oit.duke.edu/rz78/ups/db"
	"gitlab.oit.duke.edu/rz78/ups/pb/bridge"
)

func TestPackageIdReqs(t *testing.T) {
	resp := server.PackageIdReq(&bridge.Package{
		WarehouseId: proto.Int32(1),
		X:           proto.Int32(233),
		Y:           proto.Int32(666),
	})
	if e := resp.Error; e != nil {
		t.Error(*e)
	}
	var (
		whId  int64
		coord db.Coord
	)
	err := database.QueryRow(`SELECT warehouse_id, destination FROM package WHERE id = $1`, resp.GetPackageId()).Scan(&whId, &coord)
	if err != nil {
		t.Error(err)
	}
	if whId != 1 {
		t.Error("warehouse_id != 1")
	}
	if (coord != db.Coord{X: 233, Y: 666}) {
		t.Error("coord != (233,666)")
	}
}

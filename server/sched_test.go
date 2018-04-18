package server

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"gitlab.oit.duke.edu/rz78/ups/pb/bridge"
)

func TestSchedTruck(t *testing.T) {
	_, err := server.PackageIdReq(&bridge.Package{
		WarehouseId: proto.Int64(1),
		X:           proto.Int32(2333),
		Y:           proto.Int32(6666),
	})
	if err != nil {
		t.Error(err)
	}
	err = server.schedPickup()
	if err != nil {
		t.Error(err)
	}

	_, err = server.PackageIdReq(&bridge.Package{
		WarehouseId: proto.Int64(1),
		X:           proto.Int32(2),
		Y:           proto.Int32(3),
	})
	if err != nil {
		t.Error(err)
	}
	err = server.schedPickup() // should be on the same truck
	if err != nil {
		t.Error(err)
	}

	_, err = server.PackageIdReq(&bridge.Package{
		WarehouseId: proto.Int64(2),
		X:           proto.Int32(2),
		Y:           proto.Int32(3),
	})
	if err != nil {
		t.Error(err)
	}
	err = server.schedPickup() // should be on a different truck
	if err != nil {
		t.Error(err)
	}
}

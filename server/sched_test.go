package server

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"gitlab.oit.duke.edu/rz78/ups/pb/bridge"
)

func TestSchedAny(t *testing.T) {
	resp := server.PackageIdReq(&bridge.Package{
		WarehouseId: proto.Int32(1),
		X:           proto.Int32(2333),
		Y:           proto.Int32(6666),
	})
	if e := resp.Error; e != nil {
		t.Error(*e)
	}
	err := server.schedAny()
	if err != nil {
		t.Error(err)
	}

	resp = server.PackageIdReq(&bridge.Package{
		WarehouseId: proto.Int32(1),
		X:           proto.Int32(2),
		Y:           proto.Int32(3),
	})
	if e := resp.Error; e != nil {
		t.Error(*e)
	}
	err = server.schedAny()
	if err != nil {
		t.Error(err)
	}

	resp = server.PackageIdReq(&bridge.Package{
		WarehouseId: proto.Int32(2),
		X:           proto.Int32(2),
		Y:           proto.Int32(3),
	})
	if e := resp.Error; e != nil {
		t.Error(*e)
	}
	err = server.schedAny()
	if err != nil {
		t.Error(err)
	}
}

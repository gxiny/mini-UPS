package main

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"gitlab.oit.duke.edu/rz78/ups/pb/amz"
	"gitlab.oit.duke.edu/rz78/ups/pb/bridge"
)

func TestParseConnect(t *testing.T) {
	const s = "connect 123 4 5 6 7"
	ref := &amz.Connect{
		WorldId: proto.Int64(123),
		InitWarehouses: []*amz.InitWarehouse{
			{X: proto.Int32(4), Y: proto.Int32(5)},
			{X: proto.Int32(6), Y: proto.Int32(7)},
		},
	}
	msg := ParseProto(s)
	if !proto.Equal(msg, ref) {
		t.Error(msg, "!=", ref)
	}
}

func TestParsePurchase(t *testing.T) {
	const s = "purchase 1 2 \"super mario\" 5"
	ref := &amz.Commands{
		Buy: []*amz.PurchaseMore{{
			WarehouseId: proto.Int32(1),
			Things: []*amz.Product{
				{
					Id:          proto.Int64(2),
					Description: proto.String("super mario"),
					Count:       proto.Int32(5),
				},
			},
		}},
	}
	msg := ParseProto(s)
	if !proto.Equal(msg, ref) {
		t.Error(msg, "!=", ref)
	}
}

func TestParsePackage(t *testing.T) {
	const s = "pkg 1 -1 10 11 2 \"super mario\" 5"
	ref := &bridge.ACommands{
		PackageIdReq: &bridge.Package{
			WarehouseId: proto.Int32(1),
			X:           proto.Int32(10),
			Y:           proto.Int32(11),
			Items: []*bridge.Item{{
				ItemId:      proto.Int64(2),
				Description: proto.String("super mario"),
				Amount:      proto.Int32(5),
			}},
		},
	}
	msg := ParseProto(s)
	if !proto.Equal(msg, ref) {
		t.Error(msg, "!=", ref)
	}
}

func TestParseTruckReq(t *testing.T) {
	const s = "truckreq 123 10 11"
	ref := &bridge.ACommands{
		TruckReq: &bridge.RequestTruck{
			WarehouseId: proto.Int32(123),
			X:           proto.Int32(10),
			Y:           proto.Int32(11),
		},
	}
	msg := ParseProto(s)
	if !proto.Equal(msg, ref) {
		t.Error(msg, "!=", ref)
	}
}

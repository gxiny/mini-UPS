// Code generated by protoc-gen-go. DO NOT EDIT.
// source: amz.proto

/*
Package amz is a generated protocol buffer package.

Adopted from Drew's amazon.proto for better naming

It is generated from these files:
	amz.proto

It has these top-level messages:
	Product
	InitWarehouse
	Connect
	Connected
	Pack
	PutOnTruck
	PurchaseMore
	Commands
	Responses
*/
package amz

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Product struct {
	Id               *int64  `protobuf:"varint,1,req,name=id" json:"id,omitempty"`
	Description      *string `protobuf:"bytes,2,req,name=description" json:"description,omitempty"`
	Count            *int32  `protobuf:"varint,3,req,name=count" json:"count,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Product) Reset()                    { *m = Product{} }
func (m *Product) String() string            { return proto.CompactTextString(m) }
func (*Product) ProtoMessage()               {}
func (*Product) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Product) GetId() int64 {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return 0
}

func (m *Product) GetDescription() string {
	if m != nil && m.Description != nil {
		return *m.Description
	}
	return ""
}

func (m *Product) GetCount() int32 {
	if m != nil && m.Count != nil {
		return *m.Count
	}
	return 0
}

type InitWarehouse struct {
	X                *int32 `protobuf:"varint,1,req,name=x" json:"x,omitempty"`
	Y                *int32 `protobuf:"varint,2,req,name=y" json:"y,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *InitWarehouse) Reset()                    { *m = InitWarehouse{} }
func (m *InitWarehouse) String() string            { return proto.CompactTextString(m) }
func (*InitWarehouse) ProtoMessage()               {}
func (*InitWarehouse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *InitWarehouse) GetX() int32 {
	if m != nil && m.X != nil {
		return *m.X
	}
	return 0
}

func (m *InitWarehouse) GetY() int32 {
	if m != nil && m.Y != nil {
		return *m.Y
	}
	return 0
}

type Connect struct {
	WorldId          *int64           `protobuf:"varint,1,req,name=world_id" json:"world_id,omitempty"`
	InitWarehouses   []*InitWarehouse `protobuf:"bytes,2,rep,name=init_warehouses" json:"init_warehouses,omitempty"`
	XXX_unrecognized []byte           `json:"-"`
}

func (m *Connect) Reset()                    { *m = Connect{} }
func (m *Connect) String() string            { return proto.CompactTextString(m) }
func (*Connect) ProtoMessage()               {}
func (*Connect) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Connect) GetWorldId() int64 {
	if m != nil && m.WorldId != nil {
		return *m.WorldId
	}
	return 0
}

func (m *Connect) GetInitWarehouses() []*InitWarehouse {
	if m != nil {
		return m.InitWarehouses
	}
	return nil
}

type Connected struct {
	Error            *string `protobuf:"bytes,1,opt,name=error" json:"error,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Connected) Reset()                    { *m = Connected{} }
func (m *Connected) String() string            { return proto.CompactTextString(m) }
func (*Connected) ProtoMessage()               {}
func (*Connected) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *Connected) GetError() string {
	if m != nil && m.Error != nil {
		return *m.Error
	}
	return ""
}

type Pack struct {
	WarehouseId      *int32     `protobuf:"varint,1,req,name=warehouse_id" json:"warehouse_id,omitempty"`
	Things           []*Product `protobuf:"bytes,2,rep,name=things" json:"things,omitempty"`
	ShipId           *int64     `protobuf:"varint,3,req,name=ship_id" json:"ship_id,omitempty"`
	XXX_unrecognized []byte     `json:"-"`
}

func (m *Pack) Reset()                    { *m = Pack{} }
func (m *Pack) String() string            { return proto.CompactTextString(m) }
func (*Pack) ProtoMessage()               {}
func (*Pack) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *Pack) GetWarehouseId() int32 {
	if m != nil && m.WarehouseId != nil {
		return *m.WarehouseId
	}
	return 0
}

func (m *Pack) GetThings() []*Product {
	if m != nil {
		return m.Things
	}
	return nil
}

func (m *Pack) GetShipId() int64 {
	if m != nil && m.ShipId != nil {
		return *m.ShipId
	}
	return 0
}

type PutOnTruck struct {
	WarehouseId      *int32 `protobuf:"varint,1,req,name=warehouse_id" json:"warehouse_id,omitempty"`
	TruckId          *int32 `protobuf:"varint,2,req,name=truck_id" json:"truck_id,omitempty"`
	ShipId           *int64 `protobuf:"varint,3,req,name=ship_id" json:"ship_id,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *PutOnTruck) Reset()                    { *m = PutOnTruck{} }
func (m *PutOnTruck) String() string            { return proto.CompactTextString(m) }
func (*PutOnTruck) ProtoMessage()               {}
func (*PutOnTruck) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *PutOnTruck) GetWarehouseId() int32 {
	if m != nil && m.WarehouseId != nil {
		return *m.WarehouseId
	}
	return 0
}

func (m *PutOnTruck) GetTruckId() int32 {
	if m != nil && m.TruckId != nil {
		return *m.TruckId
	}
	return 0
}

func (m *PutOnTruck) GetShipId() int64 {
	if m != nil && m.ShipId != nil {
		return *m.ShipId
	}
	return 0
}

type PurchaseMore struct {
	WarehouseId      *int32     `protobuf:"varint,1,req,name=warehouse_id" json:"warehouse_id,omitempty"`
	Things           []*Product `protobuf:"bytes,2,rep,name=things" json:"things,omitempty"`
	XXX_unrecognized []byte     `json:"-"`
}

func (m *PurchaseMore) Reset()                    { *m = PurchaseMore{} }
func (m *PurchaseMore) String() string            { return proto.CompactTextString(m) }
func (*PurchaseMore) ProtoMessage()               {}
func (*PurchaseMore) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *PurchaseMore) GetWarehouseId() int32 {
	if m != nil && m.WarehouseId != nil {
		return *m.WarehouseId
	}
	return 0
}

func (m *PurchaseMore) GetThings() []*Product {
	if m != nil {
		return m.Things
	}
	return nil
}

type Commands struct {
	Buy              []*PurchaseMore `protobuf:"bytes,1,rep,name=buy" json:"buy,omitempty"`
	Load             []*PutOnTruck   `protobuf:"bytes,2,rep,name=load" json:"load,omitempty"`
	ToPack           []*Pack         `protobuf:"bytes,3,rep,name=to_pack" json:"to_pack,omitempty"`
	SimSpeed         *uint32         `protobuf:"varint,4,opt,name=sim_speed" json:"sim_speed,omitempty"`
	Disconnect       *bool           `protobuf:"varint,5,opt,name=disconnect" json:"disconnect,omitempty"`
	XXX_unrecognized []byte          `json:"-"`
}

func (m *Commands) Reset()                    { *m = Commands{} }
func (m *Commands) String() string            { return proto.CompactTextString(m) }
func (*Commands) ProtoMessage()               {}
func (*Commands) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *Commands) GetBuy() []*PurchaseMore {
	if m != nil {
		return m.Buy
	}
	return nil
}

func (m *Commands) GetLoad() []*PutOnTruck {
	if m != nil {
		return m.Load
	}
	return nil
}

func (m *Commands) GetToPack() []*Pack {
	if m != nil {
		return m.ToPack
	}
	return nil
}

func (m *Commands) GetSimSpeed() uint32 {
	if m != nil && m.SimSpeed != nil {
		return *m.SimSpeed
	}
	return 0
}

func (m *Commands) GetDisconnect() bool {
	if m != nil && m.Disconnect != nil {
		return *m.Disconnect
	}
	return false
}

type Responses struct {
	Arrived          []*PurchaseMore `protobuf:"bytes,1,rep,name=arrived" json:"arrived,omitempty"`
	Ready            []int64         `protobuf:"varint,2,rep,name=ready" json:"ready,omitempty"`
	Loaded           []int64         `protobuf:"varint,3,rep,name=loaded" json:"loaded,omitempty"`
	Error            *string         `protobuf:"bytes,4,opt,name=error" json:"error,omitempty"`
	Finished         *bool           `protobuf:"varint,5,opt,name=finished" json:"finished,omitempty"`
	XXX_unrecognized []byte          `json:"-"`
}

func (m *Responses) Reset()                    { *m = Responses{} }
func (m *Responses) String() string            { return proto.CompactTextString(m) }
func (*Responses) ProtoMessage()               {}
func (*Responses) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *Responses) GetArrived() []*PurchaseMore {
	if m != nil {
		return m.Arrived
	}
	return nil
}

func (m *Responses) GetReady() []int64 {
	if m != nil {
		return m.Ready
	}
	return nil
}

func (m *Responses) GetLoaded() []int64 {
	if m != nil {
		return m.Loaded
	}
	return nil
}

func (m *Responses) GetError() string {
	if m != nil && m.Error != nil {
		return *m.Error
	}
	return ""
}

func (m *Responses) GetFinished() bool {
	if m != nil && m.Finished != nil {
		return *m.Finished
	}
	return false
}

func init() {
	proto.RegisterType((*Product)(nil), "amz.Product")
	proto.RegisterType((*InitWarehouse)(nil), "amz.InitWarehouse")
	proto.RegisterType((*Connect)(nil), "amz.Connect")
	proto.RegisterType((*Connected)(nil), "amz.Connected")
	proto.RegisterType((*Pack)(nil), "amz.Pack")
	proto.RegisterType((*PutOnTruck)(nil), "amz.PutOnTruck")
	proto.RegisterType((*PurchaseMore)(nil), "amz.PurchaseMore")
	proto.RegisterType((*Commands)(nil), "amz.Commands")
	proto.RegisterType((*Responses)(nil), "amz.Responses")
}

func init() { proto.RegisterFile("amz.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 404 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x91, 0x5f, 0x6b, 0x14, 0x31,
	0x14, 0xc5, 0xc9, 0x64, 0xa7, 0xb3, 0xb9, 0xdd, 0x75, 0xdb, 0xe8, 0x43, 0x28, 0x2a, 0x43, 0x40,
	0x18, 0x10, 0xfa, 0xe0, 0x9b, 0xaf, 0x16, 0x41, 0x11, 0x71, 0x11, 0xc1, 0xc7, 0x25, 0x26, 0x57,
	0x27, 0xb4, 0x93, 0x0c, 0x49, 0xc6, 0xba, 0xfd, 0x0c, 0x7e, 0x68, 0x49, 0xf6, 0x8f, 0x8a, 0xf4,
	0xc5, 0xb7, 0x90, 0x9c, 0xfb, 0x3b, 0xe7, 0xdc, 0x00, 0x53, 0xc3, 0xdd, 0xe5, 0x18, 0x7c, 0xf2,
	0x9c, 0xaa, 0xe1, 0x4e, 0xbe, 0x84, 0x66, 0x1d, 0xbc, 0x99, 0x74, 0xe2, 0x00, 0x95, 0x35, 0x82,
	0xb4, 0x55, 0x47, 0xf9, 0x43, 0x38, 0x35, 0x18, 0x75, 0xb0, 0x63, 0xb2, 0xde, 0x89, 0xaa, 0xad,
	0x3a, 0xc6, 0x97, 0x50, 0x6b, 0x3f, 0xb9, 0x24, 0x68, 0x5b, 0x75, 0xb5, 0x7c, 0x06, 0xcb, 0xb7,
	0xce, 0xa6, 0xcf, 0x2a, 0x60, 0xef, 0xa7, 0x88, 0x9c, 0x01, 0xf9, 0x51, 0xe6, 0xeb, 0x7c, 0xdc,
	0x96, 0xa9, 0x5a, 0xbe, 0x81, 0xe6, 0xca, 0x3b, 0x87, 0x3a, 0xf1, 0x33, 0x98, 0xdf, 0xfa, 0x70,
	0x63, 0x36, 0x47, 0x9f, 0xe7, 0xb0, 0xb2, 0xce, 0xa6, 0xcd, 0xed, 0x01, 0x12, 0x45, 0xd5, 0xd2,
	0xee, 0xf4, 0x05, 0xbf, 0xcc, 0x41, 0xff, 0xe2, 0xcb, 0x0b, 0x60, 0x7b, 0x12, 0x9a, 0x1c, 0x06,
	0x43, 0xf0, 0x41, 0x90, 0x96, 0x74, 0x4c, 0xbe, 0x83, 0xd9, 0x5a, 0xe9, 0x6b, 0xfe, 0x08, 0x16,
	0x47, 0xd6, 0xc1, 0xa6, 0xe6, 0x8f, 0xe1, 0x24, 0xf5, 0xd6, 0x7d, 0x3b, 0xd0, 0x17, 0x85, 0x7e,
	0x28, 0xbe, 0x82, 0x26, 0xf6, 0x76, 0xcc, 0xf2, 0xdc, 0x8c, 0xca, 0xd7, 0x00, 0xeb, 0x29, 0x7d,
	0x70, 0x9f, 0xc2, 0x74, 0x2f, 0xf2, 0x0c, 0xe6, 0x29, 0x3f, 0xe7, 0x9b, 0x52, 0xf4, 0x5f, 0xcc,
	0x2b, 0x58, 0xac, 0xa7, 0xa0, 0x7b, 0x15, 0xf1, 0xbd, 0x0f, 0xf8, 0x3f, 0xd9, 0xe4, 0x4f, 0x02,
	0xf3, 0x2b, 0x3f, 0x0c, 0xca, 0x99, 0xc8, 0x9f, 0x02, 0xfd, 0x32, 0x6d, 0x05, 0x29, 0xba, 0xf3,
	0x9d, 0xee, 0x4f, 0x83, 0x27, 0x30, 0xbb, 0xf1, 0xca, 0xec, 0x41, 0xab, 0xbd, 0xe0, 0x58, 0xe4,
	0x02, 0x9a, 0xe4, 0x37, 0xa3, 0xd2, 0xd7, 0x82, 0x16, 0x05, 0xdb, 0x29, 0xf2, 0xde, 0xce, 0x81,
	0x45, 0x3b, 0x6c, 0xe2, 0x88, 0x68, 0xc4, 0xac, 0x25, 0xdd, 0x92, 0x73, 0x00, 0x63, 0xa3, 0xde,
	0x6d, 0x5c, 0xd4, 0x2d, 0xe9, 0xe6, 0xd2, 0x03, 0xfb, 0x88, 0x71, 0xf4, 0x2e, 0x62, 0xe4, 0x12,
	0x1a, 0x15, 0x82, 0xfd, 0x8e, 0xe6, 0xfe, 0x48, 0x4b, 0xa8, 0x03, 0x2a, 0xb3, 0x2d, 0x99, 0x28,
	0x7f, 0x00, 0x27, 0x39, 0x21, 0x9a, 0x92, 0x80, 0xfe, 0xfe, 0xc5, 0x6c, 0xc9, 0xf2, 0x52, 0xbf,
	0x5a, 0x67, 0x63, 0x8f, 0x66, 0x67, 0xf8, 0x2b, 0x00, 0x00, 0xff, 0xff, 0xb8, 0xe1, 0xdb, 0x2e,
	0xb0, 0x02, 0x00, 0x00,
}

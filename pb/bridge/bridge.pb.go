// Code generated by protoc-gen-go. DO NOT EDIT.
// source: bridge.proto

/*
Package bridge is a generated protocol buffer package.

It is generated from these files:
	bridge.proto

It has these top-level messages:
	Item
	Package
	UCommands
	UResponses
	ACommands
	AResponses
	ResponsePackageId
	RequestTruck
	TruckArrival
	PackagesLoaded
*/
package bridge

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

type Item struct {
	ItemId           *int64  `protobuf:"varint,1,req,name=item_id" json:"item_id,omitempty"`
	Description      *string `protobuf:"bytes,2,req,name=description" json:"description,omitempty"`
	Amount           *int32  `protobuf:"varint,3,req,name=amount" json:"amount,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Item) Reset()                    { *m = Item{} }
func (m *Item) String() string            { return proto.CompactTextString(m) }
func (*Item) ProtoMessage()               {}
func (*Item) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Item) GetItemId() int64 {
	if m != nil && m.ItemId != nil {
		return *m.ItemId
	}
	return 0
}

func (m *Item) GetDescription() string {
	if m != nil && m.Description != nil {
		return *m.Description
	}
	return ""
}

func (m *Item) GetAmount() int32 {
	if m != nil && m.Amount != nil {
		return *m.Amount
	}
	return 0
}

type Package struct {
	WarehouseId      *int32  `protobuf:"varint,1,req,name=warehouse_id" json:"warehouse_id,omitempty"`
	UpsUserId        *int64  `protobuf:"varint,2,opt,name=ups_user_id" json:"ups_user_id,omitempty"`
	X                *int32  `protobuf:"varint,3,req,name=x" json:"x,omitempty"`
	Y                *int32  `protobuf:"varint,4,req,name=y" json:"y,omitempty"`
	Items            []*Item `protobuf:"bytes,5,rep,name=items" json:"items,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Package) Reset()                    { *m = Package{} }
func (m *Package) String() string            { return proto.CompactTextString(m) }
func (*Package) ProtoMessage()               {}
func (*Package) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Package) GetWarehouseId() int32 {
	if m != nil && m.WarehouseId != nil {
		return *m.WarehouseId
	}
	return 0
}

func (m *Package) GetUpsUserId() int64 {
	if m != nil && m.UpsUserId != nil {
		return *m.UpsUserId
	}
	return 0
}

func (m *Package) GetX() int32 {
	if m != nil && m.X != nil {
		return *m.X
	}
	return 0
}

func (m *Package) GetY() int32 {
	if m != nil && m.Y != nil {
		return *m.Y
	}
	return 0
}

func (m *Package) GetItems() []*Item {
	if m != nil {
		return m.Items
	}
	return nil
}

type UCommands struct {
	PackageIdReq     []*Package      `protobuf:"bytes,1,rep,name=package_id_req" json:"package_id_req,omitempty"`
	TruckReq         *RequestTruck   `protobuf:"bytes,2,opt,name=truck_req" json:"truck_req,omitempty"`
	Loaded           *PackagesLoaded `protobuf:"bytes,3,opt,name=loaded" json:"loaded,omitempty"`
	XXX_unrecognized []byte          `json:"-"`
}

func (m *UCommands) Reset()                    { *m = UCommands{} }
func (m *UCommands) String() string            { return proto.CompactTextString(m) }
func (*UCommands) ProtoMessage()               {}
func (*UCommands) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *UCommands) GetPackageIdReq() []*Package {
	if m != nil {
		return m.PackageIdReq
	}
	return nil
}

func (m *UCommands) GetTruckReq() *RequestTruck {
	if m != nil {
		return m.TruckReq
	}
	return nil
}

func (m *UCommands) GetLoaded() *PackagesLoaded {
	if m != nil {
		return m.Loaded
	}
	return nil
}

type UResponses struct {
	Error            *string              `protobuf:"bytes,1,req,name=error" json:"error,omitempty"`
	PackageIds       []*ResponsePackageId `protobuf:"bytes,2,rep,name=package_ids" json:"package_ids,omitempty"`
	XXX_unrecognized []byte               `json:"-"`
}

func (m *UResponses) Reset()                    { *m = UResponses{} }
func (m *UResponses) String() string            { return proto.CompactTextString(m) }
func (*UResponses) ProtoMessage()               {}
func (*UResponses) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *UResponses) GetError() string {
	if m != nil && m.Error != nil {
		return *m.Error
	}
	return ""
}

func (m *UResponses) GetPackageIds() []*ResponsePackageId {
	if m != nil {
		return m.PackageIds
	}
	return nil
}

type ACommands struct {
	Arrival          *TruckArrival `protobuf:"bytes,1,opt,name=arrival" json:"arrival,omitempty"`
	XXX_unrecognized []byte        `json:"-"`
}

func (m *ACommands) Reset()                    { *m = ACommands{} }
func (m *ACommands) String() string            { return proto.CompactTextString(m) }
func (*ACommands) ProtoMessage()               {}
func (*ACommands) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *ACommands) GetArrival() *TruckArrival {
	if m != nil {
		return m.Arrival
	}
	return nil
}

type AResponses struct {
	Error            *string `protobuf:"bytes,1,req,name=error" json:"error,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *AResponses) Reset()                    { *m = AResponses{} }
func (m *AResponses) String() string            { return proto.CompactTextString(m) }
func (*AResponses) ProtoMessage()               {}
func (*AResponses) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *AResponses) GetError() string {
	if m != nil && m.Error != nil {
		return *m.Error
	}
	return ""
}

type ResponsePackageId struct {
	PackageId        *int64  `protobuf:"varint,1,req,name=package_id" json:"package_id,omitempty"`
	Error            *string `protobuf:"bytes,2,opt,name=error" json:"error,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *ResponsePackageId) Reset()                    { *m = ResponsePackageId{} }
func (m *ResponsePackageId) String() string            { return proto.CompactTextString(m) }
func (*ResponsePackageId) ProtoMessage()               {}
func (*ResponsePackageId) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *ResponsePackageId) GetPackageId() int64 {
	if m != nil && m.PackageId != nil {
		return *m.PackageId
	}
	return 0
}

func (m *ResponsePackageId) GetError() string {
	if m != nil && m.Error != nil {
		return *m.Error
	}
	return ""
}

type RequestTruck struct {
	WarehouseId      *int32 `protobuf:"varint,1,req,name=warehouse_id" json:"warehouse_id,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *RequestTruck) Reset()                    { *m = RequestTruck{} }
func (m *RequestTruck) String() string            { return proto.CompactTextString(m) }
func (*RequestTruck) ProtoMessage()               {}
func (*RequestTruck) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *RequestTruck) GetWarehouseId() int32 {
	if m != nil && m.WarehouseId != nil {
		return *m.WarehouseId
	}
	return 0
}

type TruckArrival struct {
	TruckId          *int32 `protobuf:"varint,1,req,name=truck_id" json:"truck_id,omitempty"`
	WarehouseId      *int32 `protobuf:"varint,2,req,name=warehouse_id" json:"warehouse_id,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *TruckArrival) Reset()                    { *m = TruckArrival{} }
func (m *TruckArrival) String() string            { return proto.CompactTextString(m) }
func (*TruckArrival) ProtoMessage()               {}
func (*TruckArrival) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *TruckArrival) GetTruckId() int32 {
	if m != nil && m.TruckId != nil {
		return *m.TruckId
	}
	return 0
}

func (m *TruckArrival) GetWarehouseId() int32 {
	if m != nil && m.WarehouseId != nil {
		return *m.WarehouseId
	}
	return 0
}

type PackagesLoaded struct {
	PackageIds       []int64 `protobuf:"varint,1,rep,name=package_ids" json:"package_ids,omitempty"`
	TruckId          *int32  `protobuf:"varint,2,req,name=truck_id" json:"truck_id,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *PackagesLoaded) Reset()                    { *m = PackagesLoaded{} }
func (m *PackagesLoaded) String() string            { return proto.CompactTextString(m) }
func (*PackagesLoaded) ProtoMessage()               {}
func (*PackagesLoaded) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func (m *PackagesLoaded) GetPackageIds() []int64 {
	if m != nil {
		return m.PackageIds
	}
	return nil
}

func (m *PackagesLoaded) GetTruckId() int32 {
	if m != nil && m.TruckId != nil {
		return *m.TruckId
	}
	return 0
}

func init() {
	proto.RegisterType((*Item)(nil), "bridge.Item")
	proto.RegisterType((*Package)(nil), "bridge.Package")
	proto.RegisterType((*UCommands)(nil), "bridge.UCommands")
	proto.RegisterType((*UResponses)(nil), "bridge.UResponses")
	proto.RegisterType((*ACommands)(nil), "bridge.ACommands")
	proto.RegisterType((*AResponses)(nil), "bridge.AResponses")
	proto.RegisterType((*ResponsePackageId)(nil), "bridge.ResponsePackageId")
	proto.RegisterType((*RequestTruck)(nil), "bridge.RequestTruck")
	proto.RegisterType((*TruckArrival)(nil), "bridge.TruckArrival")
	proto.RegisterType((*PackagesLoaded)(nil), "bridge.PackagesLoaded")
}

func init() { proto.RegisterFile("bridge.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 383 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x92, 0x3f, 0x6f, 0xea, 0x30,
	0x14, 0xc5, 0x95, 0x84, 0xc0, 0xcb, 0x4d, 0x1e, 0xbc, 0xe7, 0xa2, 0xca, 0x15, 0x4b, 0x64, 0xb5,
	0x25, 0x13, 0x03, 0x03, 0x5d, 0xba, 0xa0, 0x4e, 0xa8, 0x1d, 0x2a, 0x54, 0xe6, 0xc8, 0xc5, 0x16,
	0x8d, 0x20, 0x71, 0xb0, 0x93, 0xfe, 0xf9, 0x00, 0xfd, 0xde, 0x95, 0x1d, 0x02, 0x84, 0xb6, 0x9b,
	0x93, 0x7b, 0xcf, 0xb9, 0x3f, 0x9f, 0x6b, 0x08, 0x9e, 0x65, 0xc2, 0x56, 0x7c, 0x94, 0x4b, 0x51,
	0x08, 0xd4, 0xae, 0xbe, 0xc8, 0x2d, 0xb4, 0x66, 0x05, 0x4f, 0x51, 0x0f, 0x3a, 0x49, 0xc1, 0xd3,
	0x38, 0x61, 0xd8, 0x0a, 0xed, 0xc8, 0x41, 0x67, 0xe0, 0x33, 0xae, 0x96, 0x32, 0xc9, 0x8b, 0x44,
	0x64, 0xd8, 0x0e, 0xed, 0xc8, 0x43, 0x5d, 0x68, 0xd3, 0x54, 0x94, 0x59, 0x81, 0x9d, 0xd0, 0x8e,
	0x5c, 0xb2, 0x82, 0xce, 0x23, 0x5d, 0xae, 0xe9, 0x8a, 0xa3, 0x3e, 0x04, 0x6f, 0x54, 0xf2, 0x17,
	0x51, 0x2a, 0x5e, 0xbb, 0xb8, 0xda, 0xa5, 0xcc, 0x55, 0x5c, 0x2a, 0x2e, 0xf5, 0x4f, 0x3b, 0xb4,
	0x22, 0x07, 0x79, 0x60, 0xbd, 0x57, 0x06, 0xfa, 0xf8, 0x81, 0x5b, 0xe6, 0x38, 0x00, 0x57, 0x13,
	0x28, 0xec, 0x86, 0x4e, 0xe4, 0x8f, 0x83, 0xd1, 0x8e, 0x57, 0xe3, 0x91, 0x4f, 0x0b, 0xbc, 0xc5,
	0x9d, 0x48, 0x53, 0x9a, 0x31, 0x85, 0x86, 0xd0, 0xcd, 0xab, 0xb1, 0x71, 0xc2, 0x62, 0xc9, 0xb7,
	0xd8, 0x32, 0x9a, 0x5e, 0xad, 0xa9, 0xa1, 0x86, 0xe0, 0x15, 0xb2, 0x5c, 0xae, 0x4d, 0x8f, 0x1e,
	0xee, 0x8f, 0xfb, 0x75, 0xcf, 0x9c, 0x6f, 0x4b, 0xae, 0x8a, 0x27, 0x5d, 0x47, 0xd7, 0xd0, 0xde,
	0x08, 0xca, 0x38, 0xc3, 0x8e, 0xe9, 0x3a, 0x3f, 0x71, 0x52, 0x0f, 0xa6, 0x4a, 0xee, 0x01, 0x16,
	0x73, 0xae, 0x72, 0x91, 0x29, 0xae, 0xd0, 0x5f, 0x70, 0xb9, 0x94, 0x42, 0x9a, 0xcb, 0x7a, 0x68,
	0x04, 0xfe, 0x01, 0x4b, 0x61, 0xdb, 0x30, 0x5d, 0x1c, 0xe6, 0x55, 0xb2, 0x9d, 0xe3, 0x8c, 0x91,
	0x31, 0x78, 0xd3, 0xfd, 0x9d, 0xae, 0xa0, 0x43, 0xa5, 0x4c, 0x5e, 0xe9, 0x06, 0x5b, 0x4d, 0x50,
	0x43, 0x38, 0xad, 0x6a, 0x64, 0x00, 0x30, 0xfd, 0x0d, 0x80, 0x4c, 0xe0, 0xff, 0xb7, 0x29, 0x08,
	0x01, 0x1c, 0xa8, 0x76, 0xcb, 0xdd, 0xeb, 0x74, 0x26, 0x1e, 0xb9, 0x84, 0xa0, 0x91, 0xc6, 0x8f,
	0xbb, 0x24, 0x13, 0x08, 0x8e, 0x51, 0xd0, 0x3f, 0xf8, 0x53, 0x85, 0xbb, 0xdf, 0xf6, 0xa9, 0xce,
	0x36, 0xba, 0x1b, 0xe8, 0x36, 0x53, 0xd4, 0xaf, 0xe2, 0x38, 0x28, 0xbd, 0x3c, 0xa7, 0x61, 0x67,
	0x84, 0x5f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x30, 0xb4, 0xe3, 0x95, 0xb2, 0x02, 0x00, 0x00,
}

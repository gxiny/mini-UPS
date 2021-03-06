// Code generated by protoc-gen-go. DO NOT EDIT.
// source: web.proto

/*
Package web is a generated protocol buffer package.

It is generated from these files:
	web.proto

It has these top-level messages:
	PkgDest
	PkgListReq
	Request
	Response
	PkgList
	PkgDetail
*/
package web

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

type PkgDest struct {
	PackageId        *int64 `protobuf:"varint,1,req,name=package_id" json:"package_id,omitempty"`
	X                *int32 `protobuf:"varint,2,req,name=x" json:"x,omitempty"`
	Y                *int32 `protobuf:"varint,3,req,name=y" json:"y,omitempty"`
	UserId           *int64 `protobuf:"varint,4,req,name=user_id" json:"user_id,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *PkgDest) Reset()                    { *m = PkgDest{} }
func (m *PkgDest) String() string            { return proto.CompactTextString(m) }
func (*PkgDest) ProtoMessage()               {}
func (*PkgDest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *PkgDest) GetPackageId() int64 {
	if m != nil && m.PackageId != nil {
		return *m.PackageId
	}
	return 0
}

func (m *PkgDest) GetX() int32 {
	if m != nil && m.X != nil {
		return *m.X
	}
	return 0
}

func (m *PkgDest) GetY() int32 {
	if m != nil && m.Y != nil {
		return *m.Y
	}
	return 0
}

func (m *PkgDest) GetUserId() int64 {
	if m != nil && m.UserId != nil {
		return *m.UserId
	}
	return 0
}

type PkgListReq struct {
	UserId     *int64  `protobuf:"varint,1,opt,name=user_id" json:"user_id,omitempty"`
	PackageIds []int64 `protobuf:"varint,2,rep,name=package_ids" json:"package_ids,omitempty"`
	// if neither: get all shipments in the system
	Offset           *int64 `protobuf:"varint,3,opt,name=offset" json:"offset,omitempty"`
	Limit            *int64 `protobuf:"varint,4,opt,name=limit" json:"limit,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *PkgListReq) Reset()                    { *m = PkgListReq{} }
func (m *PkgListReq) String() string            { return proto.CompactTextString(m) }
func (*PkgListReq) ProtoMessage()               {}
func (*PkgListReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *PkgListReq) GetUserId() int64 {
	if m != nil && m.UserId != nil {
		return *m.UserId
	}
	return 0
}

func (m *PkgListReq) GetPackageIds() []int64 {
	if m != nil {
		return m.PackageIds
	}
	return nil
}

func (m *PkgListReq) GetOffset() int64 {
	if m != nil && m.Offset != nil {
		return *m.Offset
	}
	return 0
}

func (m *PkgListReq) GetLimit() int64 {
	if m != nil && m.Limit != nil {
		return *m.Limit
	}
	return 0
}

type Request struct {
	NewUser           *string     `protobuf:"bytes,1,opt,name=new_user" json:"new_user,omitempty"`
	GetPackageList    *PkgListReq `protobuf:"bytes,2,opt,name=get_package_list" json:"get_package_list,omitempty"`
	GetPackageDetail  *int64      `protobuf:"varint,3,opt,name=get_package_detail" json:"get_package_detail,omitempty"`
	ChangeDestination *PkgDest    `protobuf:"bytes,4,opt,name=change_destination" json:"change_destination,omitempty"`
	XXX_unrecognized  []byte      `json:"-"`
}

func (m *Request) Reset()                    { *m = Request{} }
func (m *Request) String() string            { return proto.CompactTextString(m) }
func (*Request) ProtoMessage()               {}
func (*Request) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Request) GetNewUser() string {
	if m != nil && m.NewUser != nil {
		return *m.NewUser
	}
	return ""
}

func (m *Request) GetGetPackageList() *PkgListReq {
	if m != nil {
		return m.GetPackageList
	}
	return nil
}

func (m *Request) GetGetPackageDetail() int64 {
	if m != nil && m.GetPackageDetail != nil {
		return *m.GetPackageDetail
	}
	return 0
}

func (m *Request) GetChangeDestination() *PkgDest {
	if m != nil {
		return m.ChangeDestination
	}
	return nil
}

type Response struct {
	Error            *string    `protobuf:"bytes,1,opt,name=error" json:"error,omitempty"`
	UserId           *int64     `protobuf:"varint,2,opt,name=user_id" json:"user_id,omitempty"`
	PackageList      *PkgList   `protobuf:"bytes,3,opt,name=package_list" json:"package_list,omitempty"`
	PackageDetail    *PkgDetail `protobuf:"bytes,4,opt,name=package_detail" json:"package_detail,omitempty"`
	XXX_unrecognized []byte     `json:"-"`
}

func (m *Response) Reset()                    { *m = Response{} }
func (m *Response) String() string            { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()               {}
func (*Response) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *Response) GetError() string {
	if m != nil && m.Error != nil {
		return *m.Error
	}
	return ""
}

func (m *Response) GetUserId() int64 {
	if m != nil && m.UserId != nil {
		return *m.UserId
	}
	return 0
}

func (m *Response) GetPackageList() *PkgList {
	if m != nil {
		return m.PackageList
	}
	return nil
}

func (m *Response) GetPackageDetail() *PkgDetail {
	if m != nil {
		return m.PackageDetail
	}
	return nil
}

type PkgList struct {
	Total            *int64          `protobuf:"varint,1,req,name=total" json:"total,omitempty"`
	Packages         []*PkgList_Info `protobuf:"bytes,2,rep,name=packages" json:"packages,omitempty"`
	XXX_unrecognized []byte          `json:"-"`
}

func (m *PkgList) Reset()                    { *m = PkgList{} }
func (m *PkgList) String() string            { return proto.CompactTextString(m) }
func (*PkgList) ProtoMessage()               {}
func (*PkgList) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *PkgList) GetTotal() int64 {
	if m != nil && m.Total != nil {
		return *m.Total
	}
	return 0
}

func (m *PkgList) GetPackages() []*PkgList_Info {
	if m != nil {
		return m.Packages
	}
	return nil
}

type PkgList_Info struct {
	PackageId        *int64  `protobuf:"varint,1,req,name=package_id" json:"package_id,omitempty"`
	Status           *string `protobuf:"bytes,2,req,name=status" json:"status,omitempty"`
	CreateTime       *int64  `protobuf:"varint,3,req,name=create_time" json:"create_time,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *PkgList_Info) Reset()                    { *m = PkgList_Info{} }
func (m *PkgList_Info) String() string            { return proto.CompactTextString(m) }
func (*PkgList_Info) ProtoMessage()               {}
func (*PkgList_Info) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4, 0} }

func (m *PkgList_Info) GetPackageId() int64 {
	if m != nil && m.PackageId != nil {
		return *m.PackageId
	}
	return 0
}

func (m *PkgList_Info) GetStatus() string {
	if m != nil && m.Status != nil {
		return *m.Status
	}
	return ""
}

func (m *PkgList_Info) GetCreateTime() int64 {
	if m != nil && m.CreateTime != nil {
		return *m.CreateTime
	}
	return 0
}

type PkgDetail struct {
	Items            []*PkgDetail_Item   `protobuf:"bytes,1,rep,name=items" json:"items,omitempty"`
	X                *int32              `protobuf:"varint,2,req,name=x" json:"x,omitempty"`
	Y                *int32              `protobuf:"varint,3,req,name=y" json:"y,omitempty"`
	UserId           *int64              `protobuf:"varint,4,req,name=user_id" json:"user_id,omitempty"`
	Status           []*PkgDetail_Status `protobuf:"bytes,5,rep,name=status" json:"status,omitempty"`
	XXX_unrecognized []byte              `json:"-"`
}

func (m *PkgDetail) Reset()                    { *m = PkgDetail{} }
func (m *PkgDetail) String() string            { return proto.CompactTextString(m) }
func (*PkgDetail) ProtoMessage()               {}
func (*PkgDetail) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *PkgDetail) GetItems() []*PkgDetail_Item {
	if m != nil {
		return m.Items
	}
	return nil
}

func (m *PkgDetail) GetX() int32 {
	if m != nil && m.X != nil {
		return *m.X
	}
	return 0
}

func (m *PkgDetail) GetY() int32 {
	if m != nil && m.Y != nil {
		return *m.Y
	}
	return 0
}

func (m *PkgDetail) GetUserId() int64 {
	if m != nil && m.UserId != nil {
		return *m.UserId
	}
	return 0
}

func (m *PkgDetail) GetStatus() []*PkgDetail_Status {
	if m != nil {
		return m.Status
	}
	return nil
}

type PkgDetail_Item struct {
	Description      *string `protobuf:"bytes,1,req,name=description" json:"description,omitempty"`
	Amount           *int32  `protobuf:"varint,2,req,name=amount" json:"amount,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *PkgDetail_Item) Reset()                    { *m = PkgDetail_Item{} }
func (m *PkgDetail_Item) String() string            { return proto.CompactTextString(m) }
func (*PkgDetail_Item) ProtoMessage()               {}
func (*PkgDetail_Item) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5, 0} }

func (m *PkgDetail_Item) GetDescription() string {
	if m != nil && m.Description != nil {
		return *m.Description
	}
	return ""
}

func (m *PkgDetail_Item) GetAmount() int32 {
	if m != nil && m.Amount != nil {
		return *m.Amount
	}
	return 0
}

type PkgDetail_Status struct {
	Status           *string `protobuf:"bytes,1,req,name=status" json:"status,omitempty"`
	Timestamp        *int64  `protobuf:"varint,2,opt,name=timestamp" json:"timestamp,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *PkgDetail_Status) Reset()                    { *m = PkgDetail_Status{} }
func (m *PkgDetail_Status) String() string            { return proto.CompactTextString(m) }
func (*PkgDetail_Status) ProtoMessage()               {}
func (*PkgDetail_Status) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5, 1} }

func (m *PkgDetail_Status) GetStatus() string {
	if m != nil && m.Status != nil {
		return *m.Status
	}
	return ""
}

func (m *PkgDetail_Status) GetTimestamp() int64 {
	if m != nil && m.Timestamp != nil {
		return *m.Timestamp
	}
	return 0
}

func init() {
	proto.RegisterType((*PkgDest)(nil), "web.PkgDest")
	proto.RegisterType((*PkgListReq)(nil), "web.PkgListReq")
	proto.RegisterType((*Request)(nil), "web.Request")
	proto.RegisterType((*Response)(nil), "web.Response")
	proto.RegisterType((*PkgList)(nil), "web.PkgList")
	proto.RegisterType((*PkgList_Info)(nil), "web.PkgList.Info")
	proto.RegisterType((*PkgDetail)(nil), "web.PkgDetail")
	proto.RegisterType((*PkgDetail_Item)(nil), "web.PkgDetail.Item")
	proto.RegisterType((*PkgDetail_Status)(nil), "web.PkgDetail.Status")
}

func init() { proto.RegisterFile("web.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 431 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x92, 0x41, 0x6e, 0xdb, 0x30,
	0x10, 0x45, 0x41, 0xd1, 0xb2, 0xad, 0x51, 0xa2, 0x24, 0x34, 0x0a, 0x08, 0x5e, 0x09, 0x2a, 0x5a,
	0xa8, 0x08, 0xe0, 0x85, 0x2f, 0xd0, 0x4d, 0x36, 0x01, 0xba, 0x48, 0xdd, 0x03, 0x18, 0x8c, 0x3d,
	0x76, 0x09, 0x5b, 0xa2, 0x22, 0x8e, 0xea, 0x76, 0xd3, 0x13, 0xf4, 0x7e, 0xbd, 0x4e, 0xc1, 0x91,
	0x9d, 0xc8, 0x41, 0x17, 0x59, 0x99, 0x26, 0x3f, 0xff, 0x7f, 0xc3, 0x2f, 0x88, 0x0e, 0xf8, 0x38,
	0xab, 0x1b, 0x4b, 0x56, 0xc9, 0x03, 0x3e, 0xe6, 0x77, 0x30, 0x7a, 0xd8, 0x6d, 0xef, 0xd0, 0x91,
	0x52, 0x00, 0xb5, 0x5e, 0xed, 0xf4, 0x16, 0x97, 0x66, 0x9d, 0x8a, 0x2c, 0x28, 0xa4, 0x8a, 0x40,
	0xfc, 0x4c, 0x83, 0x2c, 0x28, 0x42, 0xbf, 0xfc, 0x95, 0x4a, 0x5e, 0x5e, 0xc1, 0xa8, 0x75, 0xd8,
	0x78, 0xd9, 0xc0, 0xcb, 0xf2, 0xaf, 0x00, 0x0f, 0xbb, 0xed, 0x17, 0xe3, 0x68, 0x81, 0x4f, 0xfd,
	0x63, 0x91, 0x89, 0x42, 0xaa, 0x09, 0xc4, 0x2f, 0xce, 0x2e, 0x0d, 0x32, 0x59, 0x48, 0x95, 0xc0,
	0xd0, 0x6e, 0x36, 0x0e, 0x29, 0x95, 0x2c, 0xba, 0x84, 0x70, 0x6f, 0x4a, 0x43, 0xe9, 0xc0, 0xff,
	0xcd, 0xff, 0x08, 0x18, 0x2d, 0xf0, 0xa9, 0xf5, 0x64, 0xd7, 0x30, 0xae, 0xf0, 0xb0, 0xf4, 0xa6,
	0xec, 0x18, 0xa9, 0x4f, 0x70, 0xbd, 0x45, 0x5a, 0x9e, 0x5c, 0xf7, 0xc6, 0x51, 0x1a, 0x64, 0xa2,
	0x88, 0xe7, 0x57, 0x33, 0x3f, 0x61, 0x8f, 0x66, 0x0a, 0xaa, 0x2f, 0x5d, 0x23, 0x69, 0xb3, 0x3f,
	0x66, 0x16, 0xa0, 0x56, 0xdf, 0x75, 0xc5, 0xdb, 0x8e, 0x4c, 0xa5, 0xc9, 0xd8, 0x8a, 0x01, 0xe2,
	0xf9, 0xc5, 0xc9, 0xc8, 0x3f, 0x4e, 0xfe, 0x03, 0xc6, 0x0b, 0x74, 0xb5, 0xad, 0x1c, 0x7a, 0x52,
	0x6c, 0x1a, 0x7b, 0x62, 0xe9, 0x8d, 0x1b, 0xb0, 0x6b, 0x0e, 0x17, 0x67, 0x60, 0xf2, 0xdc, 0xcf,
	0x83, 0xa9, 0x8f, 0x90, 0xbc, 0x22, 0xea, 0x52, 0x93, 0x97, 0x54, 0xbf, 0x9b, 0xff, 0xe6, 0x7e,
	0xf8, 0xca, 0x25, 0x84, 0x64, 0x49, 0xef, 0x8f, 0xd5, 0xbc, 0x87, 0xf1, 0xd1, 0xa1, 0x7b, 0xd1,
	0x78, 0x7e, 0xd3, 0x4f, 0x98, 0xdd, 0x57, 0x1b, 0x3b, 0xfd, 0x0c, 0x03, 0xff, 0xfb, 0xdf, 0x6e,
	0x13, 0x18, 0x3a, 0xd2, 0xd4, 0x3a, 0x2e, 0x38, 0xf2, 0x2d, 0xad, 0x1a, 0xd4, 0x84, 0x4b, 0x32,
	0x25, 0x72, 0xd5, 0x32, 0xff, 0x2b, 0x20, 0x7a, 0xa6, 0x51, 0x39, 0x84, 0x86, 0xb0, 0x74, 0xa9,
	0xe0, 0xc0, 0xc9, 0x39, 0xec, 0xec, 0x9e, 0xb0, 0x7c, 0xe3, 0x27, 0xa3, 0x3e, 0x3c, 0xa7, 0x87,
	0xec, 0xf5, 0xee, 0x95, 0xd7, 0x37, 0x3e, 0x9c, 0xde, 0xc2, 0x80, 0x5d, 0x27, 0x10, 0xaf, 0xd1,
	0xad, 0x1a, 0x53, 0x73, 0x45, 0x82, 0x89, 0x13, 0x18, 0xea, 0xd2, 0xb6, 0x15, 0x75, 0x79, 0xd3,
	0x5b, 0x18, 0x76, 0xd7, 0x7a, 0xb3, 0x75, 0xca, 0x1b, 0x88, 0xfc, 0x50, 0x8e, 0x74, 0x59, 0x77,
	0x2d, 0xfd, 0x0b, 0x00, 0x00, 0xff, 0xff, 0x1e, 0xb7, 0xdd, 0x4d, 0x0a, 0x03, 0x00, 0x00,
}

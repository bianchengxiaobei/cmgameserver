// Code generated by protoc-gen-go. DO NOT EDIT.
// source: G2M.proto

package message

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

type G2M_LoginToGameServer struct {
	RoleId               int64    `protobuf:"varint,1,opt,name=roleId,proto3" json:"roleId,omitempty"`
	UserId               int64    `protobuf:"varint,2,opt,name=userId,proto3" json:"userId,omitempty"`
	UserName             string   `protobuf:"bytes,3,opt,name=userName,proto3" json:"userName,omitempty"`
	GateId               int32    `protobuf:"varint,4,opt,name=gateId,proto3" json:"gateId,omitempty"`
	ServerId             int32    `protobuf:"varint,5,opt,name=serverId,proto3" json:"serverId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *G2M_LoginToGameServer) Reset()         { *m = G2M_LoginToGameServer{} }
func (m *G2M_LoginToGameServer) String() string { return proto.CompactTextString(m) }
func (*G2M_LoginToGameServer) ProtoMessage()    {}
func (*G2M_LoginToGameServer) Descriptor() ([]byte, []int) {
	return fileDescriptor_G2M_39211666745c7fe9, []int{0}
}
func (m *G2M_LoginToGameServer) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_G2M_LoginToGameServer.Unmarshal(m, b)
}
func (m *G2M_LoginToGameServer) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_G2M_LoginToGameServer.Marshal(b, m, deterministic)
}
func (dst *G2M_LoginToGameServer) XXX_Merge(src proto.Message) {
	xxx_messageInfo_G2M_LoginToGameServer.Merge(dst, src)
}
func (m *G2M_LoginToGameServer) XXX_Size() int {
	return xxx_messageInfo_G2M_LoginToGameServer.Size(m)
}
func (m *G2M_LoginToGameServer) XXX_DiscardUnknown() {
	xxx_messageInfo_G2M_LoginToGameServer.DiscardUnknown(m)
}

var xxx_messageInfo_G2M_LoginToGameServer proto.InternalMessageInfo

func (m *G2M_LoginToGameServer) GetRoleId() int64 {
	if m != nil {
		return m.RoleId
	}
	return 0
}

func (m *G2M_LoginToGameServer) GetUserId() int64 {
	if m != nil {
		return m.UserId
	}
	return 0
}

func (m *G2M_LoginToGameServer) GetUserName() string {
	if m != nil {
		return m.UserName
	}
	return ""
}

func (m *G2M_LoginToGameServer) GetGateId() int32 {
	if m != nil {
		return m.GateId
	}
	return 0
}

func (m *G2M_LoginToGameServer) GetServerId() int32 {
	if m != nil {
		return m.ServerId
	}
	return 0
}

type G2M_RoleRegisterToGateSuccess struct {
	RoleId               int64    `protobuf:"varint,1,opt,name=roleId,proto3" json:"roleId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *G2M_RoleRegisterToGateSuccess) Reset()         { *m = G2M_RoleRegisterToGateSuccess{} }
func (m *G2M_RoleRegisterToGateSuccess) String() string { return proto.CompactTextString(m) }
func (*G2M_RoleRegisterToGateSuccess) ProtoMessage()    {}
func (*G2M_RoleRegisterToGateSuccess) Descriptor() ([]byte, []int) {
	return fileDescriptor_G2M_39211666745c7fe9, []int{1}
}
func (m *G2M_RoleRegisterToGateSuccess) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_G2M_RoleRegisterToGateSuccess.Unmarshal(m, b)
}
func (m *G2M_RoleRegisterToGateSuccess) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_G2M_RoleRegisterToGateSuccess.Marshal(b, m, deterministic)
}
func (dst *G2M_RoleRegisterToGateSuccess) XXX_Merge(src proto.Message) {
	xxx_messageInfo_G2M_RoleRegisterToGateSuccess.Merge(dst, src)
}
func (m *G2M_RoleRegisterToGateSuccess) XXX_Size() int {
	return xxx_messageInfo_G2M_RoleRegisterToGateSuccess.Size(m)
}
func (m *G2M_RoleRegisterToGateSuccess) XXX_DiscardUnknown() {
	xxx_messageInfo_G2M_RoleRegisterToGateSuccess.DiscardUnknown(m)
}

var xxx_messageInfo_G2M_RoleRegisterToGateSuccess proto.InternalMessageInfo

func (m *G2M_RoleRegisterToGateSuccess) GetRoleId() int64 {
	if m != nil {
		return m.RoleId
	}
	return 0
}

type G2M_RoleQuitGameServer struct {
	RoleId               int64    `protobuf:"varint,1,opt,name=roleId,proto3" json:"roleId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *G2M_RoleQuitGameServer) Reset()         { *m = G2M_RoleQuitGameServer{} }
func (m *G2M_RoleQuitGameServer) String() string { return proto.CompactTextString(m) }
func (*G2M_RoleQuitGameServer) ProtoMessage()    {}
func (*G2M_RoleQuitGameServer) Descriptor() ([]byte, []int) {
	return fileDescriptor_G2M_39211666745c7fe9, []int{2}
}
func (m *G2M_RoleQuitGameServer) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_G2M_RoleQuitGameServer.Unmarshal(m, b)
}
func (m *G2M_RoleQuitGameServer) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_G2M_RoleQuitGameServer.Marshal(b, m, deterministic)
}
func (dst *G2M_RoleQuitGameServer) XXX_Merge(src proto.Message) {
	xxx_messageInfo_G2M_RoleQuitGameServer.Merge(dst, src)
}
func (m *G2M_RoleQuitGameServer) XXX_Size() int {
	return xxx_messageInfo_G2M_RoleQuitGameServer.Size(m)
}
func (m *G2M_RoleQuitGameServer) XXX_DiscardUnknown() {
	xxx_messageInfo_G2M_RoleQuitGameServer.DiscardUnknown(m)
}

var xxx_messageInfo_G2M_RoleQuitGameServer proto.InternalMessageInfo

func (m *G2M_RoleQuitGameServer) GetRoleId() int64 {
	if m != nil {
		return m.RoleId
	}
	return 0
}

func init() {
	proto.RegisterType((*G2M_LoginToGameServer)(nil), "message.G2M_LoginToGameServer")
	proto.RegisterType((*G2M_RoleRegisterToGateSuccess)(nil), "message.G2M_RoleRegisterToGateSuccess")
	proto.RegisterType((*G2M_RoleQuitGameServer)(nil), "message.G2M_RoleQuitGameServer")
}

func init() { proto.RegisterFile("G2M.proto", fileDescriptor_G2M_39211666745c7fe9) }

var fileDescriptor_G2M_39211666745c7fe9 = []byte{
	// 198 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x74, 0x37, 0xf2, 0xd5,
	0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0xcf, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0x55, 0x9a,
	0xce, 0xc8, 0x25, 0xea, 0x6e, 0xe4, 0x1b, 0xef, 0x93, 0x9f, 0x9e, 0x99, 0x17, 0x92, 0xef, 0x9e,
	0x98, 0x9b, 0x1a, 0x9c, 0x5a, 0x54, 0x96, 0x5a, 0x24, 0x24, 0xc6, 0xc5, 0x56, 0x94, 0x9f, 0x93,
	0xea, 0x99, 0x22, 0xc1, 0xa8, 0xc0, 0xa8, 0xc1, 0x1c, 0x04, 0xe5, 0x81, 0xc4, 0x4b, 0x8b, 0x53,
	0x8b, 0x3c, 0x53, 0x24, 0x98, 0x20, 0xe2, 0x10, 0x9e, 0x90, 0x14, 0x17, 0x07, 0x88, 0xe5, 0x97,
	0x98, 0x9b, 0x2a, 0xc1, 0xac, 0xc0, 0xa8, 0xc1, 0x19, 0x04, 0xe7, 0x83, 0xf4, 0xa4, 0x27, 0x96,
	0x80, 0xcc, 0x62, 0x51, 0x60, 0xd4, 0x60, 0x0d, 0x82, 0xf2, 0x40, 0x7a, 0x8a, 0xc1, 0xb6, 0x79,
	0xa6, 0x48, 0xb0, 0x82, 0x65, 0xe0, 0x7c, 0x25, 0x73, 0x2e, 0x59, 0x90, 0xc3, 0x82, 0xf2, 0x73,
	0x52, 0x83, 0x52, 0xd3, 0x33, 0x8b, 0x4b, 0x52, 0x8b, 0x40, 0xee, 0x2b, 0x49, 0x0d, 0x2e, 0x4d,
	0x4e, 0x4e, 0x2d, 0x2e, 0xc6, 0xe5, 0x40, 0x25, 0x03, 0x2e, 0x31, 0x98, 0xc6, 0xc0, 0xd2, 0xcc,
	0x12, 0xc2, 0x5e, 0x4a, 0x62, 0x03, 0x07, 0x8a, 0x31, 0x20, 0x00, 0x00, 0xff, 0xff, 0x5a, 0xf0,
	0x35, 0x3f, 0x21, 0x01, 0x00, 0x00,
}

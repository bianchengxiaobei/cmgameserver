// Code generated by protoc-gen-go. DO NOT EDIT.
// source: G2C.proto

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

type G2C_CharacterInfo struct {
	Role                 *Role    `protobuf:"bytes,1,opt,name=role,proto3" json:"role,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *G2C_CharacterInfo) Reset()         { *m = G2C_CharacterInfo{} }
func (m *G2C_CharacterInfo) String() string { return proto.CompactTextString(m) }
func (*G2C_CharacterInfo) ProtoMessage()    {}
func (*G2C_CharacterInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_G2C_8f8b2c49b422c98c, []int{0}
}
func (m *G2C_CharacterInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_G2C_CharacterInfo.Unmarshal(m, b)
}
func (m *G2C_CharacterInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_G2C_CharacterInfo.Marshal(b, m, deterministic)
}
func (dst *G2C_CharacterInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_G2C_CharacterInfo.Merge(dst, src)
}
func (m *G2C_CharacterInfo) XXX_Size() int {
	return xxx_messageInfo_G2C_CharacterInfo.Size(m)
}
func (m *G2C_CharacterInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_G2C_CharacterInfo.DiscardUnknown(m)
}

var xxx_messageInfo_G2C_CharacterInfo proto.InternalMessageInfo

func (m *G2C_CharacterInfo) GetRole() *Role {
	if m != nil {
		return m.Role
	}
	return nil
}

type Role struct {
	RoleId               int64    `protobuf:"varint,1,opt,name=roleId,proto3" json:"roleId,omitempty"`
	NickName             string   `protobuf:"bytes,2,opt,name=nickName,proto3" json:"nickName,omitempty"`
	AvatarId             int32    `protobuf:"varint,3,opt,name=avatarId,proto3" json:"avatarId,omitempty"`
	Level                int32    `protobuf:"varint,4,opt,name=level,proto3" json:"level,omitempty"`
	Exp                  int32    `protobuf:"varint,5,opt,name=exp,proto3" json:"exp,omitempty"`
	Gold                 int32    `protobuf:"varint,6,opt,name=gold,proto3" json:"gold,omitempty"`
	Diam                 int32    `protobuf:"varint,7,opt,name=diam,proto3" json:"diam,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Role) Reset()         { *m = Role{} }
func (m *Role) String() string { return proto.CompactTextString(m) }
func (*Role) ProtoMessage()    {}
func (*Role) Descriptor() ([]byte, []int) {
	return fileDescriptor_G2C_8f8b2c49b422c98c, []int{1}
}
func (m *Role) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Role.Unmarshal(m, b)
}
func (m *Role) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Role.Marshal(b, m, deterministic)
}
func (dst *Role) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Role.Merge(dst, src)
}
func (m *Role) XXX_Size() int {
	return xxx_messageInfo_Role.Size(m)
}
func (m *Role) XXX_DiscardUnknown() {
	xxx_messageInfo_Role.DiscardUnknown(m)
}

var xxx_messageInfo_Role proto.InternalMessageInfo

func (m *Role) GetRoleId() int64 {
	if m != nil {
		return m.RoleId
	}
	return 0
}

func (m *Role) GetNickName() string {
	if m != nil {
		return m.NickName
	}
	return ""
}

func (m *Role) GetAvatarId() int32 {
	if m != nil {
		return m.AvatarId
	}
	return 0
}

func (m *Role) GetLevel() int32 {
	if m != nil {
		return m.Level
	}
	return 0
}

func (m *Role) GetExp() int32 {
	if m != nil {
		return m.Exp
	}
	return 0
}

func (m *Role) GetGold() int32 {
	if m != nil {
		return m.Gold
	}
	return 0
}

func (m *Role) GetDiam() int32 {
	if m != nil {
		return m.Diam
	}
	return 0
}

func init() {
	proto.RegisterType((*G2C_CharacterInfo)(nil), "message.G2C_CharacterInfo")
	proto.RegisterType((*Role)(nil), "message.Role")
}

func init() { proto.RegisterFile("G2C.proto", fileDescriptor_G2C_8f8b2c49b422c98c) }

var fileDescriptor_G2C_8f8b2c49b422c98c = []byte{
	// 206 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x3c, 0x8f, 0xc1, 0x4e, 0x84, 0x30,
	0x10, 0x86, 0x53, 0x29, 0xac, 0x3b, 0xc6, 0x44, 0x27, 0xc6, 0x34, 0x9e, 0x70, 0x4f, 0x9c, 0x38,
	0x60, 0xe2, 0x0b, 0x70, 0xd8, 0x70, 0xf1, 0xd0, 0x17, 0x30, 0xe3, 0x76, 0x44, 0x62, 0xa1, 0xa4,
	0x10, 0xe2, 0x13, 0xf9, 0x9c, 0x9b, 0x16, 0xc2, 0xed, 0xff, 0xbe, 0x6f, 0x2e, 0x03, 0xc7, 0x73,
	0x55, 0x97, 0xa3, 0x77, 0xb3, 0xc3, 0x43, 0xcf, 0xd3, 0x44, 0x2d, 0x9f, 0xde, 0xe1, 0xf1, 0x5c,
	0xd5, 0x9f, 0xf5, 0x0f, 0x79, 0xba, 0xcc, 0xec, 0x9b, 0xe1, 0xdb, 0xe1, 0x2b, 0x48, 0xef, 0x2c,
	0x2b, 0x91, 0x8b, 0xe2, 0xae, 0xba, 0x2f, 0xb7, 0xe3, 0x52, 0x3b, 0xcb, 0x3a, 0xa6, 0xd3, 0xbf,
	0x00, 0x19, 0x10, 0x9f, 0x21, 0x0b, 0xa2, 0x31, 0xf1, 0x3a, 0xd1, 0x1b, 0xe1, 0x0b, 0xdc, 0x0e,
	0xdd, 0xe5, 0xf7, 0x83, 0x7a, 0x56, 0x37, 0xb9, 0x28, 0x8e, 0x7a, 0xe7, 0xd0, 0x68, 0xa1, 0x99,
	0x7c, 0x63, 0x54, 0x92, 0x8b, 0x22, 0xd5, 0x3b, 0xe3, 0x13, 0xa4, 0x96, 0x17, 0xb6, 0x4a, 0xc6,
	0xb0, 0x02, 0x3e, 0x40, 0xc2, 0x7f, 0xa3, 0x4a, 0xa3, 0x0b, 0x13, 0x11, 0x64, 0xeb, 0xac, 0x51,
	0x59, 0x54, 0x71, 0x07, 0x67, 0x3a, 0xea, 0xd5, 0x61, 0x75, 0x61, 0x7f, 0x65, 0xf1, 0xe1, 0xb7,
	0x6b, 0x00, 0x00, 0x00, 0xff, 0xff, 0x6f, 0x9b, 0x9a, 0x4a, 0xfd, 0x00, 0x00, 0x00,
}
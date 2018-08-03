// Code generated by protoc-gen-go. DO NOT EDIT.
// source: M2C.proto

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

type Hero struct {
	HeroId               int32    `protobuf:"varint,1,opt,name=heroId,proto3" json:"heroId,omitempty"`
	Level                int32    `protobuf:"varint,2,opt,name=level,proto3" json:"level,omitempty"`
	ItemIds              []int32  `protobuf:"varint,3,rep,packed,name=itemIds,proto3" json:"itemIds,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Hero) Reset()         { *m = Hero{} }
func (m *Hero) String() string { return proto.CompactTextString(m) }
func (*Hero) ProtoMessage()    {}
func (*Hero) Descriptor() ([]byte, []int) {
	return fileDescriptor_M2C_be74c7094c513b80, []int{0}
}
func (m *Hero) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Hero.Unmarshal(m, b)
}
func (m *Hero) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Hero.Marshal(b, m, deterministic)
}
func (dst *Hero) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Hero.Merge(dst, src)
}
func (m *Hero) XXX_Size() int {
	return xxx_messageInfo_Hero.Size(m)
}
func (m *Hero) XXX_DiscardUnknown() {
	xxx_messageInfo_Hero.DiscardUnknown(m)
}

var xxx_messageInfo_Hero proto.InternalMessageInfo

func (m *Hero) GetHeroId() int32 {
	if m != nil {
		return m.HeroId
	}
	return 0
}

func (m *Hero) GetLevel() int32 {
	if m != nil {
		return m.Level
	}
	return 0
}

func (m *Hero) GetItemIds() []int32 {
	if m != nil {
		return m.ItemIds
	}
	return nil
}

type Item struct {
	ItemId               int32    `protobuf:"varint,1,opt,name=itemId,proto3" json:"itemId,omitempty"`
	AttributeId          []int32  `protobuf:"varint,2,rep,packed,name=attributeId,proto3" json:"attributeId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Item) Reset()         { *m = Item{} }
func (m *Item) String() string { return proto.CompactTextString(m) }
func (*Item) ProtoMessage()    {}
func (*Item) Descriptor() ([]byte, []int) {
	return fileDescriptor_M2C_be74c7094c513b80, []int{1}
}
func (m *Item) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Item.Unmarshal(m, b)
}
func (m *Item) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Item.Marshal(b, m, deterministic)
}
func (dst *Item) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Item.Merge(dst, src)
}
func (m *Item) XXX_Size() int {
	return xxx_messageInfo_Item.Size(m)
}
func (m *Item) XXX_DiscardUnknown() {
	xxx_messageInfo_Item.DiscardUnknown(m)
}

var xxx_messageInfo_Item proto.InternalMessageInfo

func (m *Item) GetItemId() int32 {
	if m != nil {
		return m.ItemId
	}
	return 0
}

func (m *Item) GetAttributeId() []int32 {
	if m != nil {
		return m.AttributeId
	}
	return nil
}

type M2C_EnterLobby struct {
	IsInBattle bool `protobuf:"varint,1,opt,name=isInBattle,proto3" json:"isInBattle,omitempty"`
	// 玩家数据
	RoleBasicInfo *Role `protobuf:"bytes,2,opt,name=roleBasicInfo,proto3" json:"roleBasicInfo,omitempty"`
	// 玩家英雄数据
	HeroInfo             []*Hero  `protobuf:"bytes,3,rep,name=heroInfo,proto3" json:"heroInfo,omitempty"`
	Items                []*Item  `protobuf:"bytes,4,rep,name=items,proto3" json:"items,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *M2C_EnterLobby) Reset()         { *m = M2C_EnterLobby{} }
func (m *M2C_EnterLobby) String() string { return proto.CompactTextString(m) }
func (*M2C_EnterLobby) ProtoMessage()    {}
func (*M2C_EnterLobby) Descriptor() ([]byte, []int) {
	return fileDescriptor_M2C_be74c7094c513b80, []int{2}
}
func (m *M2C_EnterLobby) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_M2C_EnterLobby.Unmarshal(m, b)
}
func (m *M2C_EnterLobby) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_M2C_EnterLobby.Marshal(b, m, deterministic)
}
func (dst *M2C_EnterLobby) XXX_Merge(src proto.Message) {
	xxx_messageInfo_M2C_EnterLobby.Merge(dst, src)
}
func (m *M2C_EnterLobby) XXX_Size() int {
	return xxx_messageInfo_M2C_EnterLobby.Size(m)
}
func (m *M2C_EnterLobby) XXX_DiscardUnknown() {
	xxx_messageInfo_M2C_EnterLobby.DiscardUnknown(m)
}

var xxx_messageInfo_M2C_EnterLobby proto.InternalMessageInfo

func (m *M2C_EnterLobby) GetIsInBattle() bool {
	if m != nil {
		return m.IsInBattle
	}
	return false
}

func (m *M2C_EnterLobby) GetRoleBasicInfo() *Role {
	if m != nil {
		return m.RoleBasicInfo
	}
	return nil
}

func (m *M2C_EnterLobby) GetHeroInfo() []*Hero {
	if m != nil {
		return m.HeroInfo
	}
	return nil
}

func (m *M2C_EnterLobby) GetItems() []*Item {
	if m != nil {
		return m.Items
	}
	return nil
}

// 刷新房间列表
type M2C_RefreshRoomList struct {
	RoomList             []*Room  `protobuf:"bytes,1,rep,name=roomList,proto3" json:"roomList,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *M2C_RefreshRoomList) Reset()         { *m = M2C_RefreshRoomList{} }
func (m *M2C_RefreshRoomList) String() string { return proto.CompactTextString(m) }
func (*M2C_RefreshRoomList) ProtoMessage()    {}
func (*M2C_RefreshRoomList) Descriptor() ([]byte, []int) {
	return fileDescriptor_M2C_be74c7094c513b80, []int{3}
}
func (m *M2C_RefreshRoomList) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_M2C_RefreshRoomList.Unmarshal(m, b)
}
func (m *M2C_RefreshRoomList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_M2C_RefreshRoomList.Marshal(b, m, deterministic)
}
func (dst *M2C_RefreshRoomList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_M2C_RefreshRoomList.Merge(dst, src)
}
func (m *M2C_RefreshRoomList) XXX_Size() int {
	return xxx_messageInfo_M2C_RefreshRoomList.Size(m)
}
func (m *M2C_RefreshRoomList) XXX_DiscardUnknown() {
	xxx_messageInfo_M2C_RefreshRoomList.DiscardUnknown(m)
}

var xxx_messageInfo_M2C_RefreshRoomList proto.InternalMessageInfo

func (m *M2C_RefreshRoomList) GetRoomList() []*Room {
	if m != nil {
		return m.RoomList
	}
	return nil
}

// 改变Room状态信息
type M2C_JoinRoom struct {
	JoinerId             int64    `protobuf:"varint,1,opt,name=joinerId,proto3" json:"joinerId,omitempty"`
	JoinerName           string   `protobuf:"bytes,2,opt,name=joinerName,proto3" json:"joinerName,omitempty"`
	JoinerIconId         int32    `protobuf:"varint,3,opt,name=joinerIconId,proto3" json:"joinerIconId,omitempty"`
	GroupId              int32    `protobuf:"varint,4,opt,name=groupId,proto3" json:"groupId,omitempty"`
	RoomId               int32    `protobuf:"varint,5,opt,name=roomId,proto3" json:"roomId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *M2C_JoinRoom) Reset()         { *m = M2C_JoinRoom{} }
func (m *M2C_JoinRoom) String() string { return proto.CompactTextString(m) }
func (*M2C_JoinRoom) ProtoMessage()    {}
func (*M2C_JoinRoom) Descriptor() ([]byte, []int) {
	return fileDescriptor_M2C_be74c7094c513b80, []int{4}
}
func (m *M2C_JoinRoom) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_M2C_JoinRoom.Unmarshal(m, b)
}
func (m *M2C_JoinRoom) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_M2C_JoinRoom.Marshal(b, m, deterministic)
}
func (dst *M2C_JoinRoom) XXX_Merge(src proto.Message) {
	xxx_messageInfo_M2C_JoinRoom.Merge(dst, src)
}
func (m *M2C_JoinRoom) XXX_Size() int {
	return xxx_messageInfo_M2C_JoinRoom.Size(m)
}
func (m *M2C_JoinRoom) XXX_DiscardUnknown() {
	xxx_messageInfo_M2C_JoinRoom.DiscardUnknown(m)
}

var xxx_messageInfo_M2C_JoinRoom proto.InternalMessageInfo

func (m *M2C_JoinRoom) GetJoinerId() int64 {
	if m != nil {
		return m.JoinerId
	}
	return 0
}

func (m *M2C_JoinRoom) GetJoinerName() string {
	if m != nil {
		return m.JoinerName
	}
	return ""
}

func (m *M2C_JoinRoom) GetJoinerIconId() int32 {
	if m != nil {
		return m.JoinerIconId
	}
	return 0
}

func (m *M2C_JoinRoom) GetGroupId() int32 {
	if m != nil {
		return m.GroupId
	}
	return 0
}

func (m *M2C_JoinRoom) GetRoomId() int32 {
	if m != nil {
		return m.RoomId
	}
	return 0
}

// 准备成功
type M2C_ReadySuccess struct {
	RoleId               int64    `protobuf:"varint,1,opt,name=roleId,proto3" json:"roleId,omitempty"`
	Ready                bool     `protobuf:"varint,2,opt,name=ready,proto3" json:"ready,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *M2C_ReadySuccess) Reset()         { *m = M2C_ReadySuccess{} }
func (m *M2C_ReadySuccess) String() string { return proto.CompactTextString(m) }
func (*M2C_ReadySuccess) ProtoMessage()    {}
func (*M2C_ReadySuccess) Descriptor() ([]byte, []int) {
	return fileDescriptor_M2C_be74c7094c513b80, []int{5}
}
func (m *M2C_ReadySuccess) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_M2C_ReadySuccess.Unmarshal(m, b)
}
func (m *M2C_ReadySuccess) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_M2C_ReadySuccess.Marshal(b, m, deterministic)
}
func (dst *M2C_ReadySuccess) XXX_Merge(src proto.Message) {
	xxx_messageInfo_M2C_ReadySuccess.Merge(dst, src)
}
func (m *M2C_ReadySuccess) XXX_Size() int {
	return xxx_messageInfo_M2C_ReadySuccess.Size(m)
}
func (m *M2C_ReadySuccess) XXX_DiscardUnknown() {
	xxx_messageInfo_M2C_ReadySuccess.DiscardUnknown(m)
}

var xxx_messageInfo_M2C_ReadySuccess proto.InternalMessageInfo

func (m *M2C_ReadySuccess) GetRoleId() int64 {
	if m != nil {
		return m.RoleId
	}
	return 0
}

func (m *M2C_ReadySuccess) GetReady() bool {
	if m != nil {
		return m.Ready
	}
	return false
}

// 开始战斗加载
type M2C_StartBattleLoad struct {
	AllReady             bool     `protobuf:"varint,1,opt,name=allReady,proto3" json:"allReady,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *M2C_StartBattleLoad) Reset()         { *m = M2C_StartBattleLoad{} }
func (m *M2C_StartBattleLoad) String() string { return proto.CompactTextString(m) }
func (*M2C_StartBattleLoad) ProtoMessage()    {}
func (*M2C_StartBattleLoad) Descriptor() ([]byte, []int) {
	return fileDescriptor_M2C_be74c7094c513b80, []int{6}
}
func (m *M2C_StartBattleLoad) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_M2C_StartBattleLoad.Unmarshal(m, b)
}
func (m *M2C_StartBattleLoad) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_M2C_StartBattleLoad.Marshal(b, m, deterministic)
}
func (dst *M2C_StartBattleLoad) XXX_Merge(src proto.Message) {
	xxx_messageInfo_M2C_StartBattleLoad.Merge(dst, src)
}
func (m *M2C_StartBattleLoad) XXX_Size() int {
	return xxx_messageInfo_M2C_StartBattleLoad.Size(m)
}
func (m *M2C_StartBattleLoad) XXX_DiscardUnknown() {
	xxx_messageInfo_M2C_StartBattleLoad.DiscardUnknown(m)
}

var xxx_messageInfo_M2C_StartBattleLoad proto.InternalMessageInfo

func (m *M2C_StartBattleLoad) GetAllReady() bool {
	if m != nil {
		return m.AllReady
	}
	return false
}

// 开始战斗
type M2C_StartBattle struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *M2C_StartBattle) Reset()         { *m = M2C_StartBattle{} }
func (m *M2C_StartBattle) String() string { return proto.CompactTextString(m) }
func (*M2C_StartBattle) ProtoMessage()    {}
func (*M2C_StartBattle) Descriptor() ([]byte, []int) {
	return fileDescriptor_M2C_be74c7094c513b80, []int{7}
}
func (m *M2C_StartBattle) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_M2C_StartBattle.Unmarshal(m, b)
}
func (m *M2C_StartBattle) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_M2C_StartBattle.Marshal(b, m, deterministic)
}
func (dst *M2C_StartBattle) XXX_Merge(src proto.Message) {
	xxx_messageInfo_M2C_StartBattle.Merge(dst, src)
}
func (m *M2C_StartBattle) XXX_Size() int {
	return xxx_messageInfo_M2C_StartBattle.Size(m)
}
func (m *M2C_StartBattle) XXX_DiscardUnknown() {
	xxx_messageInfo_M2C_StartBattle.DiscardUnknown(m)
}

var xxx_messageInfo_M2C_StartBattle proto.InternalMessageInfo

// 每一帧的操作
type M2C_BattleFrame struct {
	FrameCount           int32      `protobuf:"varint,1,opt,name=frameCount,proto3" json:"frameCount,omitempty"`
	Cmd                  []*Command `protobuf:"bytes,2,rep,name=cmd,proto3" json:"cmd,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *M2C_BattleFrame) Reset()         { *m = M2C_BattleFrame{} }
func (m *M2C_BattleFrame) String() string { return proto.CompactTextString(m) }
func (*M2C_BattleFrame) ProtoMessage()    {}
func (*M2C_BattleFrame) Descriptor() ([]byte, []int) {
	return fileDescriptor_M2C_be74c7094c513b80, []int{8}
}
func (m *M2C_BattleFrame) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_M2C_BattleFrame.Unmarshal(m, b)
}
func (m *M2C_BattleFrame) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_M2C_BattleFrame.Marshal(b, m, deterministic)
}
func (dst *M2C_BattleFrame) XXX_Merge(src proto.Message) {
	xxx_messageInfo_M2C_BattleFrame.Merge(dst, src)
}
func (m *M2C_BattleFrame) XXX_Size() int {
	return xxx_messageInfo_M2C_BattleFrame.Size(m)
}
func (m *M2C_BattleFrame) XXX_DiscardUnknown() {
	xxx_messageInfo_M2C_BattleFrame.DiscardUnknown(m)
}

var xxx_messageInfo_M2C_BattleFrame proto.InternalMessageInfo

func (m *M2C_BattleFrame) GetFrameCount() int32 {
	if m != nil {
		return m.FrameCount
	}
	return 0
}

func (m *M2C_BattleFrame) GetCmd() []*Command {
	if m != nil {
		return m.Cmd
	}
	return nil
}

func init() {
	proto.RegisterType((*Hero)(nil), "message.Hero")
	proto.RegisterType((*Item)(nil), "message.Item")
	proto.RegisterType((*M2C_EnterLobby)(nil), "message.M2C_EnterLobby")
	proto.RegisterType((*M2C_RefreshRoomList)(nil), "message.M2C_RefreshRoomList")
	proto.RegisterType((*M2C_JoinRoom)(nil), "message.M2C_JoinRoom")
	proto.RegisterType((*M2C_ReadySuccess)(nil), "message.M2C_ReadySuccess")
	proto.RegisterType((*M2C_StartBattleLoad)(nil), "message.M2C_StartBattleLoad")
	proto.RegisterType((*M2C_StartBattle)(nil), "message.M2C_StartBattle")
	proto.RegisterType((*M2C_BattleFrame)(nil), "message.M2C_BattleFrame")
}

func init() { proto.RegisterFile("M2C.proto", fileDescriptor_M2C_be74c7094c513b80) }

var fileDescriptor_M2C_be74c7094c513b80 = []byte{
	// 456 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x53, 0xcb, 0x6e, 0xdb, 0x30,
	0x10, 0x84, 0x23, 0x2b, 0xb1, 0xd7, 0x49, 0x9b, 0xb2, 0x45, 0x41, 0xf8, 0x10, 0x18, 0xec, 0x25,
	0xbd, 0x18, 0xa8, 0xf2, 0x03, 0x46, 0x84, 0x3e, 0x54, 0xd8, 0x39, 0x30, 0xe8, 0xb9, 0xa0, 0xc5,
	0x75, 0xa2, 0x42, 0x12, 0x03, 0x92, 0x2e, 0x90, 0xbf, 0xe9, 0x4f, 0xf4, 0xff, 0x82, 0x25, 0x65,
	0x41, 0xf1, 0x8d, 0x33, 0x3b, 0x5c, 0xed, 0xec, 0x50, 0x30, 0xdd, 0x64, 0xf9, 0xf2, 0xc9, 0x1a,
	0x6f, 0xd8, 0x59, 0x83, 0xce, 0xa9, 0x07, 0x9c, 0x4f, 0xf3, 0x6c, 0x13, 0xb9, 0xf9, 0xf4, 0xfb,
	0xa1, 0x2c, 0xee, 0x60, 0xfc, 0x03, 0xad, 0x61, 0x1f, 0xe1, 0xf4, 0x11, 0xad, 0x29, 0x34, 0x1f,
	0x2d, 0x46, 0xd7, 0xa9, 0xec, 0x10, 0xfb, 0x00, 0x69, 0x8d, 0x7f, 0xb1, 0xe6, 0x27, 0x81, 0x8e,
	0x80, 0x71, 0x38, 0xab, 0x3c, 0x36, 0x85, 0x76, 0x3c, 0x59, 0x24, 0xd7, 0xa9, 0x3c, 0x40, 0xb1,
	0x82, 0x71, 0xe1, 0xb1, 0xa1, 0x7e, 0x91, 0x3a, 0xf4, 0x8b, 0x88, 0x2d, 0x60, 0xa6, 0xbc, 0xb7,
	0xd5, 0x76, 0xef, 0xb1, 0xd0, 0xfc, 0x24, 0xdc, 0x1e, 0x52, 0xe2, 0xff, 0x08, 0xde, 0x6c, 0xb2,
	0xfc, 0xf7, 0xd7, 0xd6, 0xa3, 0x5d, 0x9b, 0xed, 0xf6, 0x99, 0x5d, 0x01, 0x54, 0xae, 0x68, 0x6f,
	0x95, 0xf7, 0x35, 0x86, 0x86, 0x13, 0x39, 0x60, 0xd8, 0x0d, 0x5c, 0x58, 0x53, 0xe3, 0xad, 0x72,
	0x55, 0x59, 0xb4, 0x3b, 0x13, 0x86, 0x9d, 0x65, 0x17, 0xcb, 0xce, 0xfb, 0x52, 0x9a, 0x1a, 0xe5,
	0x6b, 0x0d, 0xfb, 0x0c, 0x93, 0xe0, 0x91, 0xf4, 0x64, 0x62, 0xa8, 0xa7, 0x95, 0xc8, 0xbe, 0xcc,
	0x3e, 0x41, 0x4a, 0xe3, 0x3b, 0x3e, 0x3e, 0xd2, 0x91, 0x55, 0x19, 0x6b, 0x62, 0x05, 0xef, 0x69,
	0x6c, 0x89, 0x3b, 0x8b, 0xee, 0x51, 0x1a, 0xd3, 0xac, 0x2b, 0xe7, 0xe9, 0x33, 0xb6, 0x3b, 0xf3,
	0xd1, 0xd1, 0x75, 0x12, 0xc9, 0xbe, 0x2c, 0xfe, 0x8d, 0xe0, 0x9c, 0x5a, 0xfc, 0x34, 0x55, 0x4b,
	0x25, 0x36, 0x87, 0xc9, 0x1f, 0x53, 0xb5, 0x68, 0xbb, 0x35, 0x26, 0xb2, 0xc7, 0xb4, 0x93, 0x78,
	0xbe, 0x53, 0x0d, 0x06, 0xc3, 0x53, 0x39, 0x60, 0x98, 0x80, 0xf3, 0x4e, 0x5b, 0x9a, 0xb6, 0xd0,
	0x3c, 0x09, 0x31, 0xbc, 0xe2, 0x28, 0xc6, 0x07, 0x6b, 0xf6, 0x4f, 0x85, 0xe6, 0xe3, 0x50, 0x3e,
	0x40, 0x8a, 0x8f, 0xc6, 0x2a, 0x34, 0x4f, 0x63, 0x7c, 0x11, 0x89, 0x15, 0x5c, 0x46, 0x93, 0x4a,
	0x3f, 0xdf, 0xef, 0xcb, 0x12, 0x9d, 0x8b, 0xda, 0x1a, 0xfb, 0x19, 0x3b, 0x44, 0x4f, 0xc7, 0x92,
	0x2e, 0x0c, 0x37, 0x91, 0x11, 0x88, 0x2f, 0x71, 0x4d, 0xf7, 0x5e, 0x59, 0x1f, 0xe3, 0x5b, 0x1b,
	0xa5, 0xc9, 0xaa, 0xaa, 0xeb, 0xd0, 0xb7, 0x0b, 0xb8, 0xc7, 0xe2, 0x1d, 0xbc, 0x3d, 0xba, 0x22,
	0x7e, 0x45, 0x2a, 0xa2, 0x6f, 0x96, 0x0c, 0x5f, 0x01, 0xec, 0xe8, 0x90, 0x9b, 0x7d, 0xeb, 0xbb,
	0x57, 0x37, 0x60, 0x98, 0x80, 0xa4, 0x6c, 0xe2, 0x8b, 0x9b, 0x65, 0x97, 0x7d, 0x06, 0xb9, 0x69,
	0x1a, 0xd5, 0x6a, 0x49, 0xc5, 0xed, 0x69, 0xf8, 0x29, 0x6e, 0x5e, 0x02, 0x00, 0x00, 0xff, 0xff,
	0xbb, 0x56, 0x81, 0x3e, 0x40, 0x03, 0x00, 0x00,
}

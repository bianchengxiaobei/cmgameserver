package face

import "cmgameserver/message"

type IRoom interface {
	GetRoomName() string
	SetRoomName(name string)
	GetRoomId() int32
	GetRoomOwnerId() int64
	GetRoomMaxPlayerNum() int32
	SetRoomMaxPlayerNum(num int32)
	SetGameType(gameType int32)
	GetGameType() int32
	SetMapId(mapId int32)
	GetMapId() int32
	JoinOneMember(roleId int64) (int32, bool)
	LeaveOneMember(roleId int64) bool
	GetRoomRoleIds() [4]int64
	GetCurPlayerNum() int32
	GetRoomOwnerGroupId() int32
	SetRoomOwnerGroupId(id int32)
	GetRoomOwnerName()string
	CheckRoomReady() bool
	SetRoomMemberReady(ready bool, roleId int64) bool
	IsRoomOwner(roleId int64) bool
	SetIsWarFow(value bool)
	GetIsWarFow() bool
	GetSeed()int32
	SetSeed(seed int32)
	GetRoomOwnerAvatarId()int32
	GetRoomMemberGroupId(roleId int64) int32
	GetRoomMemberReady(roleId int64) bool
	GetArrowerData() *message.FreeSoldierData
	GetDaodunData() *message.FreeSoldierData
	GetSpearData() *message.FreeSoldierData
	GetFashiData() *message.FreeSoldierData
}

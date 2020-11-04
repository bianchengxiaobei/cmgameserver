package face

import "cmgameserver/message"

type BattleType int32

const(
	SimulateBattleType BattleType = 0
	BattleBattleType BattleType = 1
	PaiWeiBattleType BattleType = 2
	FreeRoomBattleType BattleType = 3
)

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
	GetInBattle()bool
	SetInBattle(value bool)
	JoinOneMember(roleId int64) (int32, bool)
	LeaveOneMember(roleId int64) bool
	GetRoomRoleIds() [4]int64
	GetCurPlayerNum() int32
	GetRoomOwnerGroupId() int32
	SetRoomOwnerGroupId(id int32)
	SetRoomOwnerCityId(city int32)
	GetRoomOwnerCityId() int32
	GetRoomOwnerLevel()int32
	GetRoomOwnerName()string
	CheckRoomReady() bool
	SetRoomMemberReady(ready bool, roleId int64) bool
	SetRoomMemberCityId(cityId int32,roleId int64)bool
	IsRoomOwner(roleId int64) bool
	SetIsWarFow(value bool)
	GetIsWarFow() bool
	GetIsOutsideMonster() bool
	SetIsOutsideMonster(value bool)
	GetSeed()int32
	SetSeed(seed int32)
	GetRoomOwnerAvatarId()int32
	GetRoomMemberGroupId(roleId int64) int32
	GetRoomMemberReady(roleId int64) bool
	GetRoomMemberCityId(roleId int64)int32
	GetArrowerData() *message.FreeSoldierData
	GetDaodunData() *message.FreeSoldierData
	GetSpearData() *message.FreeSoldierData
	GetFashiData() *message.FreeSoldierData
}

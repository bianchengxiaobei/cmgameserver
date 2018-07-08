package face

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
	JoinOneMember(roleId int64) (int32,bool)
	LeaveOneMember(roleId int64) bool
	GetRoomRoleIds() [4]int64
	GetCurPlayerNum() int32
	GetRoomOwnerGroupId() int32
	SetRoomOwnerGroupId(id int32)
	CheckRoomReady()bool
	SetRoomMemberReady(ready bool,roleId int64) bool
}

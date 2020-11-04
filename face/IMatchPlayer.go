package face

type IMatchPlayer interface {
	SetBOwner(value bool)
	GetBOwner()bool
	GetRoleId()int64
	GetOnlineRole()IOnlineRole
	GetPosIndex()int32
	SetPosIndex(value int32)
	GetMatchTeamId()int32
	SetMatchTeamId(value int32)
	SetGroupId(value int32)
	GetGroupId()int32
	SetMatchRoomId(id int32)
	GetMatchRoomId()int32
	SetPrepare(value bool)
	GetPrepare()bool
	GetBInBattleRoom()bool
	SetBInBattleRoom(value bool)
	Clear()
	GetBInMatching()bool
	SetBInMatching(value bool)
} 
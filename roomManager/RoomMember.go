package roomManager


type RoomMember struct {
	RoleId	int64
	GroupId	int32//队伍id
	Prepare	bool//是否准备好
}

func (mem RoomMember) GetGroupId() int32 {
	return mem.GroupId
}

func (mem RoomMember) GetPrepare() bool {
	return mem.Prepare
}

package roomManager


type RoomMember struct {
	RoleId	int64
	GroupId	int32//队伍id
	Prepare	bool//是否准备好
	CityId int32//选择的城市id
}

func (mem RoomMember) GetGroupId() int32 {
	return mem.GroupId
}

func (mem RoomMember) GetPrepare() bool {
	return mem.Prepare
}
func (mem RoomMember) GetCityId() int32 {
	return mem.CityId
}

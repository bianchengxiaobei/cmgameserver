package face

type IRoomMember interface {
	GetGroupId()int32
	GetPrepare()bool
	GetCityId() int32
}

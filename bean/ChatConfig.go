package bean
//保存聊天记录，重启服务器的时候

type ChatConfig struct {
	Datas  [10]ChatData
}
type ChatData struct {
	AvatarId             int32
	Name                 string
	Chat                 string
	Time                 int64
	Level int32
	RankScore int32
	RoomId int32
	Sex int32
}
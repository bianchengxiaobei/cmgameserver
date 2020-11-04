package face

import "github.com/golang/protobuf/proto"

type IMatchBattleRoom interface {
	AddAllPlayer(all map[int]IMatchTeam)
	SetPlayerPrepare(roleId int64,value bool)bool
	SendMstToAllPlayer(msgId int,msg proto.Message)
	CheckAllPlayerLoadFinished()(bool,*[4]IOnlineRole)
	Clear()
	ClearNotifyAllPlayer()
	GetRoomId()int32
	GetRoomSeed()int32
	GetAllPlayer()[]IMatchPlayer
}

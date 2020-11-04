package face

import "github.com/golang/protobuf/proto"


type PaiWeiGameMode int32

const (
	Classics PaiWeiGameMode = 0//经典排位
	TaFang PaiWeiGameMode = 1//塔防排位
)
type IMatchManager interface {
	ReqAutoCreateMatchTeam(role IOnlineRole)
	TeamStartMatching(teamId int32,mode PaiWeiGameMode)
	RemovePlayerFromMatchTeam(matchPlayer IMatchPlayer)
	EnterMatchBattleRoom(teams map[int]IMatchTeam)
	SetMatchRoomPlayerPrepare(roomId int32,roleId int64,value bool)bool
	SendMsgToAllBattleRoomPlayer(msgId int,msg proto.Message,roomId int32)
	CheckAllPlayerLoadFinished(roomId int32)(bool,*[4]IOnlineRole)
	CancelStartMatch(matchPlayer IMatchPlayer)
	OnePlayerQuitMatchBattleRoom(player IMatchPlayer)
	GetMatchBattleRoom(roomId int32)IMatchBattleRoom
}
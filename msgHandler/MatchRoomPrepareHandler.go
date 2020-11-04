package msgHandler

import (
	"cmgameserver/face"
	"github.com/bianchengxiaobei/cmgo/network"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/log4g"
)

type MatchRoomPrepareHandler struct {
	GameServer face.IGameServer
}

func (handler *MatchRoomPrepareHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if msg, ok := innerMsg.MsgData.(*message.C2M2C_MatchRoomPrepare); ok {
			role := handler.GameServer.GetRoleManager().GetOnlineRole(msg.RoleId)
			if role != nil {
				matchPlayer := role.GetMatchPlayer()
				if matchPlayer != nil{
					roomId := matchPlayer.GetMatchRoomId()
					if handler.GameServer.GetMatchManager().SetMatchRoomPlayerPrepare(roomId,role.GetRoleId(),msg.Value){
						handler.GameServer.GetMatchManager().SendMsgToAllBattleRoomPlayer(5069,msg,roomId)
					}
				}
			}
		} else {
			log4g.Error("不是C2M2C_MatchRoomPrepare！")
		}
	}
}


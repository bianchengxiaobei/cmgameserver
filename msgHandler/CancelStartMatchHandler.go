package msgHandler

import (
	"cmgameserver/face"
	"github.com/bianchengxiaobei/cmgo/network"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/log4g"
)

type CancelStartMatchHandler struct {
	GameServer face.IGameServer
}

func (handler *CancelStartMatchHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if _, ok := innerMsg.MsgData.(*message.C2M_CancelStartMatch); ok {
			role := handler.GameServer.GetRoleManager().GetOnlineRole(innerMsg.RoleId)
			if role != nil {
				matchPlayer := role.GetMatchPlayer()
				if matchPlayer != nil{
					//匹配管理器取消正在匹配
					handler.GameServer.GetMatchManager().CancelStartMatch(matchPlayer)
				}
			}
		} else {
			log4g.Error("不是C2M_CancelStartMatch！")
		}
	}
}
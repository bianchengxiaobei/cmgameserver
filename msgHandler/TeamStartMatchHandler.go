package msgHandler

import (
	"cmgameserver/face"
	"github.com/bianchengxiaobei/cmgo/network"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/log4g"
)

type TeamStartMatchHandler struct {
	GameServer face.IGameServer
}

func (handler *TeamStartMatchHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if _, ok := innerMsg.MsgData.(*message.C2M_TeamStartMatch); ok {
			role := handler.GameServer.GetRoleManager().GetOnlineRole(innerMsg.RoleId)
			if role != nil {
				matchPlayer := role.GetMatchPlayer()
				if matchPlayer != nil{
					//匹配管理器加入该匹配队伍(默认是经典排位)
					handler.GameServer.GetMatchManager().
						TeamStartMatching(matchPlayer.GetMatchTeamId(),face.Classics)
				}
			}
		} else {
			log4g.Error("不是C2M_TeamStartMatch！")
		}
	}
}

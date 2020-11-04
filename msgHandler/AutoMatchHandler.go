package msgHandler

import (
	"cmgameserver/face"
	"github.com/bianchengxiaobei/cmgo/network"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/log4g"
)

type AutoMatchHandler struct {
	GameServer face.IGameServer
}

func (handler *AutoMatchHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if _, ok := innerMsg.MsgData.(*message.C2M_ReqAutoMatch); ok {
			role := handler.GameServer.GetRoleManager().GetOnlineRole(innerMsg.RoleId)
			if role != nil {
				//匹配管理器加入该匹配队伍
				handler.GameServer.GetMatchManager().ReqAutoCreateMatchTeam(role)
			}
		} else {
			log4g.Error("不是C2M_AgreePauseBattle！")
		}
	}
}
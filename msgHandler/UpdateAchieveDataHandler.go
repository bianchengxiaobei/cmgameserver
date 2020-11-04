package msgHandler

import (
	"cmgameserver/face"
	"github.com/bianchengxiaobei/cmgo/network"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/log4g"
)

type UpdateAchieveDataHandler struct {
	GameServer face.IGameServer
}

func (handler *UpdateAchieveDataHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if protoMsg, ok := innerMsg.MsgData.(*message.C2M_UpdateAchievementData); ok {
			role := handler.GameServer.GetRoleManager().GetOnlineRole(innerMsg.RoleId)
			if role != nil {
				role.SetAchieveRecord(protoMsg.Data)
			}
		} else {
			log4g.Error("不是C2M_UpdateAchievementData！")
		}
	}
}


package msgHandler

import (
	"cmgameserver/face"
	"github.com/bianchengxiaobei/cmgo/network"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/log4g"
)

type EnterBattleStateHandler struct {
	GameServer face.IGameServer
}

func (handler *EnterBattleStateHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if protoMsg, ok := innerMsg.MsgData.(*message.C2M_EnterBattleState); ok {
			role := handler.GameServer.GetRoleManager().GetOnlineRole(innerMsg.RoleId)
			if role != nil {
				role.SetInSimulateBattle(protoMsg.Value)
			}
		} else {
			log4g.Error("C2M_EnterBattleState！")
		}
	}
}
package msgHandler

import (
	"cmgameserver/face"
	"github.com/bianchengxiaobei/cmgo/network"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/log4g"
)

type GetSingHandler struct {
	GameServer face.IGameServer
}

func (handler *GetSingHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if protoMsg, ok := innerMsg.MsgData.(*message.C2M2C_GetSign); ok {
			role := handler.GameServer.GetRoleManager().GetOnlineRole(innerMsg.RoleId)
			if role != nil {
				//看角色是否已经领取过了
				if role.GetSignAward(){
					//回送
					handler.GameServer.WriteInnerMsg(role.GetGateSession(), role.GetRoleId(), 5038, protoMsg)
				}
			} else {
				log4g.Errorf("不存在RoleId:%d", innerMsg.RoleId)
			}
		} else {
			log4g.Error("不是C2M2C_GetSign！")
		}
	}
}
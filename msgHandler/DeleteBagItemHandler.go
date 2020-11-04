package msgHandler

import (
	"cmgameserver/face"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/log4g"
	"github.com/bianchengxiaobei/cmgo/network"
)

type DeleteBagItemHandler struct {
	GameServer face.IGameServer
}

func (handler *DeleteBagItemHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if protoMsg, ok := innerMsg.MsgData.(*message.C2M2C_DeleteBagItem); ok {
			role := handler.GameServer.GetRoleManager().GetOnlineRole(innerMsg.RoleId)
			if role != nil {
				role.DeleteItem(protoMsg.ItemIndex)
				//暂时不回送
				//handler.GameServer.WriteInnerMsg(role.GetGateSession(), role.GetRoleId(), 5096, protoMsg)
			} else {
				log4g.Errorf("不存在RoleId:%d", innerMsg.RoleId)
			}
		} else {
			log4g.Error("不是C2M_GetBoxAwardItem！")
		}
	}
}

package msgHandler

import (
	"cmgameserver/face"
	"github.com/bianchengxiaobei/cmgo/network"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/log4g"
)

type BuyShopDiamHandler struct {
	GameServer face.IGameServer
}

func (handler *BuyShopDiamHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if protoMsg, ok := innerMsg.MsgData.(*message.C2M_BuyShopDiam); ok {
			if handler.GameServer.GetRoleManager().BuyDiamIOS(innerMsg.RoleId,protoMsg.DiamType,protoMsg.TranId) == false{
				role := handler.GameServer.GetRoleManager().GetOnlineRole(innerMsg.RoleId)
				if role != nil {
					RMsg := new(message.M2C_BuyDiamResult)
					RMsg.DiamType = protoMsg.DiamType
					RMsg.AddDiam = 0
					handler.GameServer.WriteInnerMsg(role.GetGateSession(),role.GetRoleId(),5081,RMsg)
				}
			}
		} else {
			log4g.Error("不是C2M_BuyShopDiam！")
		}
	}
}

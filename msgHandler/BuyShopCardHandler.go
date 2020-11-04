package msgHandler

import (
	"cmgameserver/face"
	"github.com/bianchengxiaobei/cmgo/network"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/log4g"
)

type BuyShopCardHandler struct {
	GameServer face.IGameServer
}

func (handler *BuyShopCardHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if protoMsg, ok := innerMsg.MsgData.(*message.C2M_BuyShopCard); ok {
			if handler.GameServer.GetRoleManager().BuyCardIOS(innerMsg.RoleId,protoMsg.CardType,protoMsg.TranId) == false{
				role := handler.GameServer.GetRoleManager().GetOnlineRole(innerMsg.RoleId)
				if role != nil {
					RMsg := new(message.M2C_BuyCardResult)
					RMsg.CardType = protoMsg.CardType
					RMsg.Success = false
					handler.GameServer.WriteInnerMsg(role.GetGateSession(),role.GetRoleId(),5079,RMsg)
				}
			}
		} else {
			log4g.Error("不是C2M_BuyShopCard！")
		}
	}
}

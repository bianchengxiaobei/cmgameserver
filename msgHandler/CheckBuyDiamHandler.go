package msgHandler

import (
	"cmgameserver/face"
	"github.com/bianchengxiaobei/cmgo/network"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/log4g"
	"cmgameserver/bean"
	"cmgameserver/tool"
)

type CheckBuyDiamHandler struct {
	GameServer face.IGameServer
}

func (handler *CheckBuyDiamHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if protoMsg, ok := innerMsg.MsgData.(*message.C2M_CheckBuyShopDiam); ok {
			role := handler.GameServer.GetRoleManager().GetOnlineRole(innerMsg.RoleId)
			if role != nil {
				success,tranType,tranValue := role.RemoveTran(protoMsg.TranId)
				if success == false{
					log4g.Infof("移除订单失败:[%s]",protoMsg.TranId)
					return
				}
				if tranType == bean.DiamTranType{
					addDiam := tool.GetDiamValueByDiamType(tranValue)
					role.AddDiam(addDiam)
				}else {
					if role.BuyCard(bean.CardType(tranValue)) == false{
						log4g.Infof("增加功能卡失败:[%s]",protoMsg.TranId)
					}
				}
			}
		} else {
			log4g.Error("不是C2M_CheckBuyShopDiam！")
		}
	}
}
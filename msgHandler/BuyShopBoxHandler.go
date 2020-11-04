package msgHandler

import (
	"cmgameserver/face"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/log4g"
	"github.com/bianchengxiaobei/cmgo/network"
	"cmgameserver/tool"
	"time"
)

type BuyShopBoxHandler struct {
	GameServer face.IGameServer
}

func (handler *BuyShopBoxHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if protoMsg, ok := innerMsg.MsgData.(*message.C2M_BuyShopBox); ok {
			role := handler.GameServer.GetRoleManager().GetOnlineRole(innerMsg.RoleId)
			if role != nil {
				value := tool.GetBoxPriceByBoxId(protoMsg.BoxId,protoMsg.BuyType)
				RMsg := new(message.M2C_BuyShopBoxResult)
				RMsg.BuyType = protoMsg.BuyType
				RMsg.BoxId = protoMsg.BoxId
				if protoMsg.BuyType == 0{
					//铜币购买
					if role.GetGold() >= value{
						role.AddGold(-value)
						RMsg.RoleValue = role.GetGold()
					}else{
						RMsg.BoxId = -1
					}
				}else {
					//钻石
					if role.GetDiam() >= value{
						role.AddDiam(-value)
						RMsg.RoleValue = role.GetDiam()
					}else{
						RMsg.BoxId = -1
					}
				}
				if RMsg.BoxId > 0{
					now := time.Now().Unix()
					seed := int32(now)
					index := role.AddItemNoMsg(protoMsg.BoxId,seed,now,false)
					if index < 0{
						//背包已经满了
						RMsg.BoxId = -1
						if protoMsg.BuyType == 0{
							//铜币购买
							role.AddGold(value)
						}else {
							//钻石
							role.AddDiam(value)
						}
					}
					RMsg.BoxSeed = seed
					RMsg.BoxIndex = index
				}
				handler.GameServer.WriteInnerMsg(role.GetGateSession(), role.GetRoleId(), 5087, RMsg)
			}
		} else {
			log4g.Error("不是C2M_BuyShopCard！")
		}
	}
}

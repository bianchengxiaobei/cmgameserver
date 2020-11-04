package msgHandler

import (
	"cmgameserver/face"
	"github.com/bianchengxiaobei/cmgo/network"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/log4g"
)

type BuyShopHeroCardHandler struct {
	GameServer face.IGameServer
}

func (handler *BuyShopHeroCardHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if protoMsg, ok := innerMsg.MsgData.(*message.C2M_BuyHeroCard); ok {
			//通过配置文件找到价格，然后判断是否可以买
			//如果可以买，那就添加item到背包，返回消息
			role := handler.GameServer.GetRoleManager().GetOnlineRole(innerMsg.RoleId)
			if role == nil{
				return
			}
			data := handler.GameServer.GetShopHeroCardConfig().Cards[protoMsg.HeroCardId]
			if &data == nil{
				return
			}
			diam := data.Price * protoMsg.CardNum
			if diam <= 0{
				log4g.Info("diam==0")
				return
			}
			msg := new(message.M2C_BuyHeroCardResult)
			itemId := protoMsg.HeroCardId
			if itemId <= 0{
				msg.HeroCardId = -1
				handler.GameServer.WriteInnerMsg(role.GetGateSession(),role.GetRoleId(),5089,msg)
				return
			}
			if role.GetDiam() < diam{
				msg.HeroCardId = -1
				handler.GameServer.WriteInnerMsg(role.GetGateSession(),role.GetRoleId(),5089,msg)
				return
			}
			role.AddDiam(-diam)
			num := int(protoMsg.CardNum)
			var index int32
			for i:=0;i<num;i++{
				index = role.AddItemNoMsg(itemId,0,0,true)
			}
			item := role.GetItem(index)
			if &item == nil{
				msg.HeroCardId = -1
				handler.GameServer.WriteInnerMsg(role.GetGateSession(),role.GetRoleId(),5089,msg)
				return
			}
			msg.HeroCardId = protoMsg.HeroCardId
			msg.ItemIndex = index
			msg.ItemNum = item.ItemNum
			msg.RoleDiamValue = role.GetDiam()
			handler.GameServer.WriteInnerMsg(role.GetGateSession(),role.GetRoleId(),5089,msg)
		} else {
			log4g.Error("不是C2M_BuyHeroCard！")
		}
	}
}

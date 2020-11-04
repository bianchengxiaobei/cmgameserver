package msgHandler

import (
	"cmgameserver/face"
	"github.com/bianchengxiaobei/cmgo/network"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/log4g"
	"time"
)

type BuyShopEquipHandler struct {
	GameServer face.IGameServer
}
//购买商城装备处理器
func (handler *BuyShopEquipHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if protoMsg, ok := innerMsg.MsgData.(*message.C2M_BuyShopEquip); ok {
			role := handler.GameServer.GetRoleManager().GetOnlineRole(innerMsg.RoleId)
			if role != nil {
				returnMsg := new(message.M2C_BuyShopEquipResult)
				returnMsg.EquipId = -1
				config := handler.GameServer.GetShopEquipConfig()
				if &config == nil{
					handler.GameServer.WriteInnerMsg(role.GetGateSession(), role.GetRoleId(), 5078, returnMsg)
					return
				}
				data,ok := config.ShopEquip[protoMsg.EquipId]
				if ok == false{
					handler.GameServer.WriteInnerMsg(role.GetGateSession(), role.GetRoleId(), 5078, returnMsg)
					log4g.Infof("购买失败,没有该商品[%d]",protoMsg.EquipId)
					return
				}
				if role.GetDiam() < data.Price{
					//玩家的钻石不够，购买失败
					handler.GameServer.WriteInnerMsg(role.GetGateSession(), role.GetRoleId(), 5078, returnMsg)
					return
				}else{
					//扣除钻石，然后发送给客户端购买成功
					//增加装备到玩家背包
					nowTime := time.Now()
					itemSeed := int32(nowTime.Unix())
					index := role.AddItemNoMsg(protoMsg.EquipId,itemSeed,nowTime.Unix(),false)
					role.AddDiam(-data.Price)
					returnMsg.EquipId = protoMsg.EquipId
					returnMsg.ItemIndex = index
					returnMsg.ItemSeed =itemSeed
					handler.GameServer.WriteInnerMsg(role.GetGateSession(), role.GetRoleId(), 5078, returnMsg)
				}
			}
		} else {
			log4g.Error("不是C2M_BuyShopEquip！")
		}
	}
}
